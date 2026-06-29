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
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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
	Statement       string               `json:"statement"`
	Constraints     string               `json:"constraints"`
	ConstraintsNote string               `json:"constraints_note"`
	Samples         []Sample             `json:"samples"`
	StatementNote    string               `json:"statement_note"`
	EasyExplanation  string               `json:"easy_explanation"`
	Explanation      string               `json:"explanation"`
	Solutions        map[string]*Solution    `json:"solutions"`
	BadSolutions     map[string]*BadSolution `json:"bad_solutions"`
	AddedAt          string                  `json:"added_at"`
}

type ContestFile struct {
	Contest  string    `json:"contest"`
	Problems []Problem `json:"problems"`
}

type BadSolution struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type GlossaryEntry struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Short        string   `json:"short"`
	Description  string   `json:"description"`
	WithoutLabel string   `json:"without_label"`
	WithoutCode  string   `json:"without_code"`
	WithLabel    string   `json:"with_label"`
	WithCode     string   `json:"with_code"`
	WhenToUse    string   `json:"when_to_use"`
	Problems     []string `json:"problems"`
}

type CodeReadingEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Short      string `json:"short"`
	Body       string `json:"body"`
	PythonCode string `json:"python_code"`
	OtherNote  string `json:"other_note"`
}

// ── build: JSON → HTML ────────────────────────────────────────────────────────

var diffColor = map[int]string{
	1: "#78909c", 2: "#4caf50", 3: "#2196f3", 4: "#ff9800",
	5: "#ff5722", 6: "#e53935", 7: "#7b1fa2",
}

func starLabel(d int) string { return fmt.Sprintf("★%d", d) }

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
	col := diffColor[d]
	if col == "" {
		col = "#888"
	}
	return fmt.Sprintf(`<span class="diff-badge" style="background:%s">%s</span>`, col, starLabel(d))
}

func tagSpans(tags []string) string {
	var b strings.Builder
	for _, t := range tags {
		fmt.Fprintf(&b, `<span class="tag">%s</span>`, e(t))
	}
	return b.String()
}

func icon(name string) string {
	return `<span class="mi">` + name + `</span>`
}

var md = goldmark.New(goldmark.WithExtensions(extension.Table, extension.Strikethrough))

func mdToHTML(src string) string {
	if src == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "<pre>" + e(src) + "</pre>"
	}
	return buf.String()
}

