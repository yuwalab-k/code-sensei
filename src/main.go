package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

// ── types ─────────────────────────────────────────────────────────────────────

type IndexEntry struct {
	ID         string   `json:"id"`
	Contest    string   `json:"contest"`
	Problem    string   `json:"problem"`
	Title      string   `json:"title"`
	Difficulty int      `json:"difficulty"`
	Tags       []string `json:"tags"`
	File       string   `json:"file"`
}

type Sample struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type Solution struct {
	Code  string   `json:"code"`
	Steps []string `json:"steps"`
}

type Problem struct {
	ID          string               `json:"id"`
	Problem     string               `json:"problem"`
	Title       string               `json:"title"`
	AtcoderURL  string               `json:"atcoder_url"`
	Difficulty  int                  `json:"difficulty"`
	Tags        []string             `json:"tags"`
	Statement   string               `json:"statement"`
	Constraints string               `json:"constraints"`
	Samples     []Sample             `json:"samples"`
	Explanation string               `json:"explanation"`
	Solutions   map[string]*Solution `json:"solutions"`
	AddedAt     string               `json:"added_at"`
}

type ContestFile struct {
	Contest  string    `json:"contest"`
	Problems []Problem `json:"problems"`
}

// ── frontmatter parser ────────────────────────────────────────────────────────

type fmData struct {
	vals map[string]string
	tags []string
}

func parseFM(text string) (fmData, string) {
	d := fmData{vals: make(map[string]string)}
	if !strings.HasPrefix(text, "---\n") {
		return d, text
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return d, text
	}
	yaml := text[4 : end+4]
	body := text[end+9:]
	for _, line := range strings.Split(yaml, "\n") {
		i := strings.Index(line, ":")
		if i < 0 {
			continue
		}
		k, v := strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			for _, t := range strings.Split(v[1:len(v)-1], ",") {
				t = strings.Trim(strings.TrimSpace(t), `"'`)
				if t != "" {
					d.tags = append(d.tags, t)
				}
			}
		} else {
			d.vals[k] = strings.Trim(v, `"'`)
		}
	}
	return d, body
}

func (d fmData) str(k, def string) string {
	if v, ok := d.vals[k]; ok && v != "" {
		return v
	}
	return def
}

func (d fmData) num(k string, def int) int {
	var n int
	if v, ok := d.vals[k]; ok {
		fmt.Sscanf(v, "%d", &n)
	}
	if n == 0 {
		return def
	}
	return n
}

// ── markdown section helpers ──────────────────────────────────────────────────

var (
	reH1    = regexp.MustCompile(`(?m)^(# [^\n]+)`)
	reH2    = regexp.MustCompile(`(?m)^(## [^\n]+)`)
	reH3exp = regexp.MustCompile(`(?s)### 解説\n(.*?)(?:\n###|\n##|\n#|\z)`)
	reCode  = regexp.MustCompile("(?s)```[^\n]*\n(.*?)```")
	reHr    = regexp.MustCompile(`(?m)^---\s*$`)
	reOL    = regexp.MustCompile(`(?m)^\d+[.。]\s+(.+)$`)
	reH1top = regexp.MustCompile(`(?m)^# (.+)$`)
)

func splitByH1(text string) []struct{ title, body string } {
	locs := reH1.FindAllStringIndex(text, -1)
	var out []struct{ title, body string }
	for i, loc := range locs {
		title := strings.TrimPrefix(text[loc[0]:loc[1]], "# ")
		end := len(text)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		out = append(out, struct{ title, body string }{title, strings.TrimSpace(text[loc[1]:end])})
	}
	return out
}

func splitByH2(text string) map[string]string {
	locs := reH2.FindAllStringIndex(text, -1)
	out := make(map[string]string)
	for i, loc := range locs {
		title := strings.TrimPrefix(text[loc[0]:loc[1]], "## ")
		end := len(text)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		out[title] = strings.TrimSpace(text[loc[1]:end])
	}
	return out
}

