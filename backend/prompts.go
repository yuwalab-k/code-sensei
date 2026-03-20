package main

import (
	"fmt"
	"strings"
	"time"
)

func ageFromBirthYear(birthYear int) int {
	if birthYear == 0 {
		return 0
	}
	return time.Now().Year() - birthYear
}

func buildAgeInstruction(birthYear int) string {
	age := ageFromBirthYear(birthYear)
	if age == 0 {
		return "相手は小学生（9〜12歳）です。やさしい言葉で説明してください。"
	}
	switch {
	case age <= 7:
		return fmt.Sprintf(`【ぜったいまもること：%d さいの子ども むけ】
- かんじを いっさい つかわない。ひらがなと カタカナだけ つかう
- むずかしい ことばを つかわない
- 1つの ぶんを みじかく する（15もじ いない）
- やさしく、ゆっくり つたえる
- れい：「へんすう」→「いれもの」、「くりかえし」→「なんどもやる」`, age)
	case age <= 9:
		return fmt.Sprintf(`【かならずまもること：%d歳の小学生むけ】
- かんじはなるべく つかわない。つかう場合は かならず ふりがなをつける
- むずかしい専門用語は やさしいことばに かえる
- 1文を みじかく する
- れい：「変数」→「いれもの（へんすう）」`, age)
	case age <= 12:
		return fmt.Sprintf("相手は%d歳の小学生です。小学生がわかるやさしい言葉で説明してください。難しい専門用語はひらがなや身近な例えに置き換えてください。", age)
	case age <= 15:
		return fmt.Sprintf("相手は%d歳の中学生です。中学生がわかる言葉で説明してください。英語の専門用語も使ってよいですが、初出時は日本語でも説明してください。", age)
	case age <= 18:
		return fmt.Sprintf("相手は%d歳の高校生です。高校生向けのわかりやすい説明をしてください。専門用語も適切に使ってください。", age)
	default:
		return fmt.Sprintf("相手は%d歳です。わかりやすく丁寧に説明してください。", age)
	}
}

var langTips = map[string]string{
	"python": `
【Pythonの特徴とよくあるミス】
- インデント（字下げ）がとても大切。ずれると「IndentationError」が出る
- if や for や def の後には「:」（コロン）が必ず必要
- 文字（文字列）は " " か ' ' で囲む
- 表示は print() を使う。括弧を忘れずに
- よくあるエラー:
  - IndentationError → 字下げがおかしい
  - SyntaxError → 記号（:や括弧）が足りない・多い
  - NameError → 変数名のスペルミス
  - TypeError → 数字と文字を足し算しようとした`,

	"javascript": `
【p5.jsの特徴とよくあるミス】
- setup() は最初に1回だけ動く（準備の場所）
- draw() はずっと繰り返し動く（アニメーションの場所）
- 変数は let か const で作る（例: let x = 0）
- 座標は左上が(0,0)で、右に行くとxが増え、下に行くとyが増える
- よく使う関数: circle(x,y,大きさ), rect(x,y,幅,高さ), fill(色), background(色)
- キー操作: keyIsDown(LEFT_ARROW), keyIsDown(RIGHT_ARROW)
- よくあるエラー:
  - is not defined → 変数名や関数名のスペルミス
  - { } の対応がずれている
  - setup や draw のスペルミス`,

	"go": `
【Goの特徴とよくあるミス】
- 変数は := で作る（例: x := 10）
- 表示は fmt.Println() を使う
- 型（かた）が大切。数字と文字は混ぜられない
- { } でブロックを作る。{ は行末に書く
- よくあるエラー:
  - undefined → 変数名のミス
  - syntax error → { } や () の対応ミス`,
}

func langName(language string) string {
	names := map[string]string{
		"python":     "Python",
		"javascript": "JavaScript (p5.js)",
		"go":         "Go",
	}
	if n, ok := names[language]; ok {
		return n
	}
	return language
}

func buildHintPrompt(language string, reviewMode bool, birthYear int) string {
	ageInst := buildAgeInstruction(birthYear)
	age := ageFromBirthYear(birthYear)

	if reviewMode {
		return strings.TrimSpace(`あなたはプログラミング先生「センセイ」です。
今、生徒は復習問題に取り組んでいます。

【返答ルール】
- 答えは絶対に言わない
- ヒントだけ出す（1〜2文）
- 「もう少しだよ！」「いいところに気づいたね！」など励ます
- コードを書かない

` + ageInst)
	}

	// 幼い子向けは langTips（漢字まじり）を省いて言葉遣いを優先する
	var tips string
	if age == 0 || age > 9 {
		tips = langTips[language]
	}

	prompt := `あなたはプログラミング先生「センセイ」です。
今日の学習言語は「` + langName(language) + `」です。
` + tips + `
【返答ルール】
- ヒントだけ出す。答えのコードは絶対に書かない
- 1〜3文で短く答える
- 最後に考えさせる問いかけをする
- 間違えても責めず、励ます

` + ageInst

	return strings.TrimSpace(prompt)
}