func css() string {
	return `*{box-sizing:border-box;margin:0;padding:0}
[hidden]{display:none!important}
body{font-family:'DotGothic16',monospace;background:#000;color:#fff;min-height:100vh;line-height:2}
body::after{content:'';position:fixed;inset:0;pointer-events:none;z-index:9998;background:repeating-linear-gradient(0deg,transparent,transparent 3px,rgba(0,0,0,.06) 3px,rgba(0,0,0,.06) 4px)}
a{color:inherit;text-decoration:none}
.mi{font-family:'Material Symbols Outlined';font-weight:normal;font-style:normal;font-size:1.25em;display:inline-block;vertical-align:-.22em;line-height:1;letter-spacing:normal;text-transform:none;white-space:nowrap;font-feature-settings:'liga';-webkit-font-smoothing:antialiased}
.header{background:#050505;border-bottom:2px solid #444;padding:0 16px}
.header-inner{max-width:960px;margin:0 auto;display:flex;align-items:center;gap:12px;height:52px}
.header-logo{font-family:'DotGothic16',monospace;font-size:1rem;color:#0ff;flex:1;display:flex;align-items:center;gap:8px}
.header-logo:hover{color:#fff}
.header-sub{font-size:.8rem;color:#666}
.back-btn{font-size:.85rem;color:#888;white-space:nowrap;display:flex;align-items:center;gap:3px}
.back-btn:hover{color:#fff}
.list-view{max-width:960px;margin:20px auto;padding:0 16px}
.filters{display:flex;flex-direction:column;gap:10px;margin-bottom:20px}
.search-wrap{display:flex;align-items:center;border:1px solid #444;padding:0 12px;background:#050505}
.search-icon{color:#666;margin-right:8px;display:flex;align-items:center}
.search-input{flex:1;border:none;outline:none;padding:10px 0;font-size:.94rem;background:transparent;color:#fff;font-family:'DotGothic16',monospace}
.search-input::placeholder{color:#555}
.diff-filters{display:flex;flex-wrap:wrap;gap:6px}
.diff-btn{font-family:'DotGothic16',monospace;font-size:.82rem;border:1px solid #444;background:#000;padding:5px 14px;cursor:pointer;color:#666;transition:.1s}
.diff-btn.active,.diff-btn:hover{border-color:#0ff;color:#0ff}
.diff-btn.active{background:rgba(0,255,255,.06)}
.problem-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:12px}
.problem-card{display:block;background:#000;border:1px solid #333;padding:16px;transition:.1s}
.problem-card:hover{border-color:#0ff}
.card-top{display:flex;align-items:center;justify-content:space-between;margin-bottom:6px}
.card-contest{font-size:.78rem;color:#555}
.card-title{font-size:1rem;margin-bottom:8px;line-height:1.8;color:#fff}
.card-tags{display:flex;flex-wrap:wrap;gap:4px}
.tag{display:inline-block;border:1px solid #333;color:#777;padding:2px 8px;font-size:.75rem}
.diff-badge{font-family:'DotGothic16',monospace;display:inline-block;border:1px solid currentColor;padding:2px 8px;font-size:.78rem}
.glossary-top-link{font-family:'DotGothic16',monospace;font-size:.8rem;color:#888;border:1px solid #444;padding:5px 11px;display:inline-flex;align-items:center;gap:4px;white-space:nowrap}
.glossary-top-link:hover{border-color:#fff;color:#fff}
.empty{text-align:center;color:#555;padding:40px;font-size:.9rem}
.detail-view{max-width:860px;margin:20px auto;padding:0 16px}
.detail-contest{font-size:.8rem;color:#666;margin-bottom:4px}
.detail-title{font-size:1.3rem;margin-bottom:10px;line-height:1.8;color:#fff}
.detail-meta{display:flex;align-items:center;flex-wrap:wrap;gap:8px;margin-bottom:20px}
.atcoder-link,.print-link{font-size:.82rem;color:#0ff;border:1px solid #0ff;padding:4px 11px;display:inline-flex;align-items:center;gap:4px}
.atcoder-link:hover,.print-link:hover{background:#0ff;color:#000}
.detail-section{margin-bottom:28px}
.section-title{font-family:'DotGothic16',monospace;font-size:.88rem;margin-bottom:10px;color:#0ff;display:flex;align-items:center;gap:6px;letter-spacing:.1em}
.statement-box{border:1px solid #333;padding:14px;line-height:2;font-size:.94rem}
.statement-box p{margin-bottom:.6em}
.statement-box ul,.statement-box ol{padding-left:1.6em;margin:.4em 0 .6em}
.statement-box li{margin-bottom:.2em}
.statement-box code{background:#111;padding:1px 5px;font-size:.88em;font-family:'Google Sans Code','DotGothic16',monospace}
.constraints-box{border:1px solid #333;padding:10px 14px;font-size:.88rem;margin-top:10px;line-height:1.9;color:#ccc}
.constraints-box ul,.constraints-box ol{padding-left:1.4em}
.constraints-box code{background:#111;padding:1px 4px;font-size:.85em;font-family:'Google Sans Code',monospace}
.statement-note{border-left:3px solid #fff;padding:12px 14px;font-size:.88rem;margin-top:12px;line-height:1.9;background:#050505}
.statement-note p{margin-bottom:.5em}
.statement-note strong{color:#fff}
.statement-note ul,.statement-note ol{padding-left:1.6em;margin:.4em 0 .6em}
.statement-note hr{display:none}
.statement-note table{border-collapse:collapse;margin:.8em 0;font-size:.85rem;width:100%}
.statement-note th,.statement-note td{border:1px solid #333;padding:6px 10px}
.statement-note th{background:#111}
.statement-note-badge{font-family:'DotGothic16',monospace;display:inline-flex;align-items:center;gap:3px;background:#fff;color:#000;font-size:.72rem;padding:2px 9px;margin-bottom:8px}
.constraints-note{border-left:3px solid #888;padding:12px 14px;font-size:.88rem;margin-top:10px;line-height:1.9;background:#050505}
.constraints-note p{margin-bottom:.5em}
.constraints-note ul,.constraints-note ol{padding-left:1.6em;margin:.4em 0 .6em}
.constraints-note hr{display:none}
.constraints-note-badge{font-family:'DotGothic16',monospace;display:inline-flex;align-items:center;gap:3px;background:#444;color:#fff;font-size:.72rem;padding:2px 9px;margin-bottom:8px}
.constraints-note code{background:#111;padding:1px 5px;font-size:.88em;font-family:'Google Sans Code',monospace}
.sample-block{border:1px solid #333;padding:12px;margin-bottom:12px}
.sample-row{display:flex;flex-direction:column;gap:10px}
.sample-label{font-family:'DotGothic16',monospace;font-size:.76rem;color:#0ff;margin-bottom:4px}
.sample-pre{background:#050505;padding:8px 10px;font-size:.88rem;font-family:'Google Sans Code',monospace;overflow-x:auto;white-space:pre;color:#ccc}
.sample-explanation{margin-top:10px;font-size:.88rem;color:#bbb;line-height:1.9;border-top:1px solid #333;padding-top:10px}
.sample-exp-title{font-family:'DotGothic16',monospace;font-size:.76rem;color:#888;margin-bottom:6px;display:flex;align-items:center;gap:3px}
.sample-explanation p{margin-bottom:.5em}
.sample-explanation ul,.sample-explanation ol{padding-left:1.5em;margin:.3em 0 .5em}
.sample-explanation hr{display:none}
.sample-explanation code{background:#111;padding:1px 5px;font-size:.88em;font-family:'Google Sans Code',monospace}
.explanation-box{border-left:3px solid #fff;padding:14px 16px;line-height:1.9;font-size:.94rem;background:#050505}
.easy-box{border-left-color:#aaa}
.explanation-box h2{font-family:'DotGothic16',monospace;font-size:.88rem;margin:14px 0 6px;color:#0ff}
.explanation-box h3{font-size:.9rem;margin:10px 0 4px;color:#ddd}
.explanation-box p{margin-bottom:8px}
.explanation-box pre{background:#050505;border:1px solid #333;padding:8px;font-size:.82rem;overflow-x:auto;margin:8px 0;font-family:'Google Sans Code',monospace}
.explanation-box code{background:#111;padding:1px 5px;font-size:.85em;font-family:'Google Sans Code',monospace}
.explanation-box ul,.explanation-box ol{padding-left:1.4em;margin-bottom:8px}
.lang-tabs{display:flex;gap:4px;margin-bottom:0}
.lang-tab{font-family:'DotGothic16',monospace;font-size:.76rem;border:1px solid #444;background:#000;padding:6px 14px;cursor:pointer;color:#666;transition:.1s}
.lang-tab.active{background:#fff;border-color:#fff;color:#000}
.code-panel{display:none}
.code-panel.active{display:block}
pre.code-block{background:#050505;color:#cfc;border:1px solid #444;border-top:none;padding:14px;font-size:.84rem;font-family:'Google Sans Code',monospace;overflow-x:auto;white-space:pre;line-height:1.9}
pre.code-block code{background:none;padding:0;font-size:inherit;font-family:inherit}
.solution-steps{padding-left:1.5em;margin-top:12px;font-size:.88rem;line-height:1.9;color:#bbb}
.bad-solutions{margin-top:16px}
.bad-solutions-title{font-family:'DotGothic16',monospace;padding:20px 0;font-size:.82rem;display:flex;align-items:center;gap:6px;color:#f66}
.good-solutions-title{font-family:'DotGothic16',monospace;padding:20px 0 8px;font-size:.82rem;display:flex;align-items:center;gap:6px;color:#4f8}
.bad-solution{padding:12px 14px;border-top:1px solid #222}
.bad-solution:first-child{border-top:none}
.bad-solution-label{font-family:'DotGothic16',monospace;font-size:.76rem;color:#888;margin-bottom:6px;display:flex;align-items:center;gap:4px}
pre.code-block-bad{background:#050505;color:#f99;border:1px solid #844;border-top:none;padding:14px;font-size:.84rem;font-family:'Google Sans Code',monospace;overflow-x:auto;white-space:pre;line-height:1.9}
pre.code-block-bad code{background:none;padding:0;font-size:inherit;font-family:inherit}
.bad-lang-tabs{display:flex;gap:4px;margin-bottom:0}
.bad-lang-tab{font-family:'DotGothic16',monospace;font-size:.76rem;border:1px solid #333;background:#000;padding:6px 14px;cursor:pointer;color:#555}
.bad-lang-tab.active{background:#444;border-color:#444;color:#fff}
.bad-code-panel{display:none}
.bad-code-panel.active{display:block}
@media(max-width:600px){.problem-grid{grid-template-columns:1fr}.detail-title{font-size:1.1rem}}
.glossary-nav{border:1px solid #333;padding:12px 16px;margin-bottom:24px}
.glossary-nav ul{list-style:none;display:flex;flex-wrap:wrap;gap:8px}
.glossary-nav a{color:#888;font-size:.84rem;border:1px solid #333;padding:3px 12px;display:inline-block}
.glossary-nav a:hover{border-color:#0ff;color:#0ff}
.glossary-entry{border:1px solid #333;padding:20px;margin-bottom:20px}
.glossary-name{font-size:1.05rem;margin-bottom:8px;display:flex;align-items:center;gap:6px;color:#fff}
.glossary-short{color:#aaa;font-size:.88rem;margin-bottom:14px;padding:8px 12px;border-left:3px solid #0ff;background:#050505}
.glossary-desc{font-size:.9rem;line-height:1.9;margin-bottom:16px;color:#ccc}
.glossary-desc p{margin-bottom:.6em}
.glossary-desc table{border-collapse:collapse;margin:.8em 0;font-size:.85rem;width:100%}
.glossary-desc th,.glossary-desc td{border:1px solid #333;padding:6px 10px}
.glossary-desc th{background:#111}
.glossary-desc code{background:#111;padding:1px 5px;font-size:.88em;font-family:'Google Sans Code',monospace}
.glossary-desc pre{background:#050505;color:#ccc;border:1px solid #333;padding:10px;font-size:.78rem;font-family:'Google Sans Code',monospace;white-space:pre;overflow-x:auto;margin:.6em 0}
.glossary-desc pre code{background:none;padding:0;font-size:inherit}
.code-compare{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:14px}
@media(max-width:700px){.code-compare{grid-template-columns:1fr}}
.code-compare-col{}
.code-compare-label{font-family:'DotGothic16',monospace;font-size:.76rem;padding:4px 10px;display:flex;align-items:center;gap:4px}
.bad-col .code-compare-label{background:#1a0000;color:#f66}
.good-col .code-compare-label{background:#001a00;color:#4f8}
.code-compare-pre{margin:0;border:1px solid #333;border-top:none;padding:10px;font-size:.78rem;font-family:'Google Sans Code',monospace;white-space:pre;overflow-x:auto;line-height:1.9}
.bad-col .code-compare-pre{background:#050505;color:#f99}
.good-col .code-compare-pre{background:#050505;color:#cfc}
.code-compare-pre code{background:none;padding:0;font-size:inherit;font-family:inherit}
.glossary-when{border-left:3px solid #888;padding:8px 12px;font-size:.85rem;margin-bottom:12px;line-height:1.9;color:#bbb;background:#050505}
.glossary-problems{font-size:.82rem;display:flex;align-items:center;flex-wrap:wrap;gap:6px;margin-top:12px}
.glossary-problem-link{color:#888;border:1px solid #333;padding:2px 10px;font-size:.78rem;display:inline-flex;align-items:center;gap:3px}
.glossary-problem-link:hover{border-color:#0ff;color:#0ff}
.glossary-refs{display:flex;flex-wrap:wrap;align-items:center;gap:6px;margin-top:8px;margin-bottom:16px}
.glossary-ref-link{color:#888;border:1px solid #333;padding:2px 10px;font-size:.78rem;display:inline-flex;align-items:center;gap:3px}
.glossary-ref-link:hover{border-color:#0ff;color:#0ff}
.glossary-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:12px}
.glossary-card{display:block;border:1px solid #333;padding:16px;transition:.1s}
.glossary-card:hover{border-color:#0ff}
.glossary-card-name{font-size:.95rem;margin-bottom:6px;display:flex;align-items:center;gap:5px;color:#fff}
.glossary-card-short{font-size:.82rem;color:#777;line-height:1.7}
.cr-link{font-size:.82rem;color:#888;border:1px solid #444;padding:3px 10px;display:inline-flex;align-items:center;gap:3px;white-space:nowrap}
.cr-link:hover{border-color:#fff;color:#fff}
.cr-entry{border:1px solid #333;padding:20px;margin-bottom:20px}
.cr-name{font-size:1rem;margin-bottom:6px;display:flex;align-items:center;gap:6px;color:#fff}
.cr-short{color:#aaa;font-size:.88rem;margin-bottom:12px;padding:6px 10px;border-left:3px solid #aaa;background:#050505}
.cr-body{font-size:.9rem;line-height:1.9;margin-bottom:12px;color:#ccc}
.cr-body p{margin-bottom:.5em}
.cr-body table{border-collapse:collapse;margin:.6em 0;font-size:.85rem}
.cr-body th,.cr-body td{border:1px solid #333;padding:5px 10px}
.cr-body th{background:#111}
.cr-body code{background:#111;padding:1px 5px;font-size:.88em;font-family:'Google Sans Code',monospace}
.cr-code{background:#050505;color:#ccc;border:1px solid #333;padding:10px 14px;font-size:.82rem;font-family:'Google Sans Code',monospace;white-space:pre;overflow-x:auto;margin-bottom:10px;line-height:1.9}
.cr-other{font-size:.82rem;color:#888;padding:6px 10px;border-left:3px solid #444}`
}