func firstCode(text string) string {
	if m := reCode.FindStringSubmatch(text); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func orderedList(text string) []string {
	var out []string
	for _, m := range reOL.FindAllStringSubmatch(text, -1) {
		out = append(out, strings.TrimSpace(m[1]))
	}
	return out
}

func stripHr(s string) string { return strings.TrimSpace(reHr.ReplaceAllString(s, "")) }

// ── convert: md → JSON ───────────────────────────────────────────────────────

func convertFile(mdPath string) (Problem, string, string) {
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		log.Fatalf("read %s: %v", mdPath, err)
	}
	filename := strings.TrimSuffix(filepath.Base(mdPath), ".md")
	fm, body := parseFM(string(raw))

	id := fm.str("id", filename)
	fileKey := fm.str("file", filename[:max(strings.LastIndex(filename, "_"), 0)])

	h1s := splitByH1(body)
	h2 := splitByH2(body)

	// title / problem number from first h1
	title, prob := fm.str("title", ""), fm.str("problem", "")
	if len(h1s) > 0 {
		ft := h1s[0].title
		reParsed := regexp.MustCompile(`^(\S+)\s*[-–]\s*(.+?)(?:（.+?）)?$`)
		if m := reParsed.FindStringSubmatch(ft); m != nil {
			if prob == "" {
				prob = m[1]
			}
			if title == "" {
				title = strings.TrimSpace(m[2])
			}
		} else if title == "" {
			title = regexp.MustCompile(`（.+?）`).ReplaceAllString(ft, "")
			title = strings.TrimSpace(title)
		}
	}

	// samples
	var samples []Sample
	for i := 1; i <= 30; i++ {
		inKey := fmt.Sprintf("入力例%d", i)
		if _, ok := h2[inKey]; !ok {
			inKey = fmt.Sprintf("入力例 %d", i)
		}
		inText, ok := h2[inKey]
		if !ok {
			break
		}
		outKey := fmt.Sprintf("出力例%d", i)
		if _, ok := h2[outKey]; !ok {
			outKey = fmt.Sprintf("出力例 %d", i)
		}
		outText := h2[outKey]

		exp := ""
		if m := reH3exp.FindStringSubmatch(outText); m != nil {
			exp = strings.TrimSpace(m[1])
		} else {
			exp = reCode.ReplaceAllString(outText, "")
			exp = reHr.ReplaceAllString(exp, "")
			exp = reH1top.ReplaceAllString(exp, "")
			exp = strings.TrimSpace(exp)
		}

		samples = append(samples, Sample{firstCode(inText), firstCode(outText), exp})
	}

	// explanation: h1 sections from '解説' until code sections
	isCode := func(t string) bool {
		return strings.Contains(t, "Python") || strings.Contains(t, "C++") || strings.Contains(t, "C＋＋")
	}
	var expParts []string
	inExp := false
	for _, s := range h1s {
		if s.title == "解説" {
			inExp = true
		}
		if isCode(s.title) {
			inExp = false
		}
		if !inExp {
			continue
		}
		if s.title == "解説" {
			expParts = append(expParts, s.body)
		} else {
			expParts = append(expParts, "## "+s.title+"\n\n"+s.body)
		}
	}

	// solutions
	sols := make(map[string]*Solution)
	for _, s := range h1s {
		if strings.Contains(s.title, "Python") && strings.Contains(s.title, "解答") {
			sols["python"] = &Solution{firstCode(s.body), orderedList(s.body)}
		}
		if (strings.Contains(s.title, "C++") || strings.Contains(s.title, "C＋＋")) && strings.Contains(s.title, "解答") {
			sols["cpp"] = &Solution{firstCode(s.body), orderedList(s.body)}
		}
	}

	p := Problem{
		ID:          id,
		Problem:     prob,
		Title:       title,
		AtcoderURL:  fm.str("atcoder_url", ""),
		Difficulty:  fm.num("difficulty", 1),
		Tags:        fm.tags,
		Statement:   stripHr(h2["問題文"]),
		Constraints: stripHr(h2["制約"]),
		Samples:     samples,
		Explanation: strings.TrimSpace(strings.Join(expParts, "\n\n")),
		Solutions:   sols,
		AddedAt:     fm.str("added_at", time.Now().Format("2006-01-02")),
	}
	return p, fileKey, fm.str("contest", "")
}

