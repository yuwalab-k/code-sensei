import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "ruby",
    label: "Ruby",
    icon: "diamond",
    logo: "devicon-ruby-plain colored",
    description: "読みやすくて書きやすい！日本人が作った言語",
    tag: "きれいなコードを書きたい人に！",
    available: true,
  },
  outputMode: "text",
  initialCode: "# AI先生に話しかけてみよう！\n# 例：「掛け算クイズを作って」\n",
  features: [
    {
      id: "calc", icon: "calculate", title: "計算・数学", ready: true,
      description: "数字を使った計算やゲームが作れる",
      detail: "Rubyは数学計算がシンプルに書けるよ。rand()で乱数を作ったり、.times で繰り返したりできる。Pythonや JavaScript と比べると、より「英語っぽく」読めるコードになるのがRubyの特徴！",
      trivia: "RubyはTwitterの最初のバージョンを作るのに使われた言語だよ！2006年にたった数人で作られたTwitterが、今では世界中で使われてる。あの大きなサービスの土台がRubyで書かれてたんだ！",
      examples: [
        { text: "かけ算九九の表を作る", context: "timesメソッドを使うとすっきり書ける！Pythonと比べながら書くとRubyらしい書き方が身につくよ。" },
        { text: "素数を探すプログラムを作る", context: "素数はインターネットの暗号技術の基本。Rubyのシンプルな書き方で挑戦しよう！" },
        { text: "ランダムな数当てゲームを作る", context: "rand()で乱数を作る方法が学べる。シンプルなコードでゲームが作れるよ！" },
        { text: "フィボナッチ数列を表示する", context: "ひまわりの種の並び方や巻き貝の形にも現れる不思議な数列。Rubyのシンプルさを体感しよう！" },
      ],
    },
    {
      id: "string", icon: "edit_note", title: "文字列処理", ready: true,
      description: "文字を加工したり、暗号を作ったりできる",
      detail: "Rubyは文字列処理がとても得意！.reverse で逆さにしたり、.upcase で大文字にしたり、.include? で文字の検索ができるよ。メソッドが英語の文章みたいで読みやすいのがRubyの魅力！",
      trivia: "GitHubはRubyで作られたよ！世界中の3億以上のプログラムコードが、Rubyで作られたGitHubに保存されてる。今君が使ってるコードの仕組みもGitHubで管理されてることが多いんだよ！",
      examples: [
        { text: "文字を逆さにして回文チェッカーを作る", context: ".reverseメソッド一発で逆さにできる！Rubyのシンプルさを実感できるよ。" },
        { text: "文章に含まれる単語を数える", context: ".split で単語に分割、.length で数える。Rubyらしいチェーンメソッドの書き方が学べるよ！" },
        { text: "名前を大文字・小文字に変換するツールを作る", context: ".upcase と .downcase と .capitalize を使い分けよう。文字列操作の基本がまとめて学べるよ！" },
        { text: "シーザー暗号で文字を暗号化する", context: ".ord と .chr を使った暗号化。Rubyのエレガントな書き方でスパイ暗号が作れる！" },
      ],
    },
    {
      id: "array", icon: "list_alt", title: "配列・ハッシュ", ready: true,
      description: "配列やハッシュでデータをまとめて管理しよう",
      detail: "Rubyの配列とハッシュはとても使いやすい！.each でループしたり、.map で変換したり、.select で絞り込んだりできるよ。Pythonのリストや辞書に似てるけど、もっと読みやすく書けるのがRubyの特徴！",
      trivia: "Shopify（世界最大級のECプラットフォーム）もRubyで作られてるよ！世界中の200万以上のお店がShopifyを使ってる。あのオンラインショップの裏側でRubyが動いてるんだ！",
      examples: [
        { text: "点数リストから合格者だけ抜き出す", context: ".select メソッドでデータを絞り込む。英語みたいに読めるRubyらしい書き方が学べるよ！" },
        { text: "名前リストをランダムにシャッフルする", context: ".shuffle メソッドで一発シャッフル！くじ引きアプリや席替えアプリの基本が学べる。" },
        { text: "ショッピングカートの合計金額を計算する", context: ".sum や .inject で配列の合計を計算。Rubyのシンプルさを実感できるよ！" },
        { text: "テストの点数を集計して統計を出す", context: ".min .max .sum を使って最小・最大・合計を一発で！データ分析の基本が学べるよ。" },
      ],
    },
    {
      id: "iterator", icon: "loop", title: "繰り返し・ブロック", ready: true,
      description: "Rubyらしい繰り返しの書き方を学ぼう",
      detail: "Rubyには .times / .upto / .each など、英語の文章みたいに読める繰り返し方法がたくさんあるよ。ブロック（do...end）を使うと、複雑な処理もシンプルに書けるのがRubyの大きな特徴！",
      trivia: "Ruby を作った「まつもとゆきひろ（Matz）」さんは日本人だよ！1993年にプログラマーを幸せにするために作り始めたんだって。「プログラマーの幸せを大切にする」という考え方が、Rubyのシンプルで読みやすい書き方につながってるんだ！",
      examples: [
        { text: "5.times で5回メッセージを表示する", context: "Rubyらしい書き方の入り口！英語みたいに「5回繰り返す」が1行で書けるよ。" },
        { text: "1から10まで順番に表示する", context: "(1..10).each という書き方でforループが書ける。Rubyの「範囲」という便利な概念が学べるよ！" },
        { text: "リストの各要素を2倍にする", context: ".map ブロックでリスト全体を変換。Pythonのリスト内包表記と比べてみよう！" },
        { text: "FizzBuzzをRubyらしく書く", context: "条件分岐と繰り返しの組み合わせ。Rubyのシンプルさでどれくらい短く書けるか挑戦してみよう！" },
      ],
    },
    {
      id: "method", icon: "functions", title: "メソッド・クラス", ready: true,
      description: "自分だけのメソッドやクラスを作ろう",
      detail: "Rubyの def でメソッドを定義すると、処理をまとめて再利用できるよ。classを使えばゲームのキャラクターや動物のデータをオブジェクトとして管理できる。Rubyはオブジェクト指向がとても自然に書けるのが特徴！",
      trivia: "Rubyでは数字も文字も配列もすべてがオブジェクトだよ！たとえば 3.times と書けるのは、3という数字もオブジェクトだから。「すべてがオブジェクト」という考え方は、プログラミングの深いところを理解するのにとても役立つんだ！",
      examples: [
        { text: "挨拶するメソッドを作る", context: "def で自分だけの関数が作れる。同じコードを何度も書かなくてよくなる「再利用」の考え方が学べるよ！" },
        { text: "動物クラスを作って鳴き声を管理する", context: "クラスとインスタンスの基本。RPGのキャラクター管理や、動物図鑑アプリと同じ仕組みが学べる！" },
        { text: "RPGキャラクターのステータスを管理する", context: "attr_accessor を使うとゲームのキャラクターデータが簡単に作れる。Rubyのクラスの便利さが実感できるよ！" },
        { text: "電卓クラスを作る", context: "クラスにメソッドをまとめる実践練習。本格的なアプリ開発の基礎となるオブジェクト指向が身につく！" },
      ],
    },
  ],
};

export default config;