func pwGateScript() string {
	return `<script>(function(){` +
		`var H="__VIEW_PASSWORD__";` +
		`if(H.indexOf('__')===0)return;` +
		`if(sessionStorage.getItem('cps')===H)return;` +
		`var o=document.createElement('div');` +
		`o.style='position:fixed;inset:0;background:#1565c0;display:flex;align-items:center;justify-content:center;z-index:9999';` +
		`o.innerHTML='<div style="background:#fff;border-radius:12px;padding:32px;max-width:320px;width:90%;text-align:center">` +
		`<div style="font-size:1.4rem;font-weight:700;margin-bottom:8px">競プロ教材</div>` +
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
		backEl = fmt.Sprintf(`<a class="back-btn" href="%s"><span class="mi">arrow_back</span> %s</a>`, back, backLabel)
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
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=DotGothic16&family=Google+Sans+Code&family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200">
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
    <a class="header-logo" href="%s"><span class="mi">school</span> 競プロ教材</a>
    %s
  </div>
</header>
%s
<script src="https://cdn.jsdelivr.net/npm/prismjs@1/components/prism-core.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/prismjs@1/plugins/autoloader/prism-autoloader.min.js"></script>
<footer style="text-align:center;padding:24px 16px;color:#555;font-size:.78rem;border-top:1px solid #333;margin-top:40px;font-family:'DotGothic16',monospace">本サイトは個人学習目的のみで使用し、商用利用・公開配布はしていません。</footer>
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
			d, diffColor[d], starLabel(d))
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
      <span class="search-icon"><span class="mi">search</span></span>
      <input class="search-input" id="search" placeholder="タイトル・タグで検索..." autocomplete="off">
    </div>
    <div style="display:flex;align-items:center;gap:10px">
      <div class="diff-filters" id="dfilters" style="flex:1">%s</div>
      <a class="glossary-top-link" href="glossary/index.html"><span class="mi">menu_book</span> 用語集</a>
    </div>
  </div>
  <div class="problem-grid" id="grid">%s</div>
  <div class="empty" id="empty" style="display:none">問題が見つかりませんでした</div>
