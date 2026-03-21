import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "go",
    label: "Go",
    icon: "rocket_launch",
    logo: "devicon-go-original-wordmark colored",
    description: "速くて軽い！サーバー向けの言語",
    tag: "本格的なサーバー開発に挑戦！",
    available: true,
  },
  outputMode: "text",
  initialCode: "// AI先生に話しかけてみよう！\n// 例：「掛け算クイズを作って」\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n",
  features: [
    {
      id: "basics", icon: "code", title: "変数・型・基本構文", ready: true,
      description: "Goの変数と型の仕組みを学ぼう",
      detail: "Goでは := で変数を作るよ（型を自動推論してくれる！）。または var x int = 10 のように型を明示することもできる。fmt.Println() で表示するのが基本。シンプルで読みやすいのがGoの大きな特徴！",
      trivia: "DockerはGoで作られてるよ！アプリを箱（コンテナ）に入れて動かす技術のDockerは、世界中の開発者が毎日使うツール。あの便利なDockerのコアがGoで書かれてるんだ！",
      examples: [
        { text: "変数を使って自己紹介を表示する", context: ":= で変数を定義して fmt.Println で表示。Goのシンプルな書き方が身につくよ！" },
        { text: "整数と小数の計算をする", context: "int と float64 の違いを体験。Goの型システムの基本が学べるよ！" },
        { text: "かけ算九九の表を作る", context: "for ループと fmt.Printf の組み合わせ。Goの書き方でスッキリ書ける！" },
        { text: "FizzBuzzプログラムを作る", context: "条件分岐と繰り返しの練習。Goのシンプルさでどれだけ短く書けるか挑戦！" },
      ],
    },
    {
      id: "function", icon: "functions", title: "関数", ready: true,
      description: "func で処理をまとめる関数を作ろう",
      detail: "Goでは func で関数を定義するよ。Goの関数は複数の値を返せるのが特徴！エラーと結果を同時に返すのがGoのスタイル。引数の型と戻り値の型を必ず書くことで、読みやすく安全なコードが書けるよ。",
      trivia: "KubernetesもGoで作られてるよ！世界中の企業がサーバー管理に使うKubernetes（k8s）はGoで書かれてる。DockerもKubernetesも、クラウド時代のインフラを支える超重要なツールがGoで作られてるんだ！",
      examples: [
        { text: "2つの数字を足す関数を作る", context: "func の基本。引数と戻り値の書き方が学べるよ！" },
        { text: "最大値と最小値を同時に返す関数を作る", context: "Goの多値返却。複数の値を返せるGoならではの機能が体験できる！" },
        { text: "フィボナッチ数列を計算する関数を作る", context: "再帰関数の書き方。Goでも再帰がシンプルに書けることが学べる！" },
        { text: "温度を摂氏から華氏に変換する関数を作る", context: "引数と戻り値の練習。実用的な計算関数の作り方が身につくよ！" },
      ],
    },
    {
      id: "control", icon: "psychology", title: "条件分岐・ループ", ready: true,
      description: "if・switch・forを使いこなそう",
      detail: "Goには if 文、switch 文、for ループがある。Goでは for だけでいろいろなループが書けるのが特徴（while がない！）。for i := 0; i < 10; i++ のCスタイルも、for i, v := range slice のrange記法も全部 for で書けるよ！",
      trivia: "Cloudflare（世界最大のCDN・セキュリティ会社）の多くのサービスがGoで動いてるよ！毎日何兆もの通信を処理するCloudflareがGoを選んだのは、Goの「速さ」と「シンプルさ」が理由。あなたのインターネット通信を守ってるのもGoかもしれない！",
      examples: [
        { text: "switch文でじゃんけんの勝ち負けを判定する", context: "switch 式の使い方。if/else より読みやすく書けるGoのswitchが学べるよ！" },
        { text: "forループで1から100の合計を求める", context: "Goの基本的な for ループ。Python の for とくらべてどうかな？" },
        { text: "rangeを使ってスライスを繰り返す", context: "range キーワードの使い方。Goらしいシンプルなループ記法が身につくよ！" },
        { text: "素数を探すプログラムを作る", context: "条件分岐とループの組み合わせ。Goのシンプルさで素数判定が書けるよ！" },
      ],
    },
    {
      id: "slice_map", icon: "list_alt", title: "スライス・マップ", ready: true,
      description: "データを管理するスライスとマップを使おう",
      detail: "Goのスライス（[]int など）はPythonのリストに似た動的配列。append() でデータを追加できる。マップ（map[string]int など）はPythonの辞書に似たキーと値のデータ構造。この2つをマスターするとGoでほとんどのデータ処理ができるよ！",
      trivia: "TerraformというクラウドインフラツールもGoで作られてるよ！AWSやGoogleCloud、Azureのサーバーを「コードで管理する」このツールは世界中の企業が使ってる。インフラエンジニアには必須のツールの心臓部がGoなんだ！",
      examples: [
        { text: "スライスに点数を追加して平均を計算する", context: "make() と append() の基本。Pythonのリストと比較しながら学べるよ！" },
        { text: "マップで英和辞書を作る", context: "map の基本的な使い方。Pythonの辞書との違いも学べるよ！" },
        { text: "スライスをソートして上位3人を表示する", context: "sort パッケージの使い方。Goの標準ライブラリの活用法が学べる！" },
        { text: "rangeでスライスの各要素を2倍にする", context: "range と スライスの組み合わせ。Goのデータ変換の基本が身につく！" },
      ],
    },
    {
      id: "struct", icon: "data_object", title: "構造体", ready: true,
      description: "structでデータをまとめて管理しよう",
      detail: "Goの struct（構造体）は、複数のデータをまとめて管理できるよ。Pythonのクラスに似てるけど、Goにはclassがなくて struct + method の組み合わせでオブジェクト指向的なプログラムが書ける。Goらしいシンプルなデータ設計が学べるよ！",
      trivia: "Goを作ったのはGoogleの伝説的なエンジニアたちだよ！Rob Pike、Ken Thompson（UnixとC言語を作った人！）、Robert Griesemerの3人が2009年に作ったんだ。「複雑すぎるC++に疲れたから」シンプルな言語を作ろうとしたのが始まり。本当に面白い話だよね！",
      examples: [
        { text: "RPGキャラクターのデータを管理する", context: "struct の基本。名前・HP・攻撃力などをまとめて管理する方法が学べるよ！" },
        { text: "構造体にメソッドを追加する", context: "func (r *Robot) で構造体にメソッドを追加。Goのオブジェクト指向の書き方が学べる！" },
        { text: "学生の成績管理システムを作る", context: "スライスと構造体の組み合わせ。データベースの基本的な考え方が身につくよ！" },
        { text: "動物の構造体を作って鳴き声を管理する", context: "構造体とメソッドの総合練習。Goでオブジェクト指向プログラミングの基礎が学べる！" },
      ],
    },
  ],
  installGuide: {
    windows: {
      steps: [
        { title: "Go公式サイトからダウンロードする", detail: "「Go言語 ダウンロード」と検索して、公式サイト（go.dev）からWindows用のインストーラー（.msi）をダウンロードしよう。" },
        { title: "インストールする", detail: "ダウンロードした .msi ファイルを実行して、画面の指示に従ってインストールしよう。デフォルト設定でOKだよ。" },
        { title: "インストール確認", detail: "コマンドプロンプトで「go version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    mac: {
      steps: [
        { title: "Homebrewでインストールする", detail: "ターミナルで「brew install go」と入力してEnterキーを押そう。Homebrewが自動的にインストールしてくれるよ。" },
        { title: "インストール確認", detail: "ターミナルで「go version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    linux: {
      steps: [
        { title: "aptでインストールする", detail: "ターミナルで「sudo apt-get install golang-go」と入力してEnterキーを押そう。" },
        { title: "インストール確認", detail: "ターミナルで「go version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
  },
};

export default config;
