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

// tipsForAge は年齢に応じて言語 Tips を返します。
// 幼い子向け（age <= 9）は Tips に漢字が多いため省略します。
func tipsForAge(language string, age int) string {
	if age != 0 && age <= 9 {
		return ""
	}
	return getLang(language).Tips
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

	lang := getLang(language)
	tips := tipsForAge(language, age)
	if tips != "" {
		tips = tips + "\n\n"
	}

	return strings.TrimSpace(`あなたはプログラミング先生「センセイ」です。
今日の学習言語は「` + lang.Name + `」です。

` + tips + `【返答ルール】
- ヒントだけ出す。答えのコードは絶対に書かない
- 1〜3文で短く答える
- 最後に考えさせる問いかけをする
- 間違えても責めず、励ます

` + ageInst)
}

func buildGeneratePrompt(language, execMode string, birthYear int) string {
	lang := getLang(language)
	ageInst := buildAgeInstruction(birthYear)
	age := ageFromBirthYear(birthYear)

	tips := tipsForAge(language, age)
	if tips != "" {
		tips = tips + "\n\n"
	}

	modeRules := lang.BrowserRules
	if execMode == "file" {
		modeRules = lang.FileRules
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
ユーザーの要望に合った ` + lang.Name + ` プログラムを作成してください。

【最重要】必ず最初に ` + "```" + language + ` でコードブロックを書いてから、その後に解説を書くこと。コードブロックなしで説明だけ返すのは禁止。

` + tips + responseFormat + `

` + modeRules + `

` + ageInst)
}

func buildReviewPrompt(language string, birthYear int) string {
	ageInst := buildAgeInstruction(birthYear)
	return strings.TrimSpace(`あなたはプログラミング先生「センセイ」です。
以下のコードについて復習問題を3つ作ってください。

【問題作成のルール】
- コードの内容を理解しているか確認する質問
- 答えにコードを含めない
- 「Q1:」「Q2:」「Q3:」の形式で書く
- 考えさせる質問にする（〇か×ではなく、説明させる）

` + ageInst)
}

func buildBlankPrompt(language string, birthYear int) string {
	ageInst := buildAgeInstruction(birthYear)
	return strings.TrimSpace(`あなたはプログラミング先生「センセイ」です。
以下のコードの重要な部分を「___」（アンダースコア3つ）に置き換えて、穴埋め問題を作ってください。

【穴埋めのルール】
- 3〜5箇所だけ「___」に置き換える
- 変数名・数値・演算子・関数名・条件式の値など、コードの理解に関わる重要な部分を選ぶ
- コメント行（#で始まる行）は変えない
- コードの構造（インデント・行数）は絶対に変えない
- コードブロック（` + "```" + language + `\n...\n` + "```" + `）でコードのみを返す（説明は一切不要）

` + ageInst)
}

func extractCodeBlock(text string) (code, explanation string) {
	const marker = "```"

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

	// ``` がない場合：コードらしい行が含まれていればコードとみなす（フォールバック）
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

	fmt.Printf("[extractCodeBlock] fallback: returning full text as code\n")
	return strings.TrimSpace(text), ""
}