func cmdConvert() {
	entries, _ := os.ReadDir("problems-raw")
	if err := os.MkdirAll("data/problems", 0755); err != nil {
		log.Fatal(err)
	}

	idxPath := "data/index.json"
	var index []IndexEntry
	if b, err := os.ReadFile(idxPath); err == nil {
		json.Unmarshal(b, &index)
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		fmt.Printf("変換中: %s ... ", e.Name())

		p, fileKey, contest := convertFile(filepath.Join("problems-raw", e.Name()))

		// update data/problems/<file>.json
		cPath := filepath.Join("data/problems", fileKey+".json")
		var cf ContestFile
		if b, err := os.ReadFile(cPath); err == nil {
			json.Unmarshal(b, &cf)
		} else {
			cf = ContestFile{Contest: contest}
		}
		found := false
		for i, ex := range cf.Problems {
			if ex.ID == p.ID {
				cf.Problems[i] = p
				found = true
				fmt.Println("更新")
				break
			}
		}
		if !found {
			cf.Problems = append(cf.Problems, p)
			fmt.Println("追加")
		}
		writeJSON(cPath, cf)

		// update index
		entry := IndexEntry{p.ID, contest, p.Problem, p.Title, p.Difficulty, p.Tags, fileKey}
		found = false
		for i, ex := range index {
			if ex.ID == p.ID {
				index[i] = entry
				found = true
				break
			}
		}
		if !found {
			index = append(index, entry)
		}
	}
	writeJSON(idxPath, index)
	fmt.Println("index.json 更新完了 ✓")
}

// ── build: JSON → HTML ────────────────────────────────────────────────────────

var diffLabel = map[int]string{1: "A問題レベル", 2: "B問題レベル", 3: "C問題レベル", 4: "D問題レベル"}
var diffColor = map[int]string{1: "#4caf50", 2: "#2196f3", 3: "#ff9800", 4: "#f44336"}

func e(s string) string { return html.EscapeString(s) }

func badge(d int) string {
	lbl := diffLabel[d]
	if lbl == "" {
		lbl = fmt.Sprintf("Lv%d", d)
	}
	col := diffColor[d]
	if col == "" {
		col = "#888"
	}
	return fmt.Sprintf(`<span class="diff-badge" style="background:%s">%s</span>`, col, lbl)
}

func tagSpans(tags []string) string {
	var b strings.Builder
	for _, t := range tags {
		fmt.Fprintf(&b, `<span class="tag">%s</span>`, e(t))
	}
	return b.String()
}

func mdToHTML(src string) string {
	if src == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(src), &buf); err != nil {
		return "<pre>" + e(src) + "</pre>"
	}
	return buf.String()
}