</main>
<script>
(function(){
  var grid=document.getElementById('grid'),search=document.getElementById('search'),
      empty=document.getElementById('empty'),cards=Array.from(grid.querySelectorAll('.problem-card')),curD='';
  function filter(){
    var q=search.value.toLowerCase(),v=0;
    cards.forEach(function(c){
      var ok=(!curD||c.dataset.diff===curD)&&(!q||c.dataset.title.includes(q)||c.dataset.tags.includes(q)||c.dataset.contest.includes(q));
      c.style.display=ok?'':'none'; if(ok)v++;
    });
    empty.style.display=v>0?'none':'';
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

func buildProblem(p Problem, contest string, force bool, glossaryEntries []GlossaryEntry) {
	out := "docs/problems/" + p.ID + ".html"
	if !force {
		if _, err := os.Stat(out); err == nil {
			fmt.Printf("  %s (スキップ)\n", out)
			return
		}
	}
	var sampBuf strings.Builder
	for i, s := range p.Samples {
		fmt.Fprintf(&sampBuf, `
<div class="sample-block">
  <div class="sample-row">
    <div class="sample-col"><div class="sample-label">入力 %d</div><pre class="sample-pre">%s</pre></div>
    <div class="sample-col"><div class="sample-label">出力 %d</div><pre class="sample-pre">%s</pre></div>
  </div>
</div>`, i+1, e(s.Input), i+1, e(s.Output))
	}

	atcLink := ""
	if p.AtcoderURL != "" {
		atcLink = fmt.Sprintf(`<a class="atcoder-link" href="%s" target="_blank" rel="noopener">%s AtCoderで見る</a>`, e(p.AtcoderURL), icon("open_in_new"))
	}
	printLink := fmt.Sprintf(`<a class="print-link" href="../print/%s.html" target="_blank">%s 印刷用</a>`, p.ID, icon("print"))
	codeReadingLink := `<a class="cr-link" href="../code_reading.html">` + icon("menu_book") + ` プログラムの読み方</a>`

	stmtNoteSection := ""
	if p.StatementNote != "" {
		stmtNoteSection = fmt.Sprintf(
			`<div class="statement-note"><span class="statement-note-badge">%s かんたん解説</span>%s</div>`,
			icon("info"), mdToHTML(p.StatementNote))
	}

	constraintsNoteSection := ""
	if p.ConstraintsNote != "" {
		constraintsNoteSection = fmt.Sprintf(
			`<div class="constraints-note"><span class="constraints-note-badge">%s 制約の読み方</span>%s</div>`,
			icon("help"), mdToHTML(p.ConstraintsNote))
	}

	codeSection := buildCodeSection(p)

	glossaryRefsEl := ""
	if len(glossaryEntries) > 0 {
		var refsBuf strings.Builder
		fmt.Fprintf(&refsBuf, `<div class="glossary-refs"><span style="font-size:.8rem;color:#666">%s 使う考え方：</span>`, icon("psychology"))
		for _, g := range glossaryEntries {
			fmt.Fprintf(&refsBuf, `<a class="glossary-ref-link" href="../glossary/%s.html">%s %s</a>`, g.ID, icon("book_2"), e(g.Name))
		}
		refsBuf.WriteString(`</div>`)
		glossaryRefsEl = refsBuf.String()
	}

	body := fmt.Sprintf(`
<main class="detail-view">
  <div class="detail-contest">%s %s</div>
  <h1 class="detail-title">%s</h1>
  <div class="detail-meta">%s %s %s %s %s</div>
  %s
  <section class="detail-section">
    <h2 class="section-title">%s 問題文</h2>
    <div class="statement-box">%s</div>
    %s
    <h2 class="section-title" style="margin-top:20px">%s 入力・出力の例</h2>
    %s
    %s
    %s
  </section>
  %s
</main>`,
		e(contest), e(p.Problem), e(p.Title),
		badge(p.Difficulty), tagSpans(p.Tags), atcLink, printLink, codeReadingLink,
		glossaryRefsEl,
		icon("description"), mdToHTML(p.Statement),
		constraintsBlock(p.Constraints),
		icon("quiz"), sampBuf.String(),
		stmtNoteSection,
		constraintsNoteSection,
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
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=DotGothic16&family=Google+Sans+Code&family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200">
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.css">
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/katex.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.11/dist/contrib/auto-render.min.js"
  onload="renderMathInElement(document.body,{delimiters:[{left:'$$',right:'$$',display:true},{left:'$',right:'$',display:false}],ignoredTags:['script','noscript','style','pre','code']})"></script>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#fff;color:#000;max-width:820px;margin:0 auto;padding:24px 32px 24px 48px;font-size:.93rem;line-height:1.7}
h1{font-size:1.3rem;margin-bottom:8px}
h2{font-size:1rem;font-weight:700;margin:20px 0 8px;color:#000;border-bottom:1.5px solid #000;padding-bottom:4px;display:flex;align-items:center;gap:4px}
h3{font-size:.9rem;font-weight:700;margin:14px 0 6px}
p{margin-bottom:.6em}
ul,ol{padding-left:1.6em;margin:.4em 0 .6em}
li{margin-bottom:.2em}
code{background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:.88em;font-family:monospace}
.mi{font-family:'Material Symbols Outlined';font-weight:normal;font-style:normal;font-size:1.2em;display:inline-block;vertical-align:-.2em;line-height:1;letter-spacing:normal;text-transform:none;white-space:nowrap}
.diff-badge{display:inline-block;border:1.5px solid #000;border-radius:12px;padding:2px 10px;font-size:.72rem;font-weight:700;color:#000}
.tag{display:inline-block;border:1px solid #555;border-radius:12px;padding:2px 9px;font-size:.72rem;font-weight:600;color:#000}
.print-header{border-bottom:2px solid #000;padding-bottom:12px;margin-bottom:20px}
.print-contest{font-size:.8rem;color:#555;margin-bottom:4px}
.print-meta{display:flex;flex-wrap:wrap;gap:8px;margin-top:8px;align-items:center}
.print-section{margin-bottom:22px}
.print-box{border:1px solid #888;border-radius:4px;padding:12px 14px}
.print-box p{margin-bottom:.6em}
.print-box ul,.print-box ol{padding-left:1.6em;margin:.4em 0 .6em}
.print-constraints{border:1px solid #aaa;border-radius:4px;padding:8px 12px;font-size:.85rem;margin-top:8px}
.print-note,.print-easy,.print-exp,.print-constraints-note{padding:10px 0;margin:10px 0 0}
.print-note hr,.print-easy hr,.print-exp hr,.print-constraints-note hr{display:none}
.print-note p,.print-easy p,.print-exp p,.print-constraints-note p{margin-bottom:.5em}
.print-note ul,.print-note ol,.print-easy ul,.print-easy ol,.print-exp ul,.print-exp ol,.print-constraints-note ul,.print-constraints-note ol{padding-left:1.4em;margin:.3em 0 .5em}
.print-note code,.print-easy code,.print-exp code,.print-constraints-note code{background:#f0f0f0;padding:1px 5px;border-radius:3px;font-size:.88em}
.print-note table{border-collapse:collapse;margin:.6em 0;font-size:.85rem;width:100%%}
.print-note th,.print-note td{border:1px solid #aaa;padding:5px 10px}
.print-note th{background:#f0f0f0}
.print-note pre,.print-easy pre,.print-exp pre,.print-constraints-note pre{background:#f5f5f5;border:1px solid #ddd;border-radius:4px;padding:8px;font-size:.82rem;overflow-x:auto;margin:6px 0;white-space:pre-wrap}
.sample-block{border:1px solid #aaa;border-radius:4px;padding:10px;margin-bottom:10px}
.sample-label{font-size:.72rem;font-weight:600;color:#555;margin-bottom:3px}
.sample-pre{background:#f5f5f5;border:1px solid #ddd;border-radius:3px;padding:6px 10px;font-size:.82rem;font-family:monospace;white-space:pre;margin-bottom:8px}
.sample-exp{font-size:.82rem;color:#333;border-top:1px solid #ccc;padding-top:8px;margin-top:8px}
.sample-exp p{margin-bottom:.4em}
.sample-exp ul,.sample-exp ol{padding-left:1.4em;margin:.2em 0 .4em}
.sample-exp-title{font-weight:700;font-size:.78rem;color:#333;margin-bottom:6px;display:flex;align-items:center;gap:3px}
.lang-header{border:1px solid #888;border-bottom:none;border-radius:4px 4px 0 0;padding:4px 12px;font-size:.8rem;font-weight:700;margin-top:14px;background:#f0f0f0;color:#000}
pre.code-print{background:#f8f8f8;border:1px solid #aaa;border-radius:0 4px 4px 4px;padding:12px;font-size:.75rem;font-family:monospace;white-space:pre-wrap;word-break:break-all;line-height:1.55;margin-bottom:0}
pre.code-print code{background:none;padding:0;font-size:inherit}
.bad-code-section{border-top:1px dashed #aaa;margin-top:20px;padding-top:10px}
.bad-code-label{font-size:.78rem;font-weight:700;color:#555;margin-bottom:4px}
pre.code-print-bad{background:#f8f8f8;border:1px dashed #aaa;border-radius:4px;padding:12px;font-size:.75rem;font-family:monospace;white-space:pre-wrap;word-break:break-all;line-height:1.55}
.print-tagline{font-size:.92rem;color:#333;margin-bottom:18px;padding:10px 14px;border-left:3px solid #888;background:#fafafa;line-height:1.6}
.compare-block{border:2px solid #888;border-radius:6px;margin-bottom:6px;overflow:hidden;page-break-inside:avoid}
.bad-block{border-style:dashed;border-color:#666}
.good-block{border-color:#333}
.compare-block-header{padding:8px 14px;font-size:.85rem;font-weight:700;display:flex;align-items:flex-start;gap:6px;border-bottom:1.5px solid #aaa;background:#ececec;line-height:1.5}
.bad-block .compare-block-header{border-bottom-style:dashed}
.compare-block pre.code-print{border:none;border-radius:0;margin:0;background:#f8f8f8}
.compare-block.good-block pre.code-print{background:#f0f0f0}
.compare-arrow{text-align:center;padding:10px 0;font-size:.88rem;color:#333;font-weight:700;display:flex;align-items:center;justify-content:center;gap:6px}
.print-btn{position:fixed;bottom:20px;right:20px;background:#333;color:#fff;border:none;border-radius:8px;padding:10px 20px;font-size:.9rem;cursor:pointer;font-weight:600;display:flex;align-items:center;gap:6px;box-shadow:0 2px 8px rgba(0,0,0,.25)}
.print-btn:hover{background:#000}
@media print{
  @page{margin:1.5cm 1.5cm 1.5cm 2.5cm}
  body{max-width:none;padding:0}
  .print-btn{display:none}
  pre.code-print,pre.code-print-bad{page-break-inside:avoid;white-space:pre-wrap}
  h2{page-break-after:avoid}
  .compare-block{page-break-inside:avoid}
  .print-note,.print-easy,.print-exp,.sample-block,.print-section{page-break-inside:avoid}
  .lang-header{page-break-after:avoid}
}
</style>
</head>
<body>
%s
%s
<button class="print-btn" onclick="window.print()"><span class="mi">print</span> 印刷</button>
<footer style="text-align:center;padding:24px 16px;color:#666;font-size:.78rem;border-top:1px solid #ccc;margin-top:40px">本サイトは個人学習目的のみで使用し、商用利用・公開配布はしていません。</footer>
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

	fmt.Fprintf(&buf, `<section class="print-section"><h2>%s 問題文</h2><div class="print-box">%s</div>`, icon("description"), mdToHTML(p.Statement))
	if p.Constraints != "" {
		fmt.Fprintf(&buf, `<div class="print-constraints"><strong>制約</strong>%s</div>`, mdToHTML(p.Constraints))
	}
	if len(p.Samples) > 0 {
		buf.WriteString(fmt.Sprintf(`<div class="print-samples-inline"><h3>%s 入力・出力の例</h3>`, icon("quiz")))
		for i, s := range p.Samples {
			fmt.Fprintf(&buf, `<div class="sample-block"><div class="sample-label">入力 %d</div><pre class="sample-pre">%s</pre><div class="sample-label">出力 %d</div><pre class="sample-pre">%s</pre></div>`,
				i+1, e(s.Input), i+1, e(s.Output))
		}
		buf.WriteString(`</div>`)
	}
	buf.WriteString(`</section>`)

	if p.StatementNote != "" {
		fmt.Fprintf(&buf, `<section class="print-section"><h2>%s かんたん解説</h2><div class="print-note">%s</div></section>`, icon("info"), mdToHTML(p.StatementNote))
	}
	if p.ConstraintsNote != "" {
		fmt.Fprintf(&buf, `<section class="print-section"><h2>%s 制約の読み方</h2><div class="print-constraints-note">%s</div></section>`, icon("help"), mdToHTML(p.ConstraintsNote))
	}

	if sol := p.Solutions["python"]; sol != nil {
		buf.WriteString(fmt.Sprintf(`<section class="print-section"><h2>%s 解答コード（Python）</h2>`, icon("code")))
		if bad := p.BadSolutions["python"]; bad != nil {
			buf.WriteString(fmt.Sprintf(`<div class="bad-code-section"><h3>%s 悪い例（結果は合ってるけど…）</h3>`, icon("sentiment_dissatisfied")))
			fmt.Fprintf(&buf, `<pre class="code-print-bad"><code>%s</code></pre></div>`, e(bad.Code))
		}
		fmt.Fprintf(&buf, `<div class="lang-header">Python（正しい解答）</div><pre class="code-print"><code>%s</code></pre>`, e(sol.Code))
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

	// bad solutions
	badSection := ""
	if len(p.BadSolutions) > 0 {
		var badTabs, badPanels strings.Builder
		first = true
		for _, lang := range langDefs {
			bad, ok := p.BadSolutions[lang.key]
			if !ok || bad == nil {
				continue
			}
			active := first
			first = false
			cls := "bad-lang-tab"
			if active {
				cls += " active"
			}
			fmt.Fprintf(&badTabs, `<button class="%s" data-lang="bad-%s">%s</button>`, cls, lang.key, lang.label)
			panelCls := "bad-code-panel"
			if active {
				panelCls += " active"
			}
			fmt.Fprintf(&badPanels, `<div class="%s" id="panel-bad-%s"><pre class="code-block-bad language-%s"><code class="language-%s">%s</code></pre></div>`,
				panelCls, lang.key, lang.key, lang.key, e(bad.Code))
		}
		badSection = fmt.Sprintf(`
<div class="bad-solutions">
  <div class="bad-solutions-title">%s 悪い例（結果は合ってるけど…）</div>
  <div class="bad-lang-tabs" id="bltabs">%s</div>
  %s
</div>
<script>
document.getElementById('bltabs').addEventListener('click',function(e){
  var b=e.target.closest('.bad-lang-tab'); if(!b)return;
  document.querySelectorAll('.bad-lang-tab').forEach(function(x){x.classList.remove('active');});
  document.querySelectorAll('.bad-code-panel').forEach(function(x){x.classList.remove('active');});
  b.classList.add('active');
  var p=document.getElementById('panel-'+b.dataset.lang); if(p)p.classList.add('active');
});
</script>`, icon("sentiment_dissatisfied"), badTabs.String(), badPanels.String())
	}

	return fmt.Sprintf(`
<section class="detail-section">
  <h2 class="section-title">%s 解答コード</h2>
  %s
  <div class="good-solutions-title">%s 正しいコード</div>
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
</script>`, icon("code"), badSection, icon("check_circle"), tabs.String(), panels.String())
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

func buildCodeReading() {
	b, err := os.ReadFile("data/code_reading.json")
	if err != nil {
		fmt.Println("  data/code_reading.json が見つかりません（スキップ）")
		return
	}
	var entries []CodeReadingEntry
	json.Unmarshal(b, &entries)

	var buf strings.Builder
	for _, cr := range entries {
		codeEl := ""
		if cr.PythonCode != "" {
			codeEl = fmt.Sprintf(`<pre class="cr-code">%s</pre>`, e(cr.PythonCode))
		}
		otherEl := ""
		if cr.OtherNote != "" {
			otherEl = fmt.Sprintf(`<div class="cr-other">%s %s</div>`, icon("translate"), e(cr.OtherNote))
		}
		fmt.Fprintf(&buf, `
<div class="cr-entry" id="%s">
  <div class="cr-name">%s %s</div>
  <div class="cr-short">%s</div>
  <div class="cr-body">%s</div>
  %s
  %s
</div>`, cr.ID, icon("code"), e(cr.Name), e(cr.Short), mdToHTML(cr.Body), codeEl, otherEl)
	}

	printLnk := fmt.Sprintf(`<a class="print-link" href="print/code_reading.html" target="_blank">%s 印刷用</a>`, icon("print"))
	body := fmt.Sprintf(`
<main class="detail-view">
  <h1 class="detail-title" style="margin-bottom:6px">%s プログラムの読み方</h1>
  <div style="margin-bottom:20px">%s</div>
  <p style="color:#666;font-size:.88rem;margin-bottom:24px">コードに出てくる言葉の意味をまとめました。わからない言葉があったら調べてみよう。</p>
  %s
</main>`, icon("menu_book"), printLnk, buf.String())

	writeFile("docs/code_reading.html", shell("プログラムの読み方", "", "", "", "", body))
	fmt.Println("  docs/code_reading.html")

	// 印刷ページ
	var printBuf strings.Builder
	fmt.Fprintf(&printBuf, `<div class="print-header">
<h1>プログラムの読み方</h1>
<p style="font-size:.85rem;color:#555;margin-top:6px">コードに出てくる言葉の意味一覧</p>
</div>`)
	for _, cr := range entries {
		codeEl := ""
		if cr.PythonCode != "" {
			codeEl = fmt.Sprintf(`<div class="lang-header">Python の例</div><pre class="code-print"><code>%s</code></pre>`, e(cr.PythonCode))
		}
		otherEl := ""
		if cr.OtherNote != "" {
			otherEl = fmt.Sprintf(`<p style="font-size:.8rem;color:#555;margin-top:6px">%s</p>`, e(cr.OtherNote))
		}
		fmt.Fprintf(&printBuf, `<section class="print-section">
<h2>%s %s</h2>
<p class="print-tagline">%s</p>
<div class="print-box">%s</div>
%s
%s
</section>`, icon("code"), e(cr.Name), e(cr.Short), mdToHTML(cr.Body), codeEl, otherEl)
	}
	writeFile("docs/print/code_reading.html", printShell("プログラムの読み方", printBuf.String()))
	fmt.Println("  docs/print/code_reading.html")
}

func buildGlossary(index []IndexEntry) {
	b, err := os.ReadFile("data/glossary.json")
	if err != nil {
		fmt.Println("  data/glossary.json が見つかりません（スキップ）")
		return
	}
	var entries []GlossaryEntry
	json.Unmarshal(b, &entries)

	// index page
	var cardBuf strings.Builder
	for _, g := range entries {
		fmt.Fprintf(&cardBuf, `
<a class="glossary-card" href="%s.html">
  <div class="glossary-card-name">%s %s</div>
  <div class="glossary-card-short">%s</div>
</a>`, g.ID, icon("book_2"), e(g.Name), e(g.Short))
	}
	indexBody := fmt.Sprintf(`
<main class="list-view">
  <h1 class="detail-title" style="margin-bottom:20px">%s 用語集</h1>
  <p style="color:#666;font-size:.9rem;margin-bottom:20px">競プロでよく使うアルゴリズムや考え方をまとめました。</p>
  <div class="glossary-grid">%s</div>
</main>`, icon("menu_book"), cardBuf.String())
	writeFile("docs/glossary/index.html", shell("用語集", "../", "", "", "", indexBody))
	fmt.Println("  docs/glossary/index.html")

	// individual pages
	for _, g := range entries {
		printLink := fmt.Sprintf(`<a class="print-link" href="../print/glossary_%s.html" target="_blank">%s 印刷用</a>`, g.ID, icon("print"))

		body := fmt.Sprintf(`
<main class="detail-view">
  <h1 class="detail-title">%s %s</h1>
  <div class="detail-meta" style="margin-bottom:16px">%s</div>
  <p class="glossary-short" style="margin-bottom:24px">%s</p>
  <section class="detail-section">
    <h2 class="section-title">%s 解説</h2>
    <div class="glossary-desc statement-box">%s</div>
  </section>
  <section class="detail-section">
    <h2 class="section-title">%s 使わない場合 vs 使う場合</h2>
    <div class="code-compare">
      <div class="code-compare-col bad-col">
        <div class="code-compare-label">%s %s</div>
        <pre class="code-compare-pre"><code>%s</code></pre>
      </div>
      <div class="code-compare-col good-col">
        <div class="code-compare-label">%s %s</div>
        <pre class="code-compare-pre"><code>%s</code></pre>
      </div>
    </div>
  </section>
  <section class="detail-section">
    <h2 class="section-title">%s いつ使う？</h2>
    <div class="glossary-when">%s</div>
  </section>
</main>`,
			icon("book_2"), e(g.Name),
			printLink,
			e(g.Short),
			icon("auto_stories"), mdToHTML(g.Description),
			icon("compare"),
			icon("close"), e(g.WithoutLabel), e(g.WithoutCode),
			icon("check"), e(g.WithLabel), e(g.WithCode),
			icon("help_outline"), mdToHTML(g.WhenToUse),
		)

		out := "docs/glossary/" + g.ID + ".html"
		writeFile(out, shell(g.Name, "../", "index.html", "用語集", "", body))
		fmt.Printf("  %s\n", out)

		buildGlossaryPrintPage(g)
	}
}

func buildGlossaryPrintPage(g GlossaryEntry) {
	var buf strings.Builder

	fmt.Fprintf(&buf, `<div class="print-header">
<div class="print-contest">用語集</div>
<h1>%s</h1>
</div>
<p class="print-tagline">%s</p>`, e(g.Name), e(g.Short))

	fmt.Fprintf(&buf, `<section class="print-section"><h2>%s 解説</h2><div class="print-box">%s</div></section>`,
		icon("auto_stories"), mdToHTML(g.Description))

	fmt.Fprintf(&buf, `<section class="print-section"><h2>%s 使わない場合 vs 使う場合</h2>
<div class="compare-block bad-block">
  <div class="compare-block-header">%s %s</div>
  <pre class="code-print"><code>%s</code></pre>
</div>
<div class="compare-arrow">%s この書き方を改善すると…</div>
<div class="compare-block good-block">
  <div class="compare-block-header">%s %s</div>
  <pre class="code-print"><code>%s</code></pre>
</div>
</section>`,
		icon("compare"),
		icon("close"), e(g.WithoutLabel), e(g.WithoutCode),
		icon("arrow_downward"),
		icon("check"), e(g.WithLabel), e(g.WithCode))

	fmt.Fprintf(&buf, `<section class="print-section"><h2>%s いつ使う？</h2><div class="print-box">%s</div></section>`,
		icon("help_outline"), mdToHTML(g.WhenToUse))

	out := "docs/print/glossary_" + g.ID + ".html"
	writeFile(out, printShell(g.Name, buf.String()))
	fmt.Printf("  %s\n", out)
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

func buildExecution(index []IndexEntry) {
	os.MkdirAll("execution", 0755)
	seen := make(map[string]bool)
	for _, meta := range index {
		if seen[meta.File] {
			continue
		}
		seen[meta.File] = true
		b, err := os.ReadFile("data/problems/" + meta.File + ".json")
		if err != nil {
			continue
		}
		var cf ContestFile
		json.Unmarshal(b, &cf)
		for _, p := range cf.Problems {
			dir := "execution/" + p.ID
			os.MkdirAll(dir, 0755)
			for i, s := range p.Samples {
				fname := fmt.Sprintf("%s/input%d.txt", dir, i+1)
				os.WriteFile(fname, []byte(s.Input), 0644)
			}
			solPath := dir + "/solution.py"
			if _, err := os.Stat(solPath); os.IsNotExist(err) {
				os.WriteFile(solPath, []byte("# ここにコードを書いてね\n"), 0644)
			}
			fmt.Printf("  execution/%s/ (%d 入力)\n", p.ID, len(p.Samples))
		}
	}
}

func cmdBuild(force bool) {
	if force {
		os.RemoveAll("docs")
	}
	os.MkdirAll("docs/problems", 0755)
	os.MkdirAll("docs/print", 0755)
	os.MkdirAll("docs/glossary", 0755)

	idxB, err := os.ReadFile("data/index.json")
	if err != nil {
		log.Fatal("data/index.json が見つかりません。先に convert を実行してください。")
	}
	var index []IndexEntry
	json.Unmarshal(idxB, &index)

	// build problem → glossary entries map
	problemGlossary := make(map[string][]GlossaryEntry)
	if gb, err2 := os.ReadFile("data/glossary.json"); err2 == nil {
		var glossary []GlossaryEntry
		json.Unmarshal(gb, &glossary)
		for _, g := range glossary {
			for _, pid := range g.Problems {
				problemGlossary[pid] = append(problemGlossary[pid], g)
			}
		}
	}

	writeFile("docs/style.css", css())
	fmt.Println("  docs/style.css")
	fmt.Println("HTML 生成中...")
	buildIndex(index)
	buildGlossary(index)
	buildCodeReading()

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
			buildProblem(p, cf.Contest, force, problemGlossary[p.ID])
			buildPrintPage(p, cf.Contest, force)
		}
	}

	fmt.Println("\n実行環境生成中...")
	buildExecution(index)
	fmt.Println("\n完了 ✓  →  docs/  execution/")
}

func cmdServe() {
	fmt.Println("http://localhost:8080 で起動中...")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("docs"))))
}

// ── entry point ───────────────────────────────────────────────────────────────

func main() {
	// src/ 内から実行された場合はプロジェクトルートに移動
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Chdir("..")
	}

	force := false
	serve := false
	for _, a := range os.Args[1:] {
		switch a {
		case "-f", "--force":
			force = true
		case "serve":
			serve = true
		}
	}
	if serve {
		cmdServe()
	} else {
		cmdBuild(force)
	}
}
