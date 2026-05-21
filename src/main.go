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
	Samples          []Sample             `json:"samples"`
	StatementNote    string               `json:"statement_note"`
	EasyExplanation  string               `json:"easy_explanation"`
	Explanation      string               `json:"explanation"`
	Solutions        map[string]*Solution `json:"solutions"`
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
	reFence = regexp.MustCompile("(?m)^```")
	reHr    = regexp.MustCompile(`(?m)^---\s*$`)
	reOL    = regexp.MustCompile(`(?m)^\d+[.。]\s+(.+)$`)
	reH1top = regexp.MustCompile(`(?m)^# (.+)$`)
)

// コードフェンス（``` ... ```）の範囲を返す。この範囲内の # はセクション区切りではない
func fenceRanges(text string) [][2]int {
	locs := reFence.FindAllStringIndex(text, -1)
	var ranges [][2]int
	for i := 0; i+1 < len(locs); i += 2 {
		ranges = append(ranges, [2]int{locs[i][0], locs[i+1][1]})
	}
	return ranges
}

func inFence(pos int, ranges [][2]int) bool {
	for _, r := range ranges {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

func splitByH1(text string) []struct{ title, body string } {
	fences := fenceRanges(text)
	allLocs := reH1.FindAllStringIndex(text, -1)
	var locs [][]int
	for _, loc := range allLocs {
		if !inFence(loc[0], fences) {
			locs = append(locs, loc)
		}
	}
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

	// explanations: collect easy and standard separately, stop at code sections
	isCode := func(t string) bool {
		for _, lang := range langDefs {
			for _, kw := range lang.matches {
				if strings.Contains(t, kw) && strings.Contains(t, "解答") {
					return true
				}
			}
		}
		return false
	}
	var stmtNoteParts, easyParts, expParts []string
	mode := "" // "stmt" | "easy" | "standard" | ""
	for _, s := range h1s {
		if isCode(s.title) {
			mode = ""
			continue
		}
		if s.title == "問題の解説" {
			mode = "stmt"
			if s.body != "" {
				stmtNoteParts = append(stmtNoteParts, s.body)
			}
			continue
		}
		if s.title == "やさしい解説" {
			mode = "easy"
			if s.body != "" {
				easyParts = append(easyParts, s.body)
			}
			continue
		}
		if s.title == "解説" {
			mode = "standard"
			if s.body != "" {
				expParts = append(expParts, s.body)
			}
			continue
		}
		switch mode {
		case "stmt":
			stmtNoteParts = append(stmtNoteParts, "## "+s.title+"\n\n"+s.body)
		case "easy":
			easyParts = append(easyParts, "## "+s.title+"\n\n"+s.body)
		case "standard":
			expParts = append(expParts, "## "+s.title+"\n\n"+s.body)
		}
	}

	// solutions
	sols := make(map[string]*Solution)
	for _, s := range h1s {
		if !strings.Contains(s.title, "解答") {
			continue
		}
		for _, lang := range langDefs {
			for _, kw := range lang.matches {
				if strings.Contains(s.title, kw) {
					sols[lang.key] = &Solution{firstCode(s.body), orderedList(s.body)}
					break
				}
			}
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
		Samples:         samples,
		StatementNote:   strings.TrimSpace(strings.Join(stmtNoteParts, "\n\n")),
		EasyExplanation: strings.TrimSpace(strings.Join(easyParts, "\n\n")),
		Explanation:     strings.TrimSpace(strings.Join(expParts, "\n\n")),
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

type langDef struct {
	key     string
	label   string
	matches []string
}

var langDefs = []langDef{
	{"python", "Python", []string{"Python"}},
	{"cpp", "C++", []string{"C++", "C＋＋"}},
	{"typescript", "TypeScript", []string{"TypeScript", "TS"}},
	{"ruby", "Ruby", []string{"Ruby"}},
	{"php", "PHP", []string{"PHP"}},
	{"rust", "Rust", []string{"Rust"}},
	{"perl", "Perl", []string{"Perl"}},
}

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
.statement-box{background:#fff;border:1.5px solid #e0e0e0;border-radius:8px;padding:14px;line-height:1.75;font-size:.93rem}
.statement-box p{margin-bottom:.6em}
.statement-box ul,.statement-box ol{padding-left:1.6em;margin:.4em 0 .6em}
.statement-box li{margin-bottom:.2em}
.statement-box code{background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:.88em}
.constraints-box{background:#f3f4f6;border-radius:6px;padding:10px 14px;font-size:.85rem;margin-top:10px;line-height:1.7}
.statement-note{background:#e3f2fd;border-left:4px solid #1565c0;border-radius:6px;padding:12px 14px;font-size:.88rem;margin-top:12px;line-height:1.75}
.statement-note p{margin-bottom:.5em}
.statement-note strong{color:#1565c0}
.statement-note ul,.statement-note ol{padding-left:1.6em;margin:.4em 0 .6em}
.statement-note hr{display:none}
.statement-note-badge{display:inline-block;background:#1565c0;color:#fff;font-size:.72rem;font-weight:700;border-radius:4px;padding:1px 7px;margin-bottom:8px;letter-spacing:.03em}
.constraints-box ul,.constraints-box ol{padding-left:1.4em}
.constraints-box code{background:#e8e8e8;padding:1px 4px;border-radius:3px;font-size:.85em}
.sample-block{background:#fff;border:1.5px solid #e0e0e0;border-radius:8px;padding:12px;margin-bottom:12px}
.sample-row{display:flex;flex-direction:column;gap:10px}
.sample-label{font-size:.75rem;font-weight:600;color:#555;margin-bottom:4px}
.sample-pre{background:#f5f6fa;border-radius:6px;padding:8px 10px;font-size:.85rem;font-family:monospace;overflow-x:auto;white-space:pre;word-break:break-all}
.sample-explanation{margin-top:10px;font-size:.85rem;color:#444;line-height:1.7;border-top:1px solid #e0e0e0;padding-top:10px}
.sample-explanation::before{content:"📖 解説";display:block;font-weight:700;font-size:.78rem;color:#1565c0;margin-bottom:6px}
.sample-explanation p{margin-bottom:.5em}
.sample-explanation ul,.sample-explanation ol{padding-left:1.5em;margin:.3em 0 .5em}
.sample-explanation hr{display:none}
.sample-explanation code{background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:.88em}
.explanation-box{background:#e8f5e9;border-left:4px solid #43a047;border-radius:6px;padding:14px 16px;line-height:1.75;font-size:.93rem}
.easy-box{background:#fffde7;border-left:4px solid #f9a825}
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
pre.code-block{background:#1e1e2e;color:#cdd6f4;border-radius:0 8px 8px 8px;padding:14px;font-size:.82rem;font-family:monospace;overflow-x:auto;white-space:pre;line-height:1.6}
pre.code-block code{background:none;padding:0;border-radius:0;font-size:inherit;font-family:inherit}
.solution-steps{padding-left:1.5em;margin-top:12px;font-size:.88rem;line-height:1.7;color:#444}
.empty{text-align:center;color:#999;padding:40px;font-size:.95rem}
@media(max-width:600px){.problem-grid{grid-template-columns:1fr}.detail-title{font-size:1.15rem}}`
}

func pwGateScript() string {
	return `<script>(function(){` +
		`var H="__VIEW_PASSWORD__";` +
		`if(H.indexOf('__')===0)return;` +
		`if(sessionStorage.getItem('cps')===H)return;` +
		`var o=document.createElement('div');` +
		`o.style='position:fixed;inset:0;background:#1565c0;display:flex;align-items:center;justify-content:center;z-index:9999';` +
		`o.innerHTML='<div style="background:#fff;border-radius:12px;padding:32px;max-width:320px;width:90%;text-align:center">` +
		`<div style="font-size:1.4rem;font-weight:700;margin-bottom:8px">&#127891; 競プロ教材</div>` +
		`<div style="color:#666;font-size:.85rem;margin-bottom:20px">パスワードを入力してください</div>` +
		`<input type="password" id="pwi" style="width:100%;border:1.5px solid #ccc;border-radius:8px;padding:10px;font-size:1rem;margin-bottom:10px;box-sizing:border-box" placeholder="パスワード">` +
		`<div id="pwe" style="color:#f44336;font-size:.8rem;height:1.2em;margin-bottom:8px"></div>` +
		`<button id="pwb" style="background:#1565c0;color:#fff;border:none;border-radius:8px;padding:10px 32px;font-size:1rem;cursor:pointer;width:100%">入る</button>` +
		`</div>';` +
		`document.body.appendChild(o);` +
		`function chk(){` +
		`var v=document.getElementById('pwi').value;` +
		`if(v===H){sessionStorage.setItem('cps',H);o.remove();}` +
		`else document.getElementById('pwe').textContent='パスワードが違います';}` +
		`document.getElementById('pwb').addEventListener('click',chk);` +
		`document.getElementById('pwi').addEventListener('keydown',function(e){if(e.key==='Enter')chk();});` +
		`})();</script>`
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
<meta name="robots" content="noindex, nofollow">
<title>%s | 競プロ教材</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/prismjs@1/themes/prism-tomorrow.min.css">
<link rel="stylesheet" href="%sstyle.css">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js"
  onload="renderMathInElement(document.body,{delimiters:[{left:'$$',right:'$$',display:true},{left:'$',right:'$',display:false}],ignoredTags:['script','noscript','style','pre','code']})"></script>
</head>
<body>
%s
<header class="header">
  <div class="header-inner">
    %s
    <a class="header-logo" href="%s">&#127891; 競プロ教材</a>
    %s
  </div>
</header>
%s
<script src="https://cdn.jsdelivr.net/npm/prismjs@1/components/prism-core.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/prismjs@1/plugins/autoloader/prism-autoloader.min.js"></script>
</body>
</html>`, e(title), root, pwGateScript(), backEl, root, subEl, body)
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
			exp = fmt.Sprintf(`<div class="sample-explanation">%s</div>`, mdToHTML(s.Explanation))
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

	easySection := ""
	if p.EasyExplanation != "" {
		easySection = fmt.Sprintf(`
<section class="detail-section">
  <h2 class="section-title">&#128640; わかりやすく解説</h2>
  <div class="explanation-box easy-box">%s</div>
</section>`, mdToHTML(p.EasyExplanation))
	}

	expSection := ""
	if p.Explanation != "" {
		expSection = fmt.Sprintf(`
<section class="detail-section">
  <h2 class="section-title">&#128218; くわしい解説</h2>
  <div class="explanation-box">%s</div>
</section>`, mdToHTML(p.Explanation))
	}

	stmtNoteSection := ""
	if p.StatementNote != "" {
		stmtNoteSection = fmt.Sprintf(
			`<div class="statement-note"><span class="statement-note-badge">&#128221; 入力例1を使った説明</span>%s</div>`,
			mdToHTML(p.StatementNote))
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
    %s
  </section>
  <section class="detail-section">
    <h2 class="section-title">&#10067; 入力・出力の例</h2>
    %s
  </section>
  %s
  %s
  %s
</main>`,
		e(contest), e(p.Problem), e(p.Title),
		badge(p.Difficulty), tagSpans(p.Tags), atcLink,
		mdToHTML(p.Statement),
		constraintsBlock(p.Constraints),
		stmtNoteSection,
		sampBuf.String(),
		easySection,
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
	return fmt.Sprintf(`<div class="constraints-box"><strong>制約</strong>%s</div>`, mdToHTML(s))
}

func printShell(title, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="noindex, nofollow">
<title>%s | 印刷用 | 競プロ教材</title>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js"
  onload="renderMathInElement(document.body,{delimiters:[{left:'$$',right:'$$',display:true},{left:'$',right:'$',display:false}],ignoredTags:['script','noscript','style','pre','code']})"></script>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#fff;color:#111;max-width:820px;margin:0 auto;padding:24px 32px;font-size:.93rem;line-height:1.7}
h1{font-size:1.3rem;margin-bottom:8px}
h2{font-size:1rem;font-weight:700;margin:20px 0 8px;color:#333;border-bottom:1px solid #e0e0e0;padding-bottom:4px}
h3{font-size:.9rem;font-weight:700;margin:14px 0 6px}
p{margin-bottom:.6em}
ul,ol{padding-left:1.6em;margin:.4em 0 .6em}
li{margin-bottom:.2em}
code{background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:.88em;font-family:monospace}
.diff-badge{display:inline-block;border-radius:12px;padding:2px 10px;font-size:.72rem;font-weight:700;color:#fff}
.tag{display:inline-block;background:#e8eaf6;color:#3949ab;border-radius:12px;padding:2px 9px;font-size:.72rem;font-weight:600}
.print-header{border-bottom:2px solid #1565c0;padding-bottom:12px;margin-bottom:20px}
.print-contest{font-size:.8rem;color:#666;margin-bottom:4px}
.print-meta{display:flex;flex-wrap:wrap;gap:8px;margin-top:8px;align-items:center}
.print-section{margin-bottom:22px}
.print-box{border:1px solid #e0e0e0;border-radius:6px;padding:12px 14px}
.print-box p{margin-bottom:.6em}
.print-box ul,.print-box ol{padding-left:1.6em;margin:.4em 0 .6em}
.print-constraints{background:#f3f4f6;border-radius:6px;padding:8px 12px;font-size:.85rem;margin-top:8px}
.print-note{background:#e3f2fd;border-left:4px solid #1565c0;border-radius:4px;padding:10px 14px}
.print-note h2,.print-note h3{color:#1565c0}
.print-note hr{display:none}
.print-easy{background:#fffde7;border-left:4px solid #f9a825;border-radius:4px;padding:10px 14px}
.print-easy hr{display:none}
.print-exp{background:#e8f5e9;border-left:4px solid #43a047;border-radius:4px;padding:10px 14px}
.print-exp hr{display:none}
.print-note p,.print-easy p,.print-exp p{margin-bottom:.5em}
.print-note ul,.print-note ol,.print-easy ul,.print-easy ol,.print-exp ul,.print-exp ol{padding-left:1.4em;margin:.3em 0 .5em}
.print-note code,.print-easy code,.print-exp code{background:rgba(0,0,0,.06);padding:1px 5px;border-radius:3px;font-size:.88em}
.print-note pre,.print-easy pre,.print-exp pre{background:rgba(0,0,0,.04);border-radius:4px;padding:8px;font-size:.82rem;overflow-x:auto;margin:6px 0;white-space:pre-wrap}
.sample-block{border:1px solid #ddd;border-radius:6px;padding:10px;margin-bottom:10px}
.sample-label{font-size:.72rem;font-weight:600;color:#555;margin-bottom:3px}
.sample-pre{background:#f5f6fa;border-radius:4px;padding:6px 10px;font-size:.82rem;font-family:monospace;white-space:pre;margin-bottom:8px}
.sample-exp{font-size:.82rem;color:#444;border-top:1px solid #e0e0e0;padding-top:8px;margin-top:8px}
.sample-exp p{margin-bottom:.4em}
.sample-exp ul,.sample-exp ol{padding-left:1.4em;margin:.2em 0 .4em}
.lang-header{background:#1565c0;color:#fff;padding:4px 12px;border-radius:4px 4px 0 0;font-size:.8rem;font-weight:700;margin-top:14px}
pre.code-print{background:#f5f5f5;border:1px solid #ddd;border-radius:0 4px 4px 4px;padding:12px;font-size:.75rem;font-family:monospace;white-space:pre-wrap;word-break:break-all;line-height:1.55;margin-bottom:0}
pre.code-print code{background:none;padding:0;font-size:inherit}
.print-btn{position:fixed;bottom:20px;right:20px;background:#1565c0;color:#fff;border:none;border-radius:8px;padding:10px 20px;font-size:.9rem;cursor:pointer;font-weight:600;box-shadow:0 2px 8px rgba(0,0,0,.2)}
.print-btn:hover{background:#0d47a1}
@media print{
  @page{margin:1.5cm}
  body{max-width:none;padding:0}
  .print-btn{display:none}
  pre.code-print{page-break-inside:avoid;white-space:pre-wrap}
  h2{page-break-after:avoid}
  .print-note,.print-easy,.print-exp,.sample-block{page-break-inside:avoid}
  .lang-header{page-break-after:avoid}
}
</style>
</head>
<body>
%s
%s
<button class="print-btn" onclick="window.print()">&#128424; 印刷</button>
</body>
</html>`, e(title), pwGateScript(), body)
}

func buildPrintPage(p Problem, contest string, force bool) {
	out := "docs/print/" + p.ID + ".html"
	if !force {
		if _, err := os.Stat(out); err == nil {
			fmt.Printf("  %s (スキップ)\n", out)
			return
		}
	}
	var buf strings.Builder

	fmt.Fprintf(&buf, `<div class="print-header">
<div class="print-contest">%s %s</div>
<h1>%s</h1>
<div class="print-meta">%s %s</div>
</div>`, e(contest), e(p.Problem), e(p.Title), badge(p.Difficulty), tagSpans(p.Tags))

	fmt.Fprintf(&buf, `<section class="print-section"><h2>&#128196; 問題文</h2><div class="print-box">%s</div>`, mdToHTML(p.Statement))
	if p.Constraints != "" {
		fmt.Fprintf(&buf, `<div class="print-constraints"><strong>制約</strong>%s</div>`, mdToHTML(p.Constraints))
	}
	buf.WriteString(`</section>`)

	if len(p.Samples) > 0 {
		buf.WriteString(`<section class="print-section"><h2>&#10067; 入力・出力の例</h2>`)
		for i, s := range p.Samples {
			exp := ""
			if s.Explanation != "" {
				exp = fmt.Sprintf(`<div class="sample-exp">%s</div>`, mdToHTML(s.Explanation))
			}
			fmt.Fprintf(&buf, `<div class="sample-block"><div class="sample-label">入力 %d</div><pre class="sample-pre">%s</pre><div class="sample-label">出力 %d</div><pre class="sample-pre">%s</pre>%s</div>`,
				i+1, e(s.Input), i+1, e(s.Output), exp)
		}
		buf.WriteString(`</section>`)
	}

	if p.StatementNote != "" {
		fmt.Fprintf(&buf, `<section class="print-section"><h2>&#128221; 問題の解説（入力例1で説明）</h2><div class="print-note">%s</div></section>`, mdToHTML(p.StatementNote))
	}
	if p.EasyExplanation != "" {
		fmt.Fprintf(&buf, `<section class="print-section"><h2>&#128640; わかりやすく解説</h2><div class="print-easy">%s</div></section>`, mdToHTML(p.EasyExplanation))
	}
	if p.Explanation != "" {
		fmt.Fprintf(&buf, `<section class="print-section"><h2>&#128218; くわしい解説</h2><div class="print-exp">%s</div></section>`, mdToHTML(p.Explanation))
	}

	hasSols := false
	for _, lang := range langDefs {
		if sol, ok := p.Solutions[lang.key]; ok && sol != nil {
			hasSols = true
			_ = sol
			break
		}
	}
	if hasSols {
		buf.WriteString(`<section class="print-section"><h2>&#128187; 解答コード</h2>`)
		for _, lang := range langDefs {
			sol, ok := p.Solutions[lang.key]
			if !ok || sol == nil {
				continue
			}
			fmt.Fprintf(&buf, `<div class="lang-header">%s</div><pre class="code-print"><code>%s</code></pre>`,
				lang.label, e(sol.Code))
		}
		buf.WriteString(`</section>`)
	}

	writeFile(out, printShell(p.Title, buf.String()))
	fmt.Printf("  %s\n", out)
}

func buildCodeSection(p Problem) string {
	if len(p.Solutions) == 0 {
		return ""
	}
	var tabs, panels strings.Builder
	first := true
	for _, lang := range langDefs {
		sol, ok := p.Solutions[lang.key]
		if !ok || sol == nil {
			continue
		}
		active := first
		first = false
		if active {
			fmt.Fprintf(&tabs, `<button class="lang-tab active" data-lang="%s">%s</button>`, lang.key, lang.label)
		} else {
			fmt.Fprintf(&tabs, `<button class="lang-tab" data-lang="%s">%s</button>`, lang.key, lang.label)
		}
		panels.WriteString(codePanel(lang.key, active, sol))
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
	return fmt.Sprintf(`<div class="%s" id="panel-%s"><pre class="code-block language-%s"><code class="language-%s">%s</code></pre>%s</div>`,
		cls, lang, lang, lang, e(sol.Code), steps)
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
	os.MkdirAll("docs/print", 0755)

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
			buildPrintPage(p, cf.Contest, force)
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
	force := false
	for _, a := range os.Args[1:] {
		switch a {
		case "-f", "--force":
			force = true
		case "convert", "build", "serve":
			cmd = a
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