func css() string {
	return `*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#f5f6fa;color:#222;min-height:100vh}
a{color:inherit;text-decoration:none}
.header{background:#1565c0;color:#fff;padding:0 12px}
.header-inner{max-width:900px;margin:0 auto;display:flex;align-items:center;gap:12px;height:52px}
.header-logo{font-size:1.1rem;font-weight:700;flex:1}
.header-sub{font-size:.78rem;opacity:.8}
.back-btn{font-size:.85rem;color:#fff;opacity:.9;white-space:nowrap}
.back-btn:hover{opacity:1}
.list-view{max-width:900px;margin:20px auto;padding:0 12px}
.filters{display:flex;flex-direction:column;gap:10px;margin-bottom:18px}
.search-wrap{display:flex;align-items:center;border:1.5px solid #ccd;border-radius:8px;background:#fff;padding:0 10px}
.search-icon{color:#999;margin-right:6px}
.search-input{flex:1;border:none;outline:none;padding:10px 0;font-size:1rem;background:transparent}
.diff-filters{display:flex;flex-wrap:wrap;gap:6px}
.diff-btn{border:2px solid #ccc;background:#fff;border-radius:20px;padding:5px 14px;cursor:pointer;font-size:.82rem;font-weight:600;color:#555;transition:.15s}
.diff-btn.active{background:#1565c0;border-color:#1565c0;color:#fff}
.problem-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:14px}
.problem-card{display:block;background:#fff;border-radius:10px;padding:14px;box-shadow:0 1px 4px rgba(0,0,0,.08);border:1.5px solid transparent;transition:.15s}
.problem-card:hover{border-color:#1565c0;box-shadow:0 3px 12px rgba(21,101,192,.15)}
.card-top{display:flex;align-items:center;justify-content:space-between;margin-bottom:6px}
.card-contest{font-size:.75rem;color:#666}
.card-title{font-size:.97rem;font-weight:600;margin-bottom:8px;line-height:1.4}
.card-tags{display:flex;flex-wrap:wrap;gap:4px}
.tag{display:inline-block;background:#e8eaf6;color:#3949ab;border-radius:12px;padding:2px 9px;font-size:.72rem;font-weight:600}
.diff-badge{display:inline-block;border-radius:12px;padding:2px 10px;font-size:.72rem;font-weight:700;color:#fff}
.detail-view{max-width:820px;margin:20px auto;padding:0 12px}
.detail-contest{font-size:.8rem;color:#666;margin-bottom:4px}
.detail-title{font-size:1.4rem;font-weight:700;margin-bottom:10px;line-height:1.3}
.detail-meta{display:flex;align-items:center;flex-wrap:wrap;gap:8px;margin-bottom:20px}
.atcoder-link{font-size:.82rem;color:#1565c0;border:1px solid #1565c0;border-radius:6px;padding:3px 10px}
.atcoder-link:hover{background:#e3f2fd}
.detail-section{margin-bottom:28px}
.section-title{font-size:1rem;font-weight:700;margin-bottom:10px;color:#333}
.statement-box{background:#fff;border:1.5px solid #e0e0e0;border-radius:8px;padding:14px;line-height:1.75;white-space:pre-wrap;font-size:.93rem}
.constraints-box{background:#f3f4f6;border-radius:6px;padding:10px 14px;font-size:.85rem;margin-top:10px;line-height:1.7;white-space:pre-wrap}
.sample-block{background:#fff;border:1.5px solid #e0e0e0;border-radius:8px;padding:12px;margin-bottom:12px}
.sample-row{display:grid;grid-template-columns:1fr 1fr;gap:12px}
.sample-label{font-size:.75rem;font-weight:600;color:#555;margin-bottom:4px}
.sample-pre{background:#f5f6fa;border-radius:6px;padding:8px 10px;font-size:.85rem;font-family:monospace;overflow-x:auto;white-space:pre}
.sample-explanation{margin-top:10px;font-size:.85rem;color:#555;line-height:1.6;border-top:1px solid #eee;padding-top:8px}
.explanation-box{background:#fffde7;border-left:4px solid #f9a825;border-radius:6px;padding:14px 16px;line-height:1.75;font-size:.93rem}
.explanation-box h2{font-size:1rem;margin:14px 0 6px}
.explanation-box h3{font-size:.9rem;margin:10px 0 4px}
.explanation-box p{margin-bottom:8px}
.explanation-box pre{background:#fff8e1;border-radius:4px;padding:8px;font-size:.82rem;overflow-x:auto;margin:8px 0}
.explanation-box code{background:#fff8e1;padding:1px 5px;border-radius:3px;font-size:.85em}
.explanation-box ul,.explanation-box ol{padding-left:1.4em;margin-bottom:8px}
.lang-tabs{display:flex;gap:4px;margin-bottom:8px}
.lang-tab{border:1.5px solid #ccc;background:#fff;border-radius:8px 8px 0 0;padding:6px 18px;cursor:pointer;font-size:.88rem;font-weight:600;color:#555}
.lang-tab.active{background:#1565c0;border-color:#1565c0;color:#fff}
.code-panel{display:none}
.code-panel.active{display:block}
.code-block{background:#1e1e2e;color:#cdd6f4;border-radius:0 8px 8px 8px;padding:14px;font-size:.82rem;font-family:monospace;overflow-x:auto;white-space:pre;line-height:1.6}
.solution-steps{padding-left:1.5em;margin-top:12px;font-size:.88rem;line-height:1.7;color:#444}
.empty{text-align:center;color:#999;padding:40px;font-size:.95rem}
@media(max-width:480px){.sample-row{grid-template-columns:1fr}.problem-grid{grid-template-columns:1fr}.detail-title{font-size:1.15rem}}`
}