func buildGeneratePrompt(language, execMode string, birthYear int) string {
	lang := langName(language)
	ageInst := buildAgeInstruction(birthYear)
	age := ageFromBirthYear(birthYear)

	// 幼い子向けは langTips（漢字まじり）を省く
	var tips string
	if age == 0 || age > 9 {
		tips = langTips[language]
	}

	var modeRules string
	if execMode == "file" {
		modeRules = `【コードのルール（ファイル実行モード）】
- ユーザーのパソコンで実行するコードを生成する
- input() を使ってユーザーと対話してよい
- 使えるライブラリ: random, math, os（読み取りのみ）, csv, json, datetime, collections
- ファイルの書き込みはカレントディレクトリのみ（絶対パス禁止）
- os.system(), subprocess など外部コマンドは絶対に使わない
- pip install などパッケージのインストールは使わない
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする`
	} else {
		modeRules = `【実行環境の制約（ブラウザ実行モード）】
このモードはコードを一度だけ実行してprint()の出力を表示するだけです。
ユーザーとリアルタイムでやり取りする機能（入力待ち・ループ応答）は一切使えません。

【コードの絶対ルール】
- print() で結果を表示する
- input() は絶対に使わない（入力待ちができないため、プログラムがフリーズする）
- while True や無限ループは絶対に使わない
- import できるのは random と math のみ
- ファイルの読み書き（open()）はしない
- 30行以内のシンプルなコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする

【クイズ・ゲームを作るときのルール】
ユーザーが「掛け算クイズ」「数当てゲーム」など対話が必要そうなものを要望した場合、
input()を使わずに以下のパターンで作ること：
- 問題をランダムに複数生成してprint()で表示し、答えも一緒にprint()で表示する
  例：「3 × 7 = ?  →  答えは 21」のように問題と答えをセットで表示する
- または「問題だけ表示して、答えは自分でノートに書いてみよう！」と書いてから答えを最後にまとめて表示する
- ゲームの「結果」は乱数で決めてprint()で表示する（例：じゃんけんなら両者の手と勝敗をprint()で出す）`
	}

	responseFormat := `【返答の形式】必ず以下の形式で返答してください：
1. コードブロック（` + "```" + language + `\n...\n` + "```" + `）でコードを書く
2. コードの後に、以下の構成で詳しい解説を日本語で書く：

■ このプログラムで何をしているか（2〜3文で全体像を説明）

■ コードの解説（重要な行を1行ずつ説明）
- 各行やブロックを「〇〇行目：△△をしているよ」の形で全部説明する
- 変数・関数・ループ・条件分岐など、使っている概念を小学生でもわかる言葉で説明する
- 「なぜこう書くのか」「どういう意味があるのか」まで丁寧に説明する

■ プログラミングのポイント（このコードで学べること）
- このプログラムで使ったテクニックや考え方を2〜3個ピックアップして説明する
- 「このfor文は〜に使える！」のように応用も教える

解説は短くしなくていいです。小学生が「なるほど！」となるくらい丁寧に、でも難しい言葉は使わずに書いてください。`

	return strings.TrimSpace(`あなたはプログラミング先生「センセイ」です。
ユーザーの要望に合った ` + lang + ` プログラムを作成してください。

【最重要】必ず最初に ` + "```" + language + ` でコードブロックを書いてから、その後に解説を書くこと。コードブロックなしで説明だけ返すのは禁止。
` + tips + `

` + responseFormat + `

` + modeRules + `

` + ageInst)
}

func extractCodeBlock(text string) (code, explanation string) {
	const marker = "```"

	// ``` があれば通常パース
	start := strings.Index(text, marker)
	if start != -1 {
		afterMarker := text[start+3:]
		// 言語名の行をスキップ
		newline := strings.Index(afterMarker, "\n")
		if newline != -1 {
			codeContent := afterMarker[newline+1:]
			end := strings.Index(codeContent, marker)
			if end != -1 {
				code = strings.TrimSpace(codeContent[:end])
				explanation = strings.TrimSpace(text[:start] + " " + codeContent[end+3:])
				return
			}
			// 閉じ ``` がない場合はそのままコードとして扱う
			return strings.TrimSpace(codeContent), ""
		}
	}

	// ``` がない場合：Pythonらしい行が含まれていればコードとみなす
	lines := strings.Split(text, "\n")
	var codeLines, otherLines []string
	inCode := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		looksLikeCode := strings.HasPrefix(trimmed, "print(") ||
			strings.HasPrefix(trimmed, "for ") ||
			strings.HasPrefix(trimmed, "if ") ||
			strings.HasPrefix(trimmed, "def ") ||
			strings.HasPrefix(trimmed, "import ") ||
			strings.HasPrefix(trimmed, "while ") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.Contains(trimmed, " = ") ||
			strings.Contains(trimmed, "+=") ||
			strings.HasPrefix(trimmed, "return ")
		if looksLikeCode {
			inCode = true
		}
		if inCode && (looksLikeCode || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") || trimmed == "") {
			codeLines = append(codeLines, line)
		} else if !inCode {
			otherLines = append(otherLines, line)
		}
	}
	if len(codeLines) > 0 {
		fmt.Printf("[extractCodeBlock] no fence found, extracted %d code lines from response\n", len(codeLines))
		return strings.TrimSpace(strings.Join(codeLines, "\n")),
			strings.TrimSpace(strings.Join(otherLines, "\n"))
	}

	// 何も取れなかったら全文をコードとして返す
	fmt.Printf("[extractCodeBlock] fallback: returning full text as code\n")
	return strings.TrimSpace(text), ""
}
