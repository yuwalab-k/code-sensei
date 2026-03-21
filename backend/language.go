package main

// Language はプログラミング言語の定義です。
// 新しい言語を追加するには、languages マップにエントリを追加するだけです。
type Language struct {
	ID           string
	Name         string
	Tips         string // ヒント・レビュー用の言語特有の注意点
	BrowserRules string // ブラウザ実行モード用のコードルール
	FileRules    string // ファイル実行モード用のコードルール
}

var languages = map[string]Language{
	"python": {
		ID:   "python",
		Name: "Python",
		Tips: `【Pythonの特徴とよくあるミス】
- インデント（字下げ）がとても大切。ずれると「IndentationError」が出る
- if や for や def の後には「:」（コロン）が必ず必要
- 文字（文字列）は " " か ' ' で囲む
- 表示は print() を使う。括弧を忘れずに
- よくあるエラー:
  - IndentationError → 字下げがおかしい
  - SyntaxError → 記号（:や括弧）が足りない・多い
  - NameError → 変数名のスペルミス
  - TypeError → 数字と文字を足し算しようとした`,
		BrowserRules: `【実行環境の制約（ブラウザ実行モード）】
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
- ゲームの「結果」は乱数で決めてprint()で表示する（例：じゃんけんなら両者の手と勝敗をprint()で出す）`,
		FileRules: `【コードのルール（ファイル実行モード）】
- ユーザーのパソコンで実行するコードを生成する
- input() を使ってユーザーと対話してよい
- 使えるライブラリ: random, math, os（読み取りのみ）, csv, json, datetime, collections
- ファイルの書き込みはカレントディレクトリのみ（絶対パス禁止）
- os.system(), subprocess など外部コマンドは絶対に使わない
- pip install などパッケージのインストールは使わない
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする`,
	},

	"javascript": {
		ID:   "javascript",
		Name: "JavaScript (p5.js)",
		Tips: `【p5.jsの特徴とよくあるミス】
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
		BrowserRules: `【p5.jsのコードルール（ブラウザ実行モード）】
- p5.js のスケッチとして動作するコードを書く
- setup() と draw() を必ず定義する
- ブラウザ上のキャンバスで動作する
- 外部ライブラリの読み込みはしない（p5.js のみ使用可）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
		FileRules: `【p5.jsのコードルール（ファイル実行モード）】
- p5.js のスケッチとして動作するコードを書く
- setup() と draw() を必ず定義する
- 外部ライブラリの読み込みはしない（p5.js のみ使用可）
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
	},

	"go": {
		ID:   "go",
		Name: "Go",
		Tips: `【Goの特徴とよくあるミス】
- 変数は := で作る（例: x := 10）
- 表示は fmt.Println() を使う
- 型（かた）が大切。数字と文字は混ぜられない
- { } でブロックを作る。{ は行末に書く
- よくあるエラー:
  - undefined → 変数名のミス
  - syntax error → { } や () の対応ミス`,
		BrowserRules: `【Goのコードルール（ブラウザ実行モード）】
- package main と func main() を必ず書く
- fmt.Println() で結果を表示する
- 標準ライブラリのみ使う（fmt, math, strings, strconv, sort など）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
		FileRules: `【Goのコードルール（ファイル実行モード）】
- package main と func main() を必ず書く
- fmt.Println() で結果を表示する
- 標準ライブラリのみ使う
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
	},
}

// getLang は言語IDから Language を返します。未登録の場合は最低限の情報を返します。
func getLang(id string) Language {
	if l, ok := languages[id]; ok {
		return l
	}
	return Language{ID: id, Name: id}
}
