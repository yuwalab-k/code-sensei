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
		Name: "JavaScript",
		Tips: `【JavaScriptの特徴とよくあるミス】
- 変数は let か const で作る（例: let x = 10）
- 表示は console.log() を使う。括弧を忘れずに
- 文字（文字列）は " " か ' ' か ` + "`" + ` ` + "`" + ` で囲む
- { } でブロックを作る。{ は行末に書く
- よくあるエラー:
  - is not defined → 変数名のスペルミス
  - SyntaxError → { } や () の対応ミス
  - TypeError → 使えない操作をしようとした`,
		BrowserRules: `【JavaScriptのコードルール（ブラウザ実行モード）】
- console.log() で結果を表示する
- Node.js で実行するコードを書く
- prompt() / alert() など対話的な関数は使わない
- 外部ライブラリは使わない（標準の Math, Array, String のみ）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする`,
		FileRules: `【JavaScriptのコードルール（ファイル実行モード）】
- console.log() で結果を表示する
- Node.js で実行するコードを書く
- 外部ライブラリは使わない（標準機能のみ）
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
	},

	"ruby": {
		ID:   "ruby",
		Name: "Ruby",
		Tips: `【Rubyの特徴とよくあるミス】
- 変数は小文字で始める（例: name = "太郎"）
- 表示は puts か print を使う（puts は改行あり、print は改行なし）
- if / unless / while などの制御構文は end で終わる
- 文字（文字列）は " " か ' ' で囲む
- よくあるエラー:
  - undefined method → メソッド名のスペルミス
  - SyntaxError → end の数が合っていない
  - NoMethodError → 数字に文字のメソッドを使おうとした
  - NameError → 変数名のスペルミス`,
		BrowserRules: `【Rubyのコードルール（ブラウザ実行モード）】
- puts か print で結果を表示する
- gets / gets.chomp など入力待ちは使わない（プログラムがフリーズする）
- while true や無限ループは絶対に使わない
- 標準ライブラリのみ使う（require は不要な場合は書かない）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く（# から始まる）
- 変数名はわかりやすい英語（スネークケース）にする
- 小学生（9〜12歳）でも理解できる内容にする

【クイズ・ゲームを作るときのルール】
input() を使わずに以下のパターンで作ること：
- 問題をランダムに複数生成して puts で表示し、答えも一緒に表示する
- ゲームの結果は乱数で決めて puts で表示する`,
		FileRules: `【Rubyのコードルール（ファイル実行モード）】
- puts か print で結果を表示する
- gets.chomp を使ってユーザーと対話してよい
- 標準ライブラリのみ使う（gem install は不要）
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語（スネークケース）にする`,
	},

	"php": {
		ID:   "php",
		Name: "PHP",
		Tips: `【PHPの特徴とよくあるミス】
- 変数は $ マークから始まる（例: $name = "太郎"）
- 表示は echo か print を使う
- 文（ステートメント）の最後には ; （セミコロン）が必要
- コードは <?php で始める
- 文字（文字列）は " " か ' ' で囲む
- よくあるエラー:
  - Parse error → ; や { } が足りない・多い
  - Undefined variable → $ マーク忘れ、またはスペルミス
  - Call to undefined function → 関数名のスペルミス`,
		BrowserRules: `【PHPのコードルール（ブラウザ実行モード）】
- echo か print で結果を表示する
- コードは必ず <?php から始める
- fgets(STDIN) / readline() など入力待ちは使わない（フリーズする）
- while(true) や無限ループは絶対に使わない
- 標準ライブラリのみ使う（外部ライブラリのインストール不可）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く（// か # から始まる）
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする`,
		FileRules: `【PHPのコードルール（ファイル実行モード）】
- echo か print で結果を表示する
- コードは必ず <?php から始める
- 標準ライブラリのみ使う
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
	},

	"perl": {
		ID:   "perl",
		Name: "Perl",
		Tips: `【Perlの特徴とよくあるミス】
- 変数はスカラー（$）、配列（@）、ハッシュ（%）で始まる
- 表示は print か say を使う（say は改行付き）
- 文の最後には ; （セミコロン）が必要
- 文字（文字列）は " " か ' ' で囲む（"" は変数展開あり、'' はなし）
- よくあるエラー:
  - syntax error → ; や { } の対応ミス
  - Global symbol requires explicit package name → use strict 時の変数宣言ミス（my が必要）
  - Useless use of a variable → 変数を作ったけど使っていない`,
		BrowserRules: `【Perlのコードルール（ブラウザ実行モード）】
- print か say で結果を表示する
- use strict; と use warnings; を先頭に書く
- <STDIN> など入力待ちは使わない（フリーズする）
- while(1) や無限ループは絶対に使わない
- コアモジュールのみ使う（POSIX, List::Util など）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く（# から始まる）
- 変数名はわかりやすい英語にする
- 小学生（9〜12歳）でも理解できる内容にする`,
		FileRules: `【Perlのコードルール（ファイル実行モード）】
- print か say で結果を表示する
- use strict; と use warnings; を先頭に書く
- <STDIN> を使ってユーザーと対話してよい
- コアモジュールのみ使う
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語にする`,
	},

	"rust": {
		ID:   "rust",
		Name: "Rust",
		Tips: `【Rustの特徴とよくあるミス】
- 変数は let で作る（最初は変更できない: immutable）
- 変更したいときは let mut と書く
- 表示は println!() マクロを使う（関数ではなくマクロ）
- {} は変数を埋め込むプレースホルダー（例: println!("{}", x)）
- fn main() が必ずエントリポイント
- よくあるエラー:
  - cannot borrow as mutable → mut 忘れ
  - expected type / mismatched types → 型が合っていない
  - unused variable → 変数を宣言して使っていない（_ でよい）
  - ownership moved → 所有権が移動済みの変数を使おうとした`,
		BrowserRules: `【Rustのコードルール（ブラウザ実行モード）】
- println!() マクロで結果を表示する
- fn main() を必ず書く
- std::io::stdin() など標準入力は使わない（フリーズする）
- loop や while true の無限ループは絶対に使わない
- 標準ライブラリのみ使う（Cargo.toml で外部クレート追加不可）
- 30行以内のシンプルなコードにする
- コメントは日本語で書く（// から始まる）
- 変数名はわかりやすい英語（スネークケース）にする
- 小学生（9〜12歳）でも理解できる内容にする`,
		FileRules: `【Rustのコードルール（ファイル実行モード）】
- println!() マクロで結果を表示する
- fn main() を必ず書く
- 標準ライブラリのみ使う（Cargo.toml で外部クレート追加不可）
- 50行以内のコードにする
- コメントは日本語で書く
- 変数名はわかりやすい英語（スネークケース）にする`,
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