// root: "" for index.html, "../" for problems/*.html
func shell(title, root, back, backLabel, sub, body string) string {
	backEl := ""
	if back != "" {
		backEl = fmt.Sprintf(`<a class="back-btn" href="%s">&#8592; %s</a>`, back, backLabel)
	}
	subEl := ""
	if sub != "" {
		subEl = fmt.Sprintf(`<span class="header-sub">%s</span>`, sub)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s | 競プロ教材</title>
<link rel="stylesheet" href="%sstyle.css">
</head>
<body>
<header class="header">
  <div class="header-inner">
    %s
    <a class="header-logo" href="%s">&#127891; 競プロ教材</a>
    %s
  </div>
</header>
%s
</body>
</html>`, e(title), root, backEl, root, subEl, body)
}

func buildIndex(index []IndexEntry) {
	diffs := make(map[int]bool)
	for _, p := range index {
		diffs[p.Difficulty] = true
	}
	var btnBuf strings.Builder
	btnBuf.WriteString(`<button class="diff-btn active" data-diff="">すべて</button>`)
	for d := 1; d <= 4; d++ {
		if !diffs[d] {
			continue
		}
		fmt.Fprintf(&btnBuf, `<button class="diff-btn" data-diff="%d" style="--dc:%s">%s</button>`,
			d, diffColor[d], diffLabel[d])
	}

	var cardBuf strings.Builder
	for _, p := range index {
		fmt.Fprintf(&cardBuf, `
<a class="problem-card" href="problems/%s.html"
   data-diff="%d"
   data-title="%s"
   data-tags="%s"
   data-contest="%s">
  <div class="card-top">
    <span class="card-contest">%s %s</span>
    %s
  </div>
  <div class="card-title">%s</div>
  <div class="card-tags">%s</div>
</a>`,
			p.ID, p.Difficulty,
			e(strings.ToLower(p.Title)),
			e(strings.ToLower(strings.Join(p.Tags, " "))),
			e(strings.ToLower(p.Contest)),
			e(p.Contest), e(p.Problem),
			badge(p.Difficulty),
			e(p.Title),
			tagSpans(p.Tags),
		)
	}

	body := fmt.Sprintf(`
<main class="list-view">
  <div class="filters">
    <div class="search-wrap">
      <span class="search-icon">&#128269;</span>
      <input class="search-input" id="search" placeholder="タイトル・タグで検索..." autocomplete="off">
    </div>
    <div class="diff-filters" id="dfilters">%s</div>
  </div>
  <div class="problem-grid" id="grid">%s</div>
  <div class="empty" id="empty" hidden>問題が見つかりませんでした</div>
</main>
<script>
(function(){
  var grid=document.getElementById('grid'),search=document.getElementById('search'),
      empty=document.getElementById('empty'),cards=Array.from(grid.querySelectorAll('.problem-card')),curD='';
  function filter(){
    var q=search.value.toLowerCase(),v=0;
    cards.forEach(function(c){
      var ok=(!curD||c.dataset.diff===curD)&&(!q||c.dataset.title.includes(q)||c.dataset.tags.includes(q)||c.dataset.contest.includes(q));
      c.hidden=!ok; if(ok)v++;
    });
    empty.hidden=v>0;
  }
  search.addEventListener('input',filter);
  document.getElementById('dfilters').addEventListener('click',function(e){
    var b=e.target.closest('.diff-btn'); if(!b)return;
    curD=b.dataset.diff;
    document.querySelectorAll('.diff-btn').forEach(function(x){x.classList.remove('active');});
    b.classList.add('active'); filter();
  });
})();
</script>`, btnBuf.String(), cardBuf.String())

	writeFile("docs/index.html", shell("問題一覧", "", "", "", "小学生向け AtCoder 解説集", body))
	fmt.Println("  docs/index.html")
}

func buildProblem(p Problem, contest string, force bool) {
	out := "docs/problems/" + p.ID + ".html"
	if !force {
		if _, err := os.Stat(out); err == nil {
			fmt.Printf("  %s (スキップ)\n", out)
			return
		}
	}
	var sampBuf strings.Builder
	for i, s := range p.Samples {
		exp := ""
		if s.Explanation != "" {
			exp = fmt.Sprintf(`<div class="sample-explanation">%s</div>`, e(s.Explanation))
		}
		fmt.Fprintf(&sampBuf, `
<div class="sample-block">
  <div class="sample-row">
    <div class="sample-col"><div class="sample-label">入力 %d</div><pre class="sample-pre">%s</pre></div>
    <div class="sample-col"><div class="sample-label">出力 %d</div><pre class="sample-pre">%s</pre></div>
  </div>%s
</div>`, i+1, e(s.Input), i+1, e(s.Output), exp)
	}

	atcLink := ""
	if p.AtcoderURL != "" {
		atcLink = fmt.Sprintf(`<a class="atcoder-link" href="%s" target="_blank" rel="noopener">&#128279; AtCoderで見る</a>`, e(p.AtcoderURL))
	}

	expSection := ""
	if p.Explanation != "" {
		expSection = fmt.Sprintf(`
<section class="detail-section">
  <h2 class="section-title">&#128161; 解説</h2>
  <div class="explanation-box">%s</div>
</section>`, mdToHTML(p.Explanation))
	}

	codeSection := buildCodeSection(p)

	body := fmt.Sprintf(`
<main class="detail-view">
  <div class="detail-contest">%s %s</div>
  <h1 class="detail-title">%s</h1>
  <div class="detail-meta">%s %s %s</div>
  <section class="detail-section">
    <h2 class="section-title">&#128196; 問題文</h2>
    <div class="statement-box">%s</div>
    %s
  </section>
  <section class="detail-section">
    <h2 class="section-title">&#10067; 入力・出力の例</h2>
    %s
  </section>
  %s
  %s
</main>`,
		e(contest), e(p.Problem), e(p.Title),
		badge(p.Difficulty), tagSpans(p.Tags), atcLink,
		e(p.Statement),
		constraintsBlock(p.Constraints),
		sampBuf.String(),
		expSection,
		codeSection,
	)

	writeFile(out, shell(p.Title, "../", "../", "もどる", "", body))
	fmt.Printf("  %s\n", out)
}

func constraintsBlock(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="constraints-box"><strong>制約</strong><br>%s</div>`, e(s))
}

func buildCodeSection(p Problem) string {
	py, cpp := p.Solutions["python"], p.Solutions["cpp"]
	if py == nil && cpp == nil {
		return ""
	}
	var tabs, panels strings.Builder
	if py != nil {
		tabs.WriteString(`<button class="lang-tab active" data-lang="python">Python</button>`)
		panels.WriteString(codePanel("python", true, py))
	}
	if cpp != nil {
		tabs.WriteString(`<button class="lang-tab" data-lang="cpp">C++</button>`)
		panels.WriteString(codePanel("cpp", py == nil, cpp))
	}
	return fmt.Sprintf(`
<section class="detail-section">
  <h2 class="section-title">&#128187; 解答コード</h2>
  <div class="lang-tabs" id="ltabs">%s</div>
  %s
</section>
<script>
document.getElementById('ltabs').addEventListener('click',function(e){
  var b=e.target.closest('.lang-tab'); if(!b)return;
  document.querySelectorAll('.lang-tab').forEach(function(x){x.classList.remove('active');});
  document.querySelectorAll('.code-panel').forEach(function(x){x.classList.remove('active');});
  b.classList.add('active');
  var p=document.getElementById('panel-'+b.dataset.lang); if(p)p.classList.add('active');
});
</script>`, tabs.String(), panels.String())
}

func codePanel(lang string, active bool, sol *Solution) string {
	cls := "code-panel"
	if active {
		cls += " active"
	}
	steps := ""
	if len(sol.Steps) > 0 {
		var sb strings.Builder
		for _, s := range sol.Steps {
			fmt.Fprintf(&sb, "<li>%s</li>", e(s))
		}
		steps = fmt.Sprintf(`<ol class="solution-steps">%s</ol>`, sb.String())
	}
	return fmt.Sprintf(`<div class="%s" id="panel-%s"><pre class="code-block">%s</pre>%s</div>`,
		cls, lang, e(sol.Code), steps)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func writeJSON(path string, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	os.WriteFile(path, append(b, '\n'), 0644)
}

func writeFile(path, content string) {
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(content), 0644)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cmdBuild(force bool) {
	if force {
		os.RemoveAll("docs")
	}
	os.MkdirAll("docs/problems", 0755)

	idxB, err := os.ReadFile("data/index.json")
	if err != nil {
		log.Fatal("data/index.json が見つかりません。先に convert を実行してください。")
	}
	var index []IndexEntry
	json.Unmarshal(idxB, &index)

	writeFile("docs/style.css", css())
	fmt.Println("  docs/style.css")
	fmt.Println("HTML 生成中...")
	buildIndex(index)

	seen := make(map[string]bool)
	for _, meta := range index {
		if seen[meta.File] {
			continue
		}
		seen[meta.File] = true
		var cf ContestFile
		b, err := os.ReadFile("data/problems/" + meta.File + ".json")
		if err != nil {
			continue
		}
		json.Unmarshal(b, &cf)
		for _, p := range cf.Problems {
			buildProblem(p, cf.Contest, force)
		}
	}
	fmt.Println("\n完了 ✓  →  docs/")
}

func cmdServe() {
	fmt.Println("http://localhost:8080 で起動中...")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("docs"))))
}

// ── entry point ───────────────────────────────────────────────────────────────

func main() {
	// src/ 内から実行された場合はプロジェクトルートに移動
	if _, err := os.Stat("problems-raw"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	cmd := "all"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	force := false
	for _, a := range os.Args[2:] {
		if a == "-f" || a == "--force" {
			force = true
		}
	}
	switch cmd {
	case "convert":
		cmdConvert()
	case "build":
		cmdBuild(force)
	case "serve":
		cmdServe()
	default: // "all"
		cmdConvert()
		cmdBuild(force)
	}
}
