import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "javascript",
    label: "JavaScript",
    icon: "javascript",
    logo: "devicon-javascript-plain colored",
    description: "Webを動かす言語！計算やデータ処理もできる",
    tag: "Pythonの次のステップに！",
    available: true,
  },
  outputMode: "text",
  initialCode:
    "// AI先生に話しかけてみよう！\n// 例：「掛け算クイズを作って」\n",
  features: [
    {
      id: "calc", icon: "calculate", title: "計算・数学", ready: true,
      description: "数字を使った計算やゲームが作れる",
      detail: "JavaScriptは数学計算が得意！Math.random()で乱数を作ったり、Math.floor()で小数点を切り捨てたりできるよ。数字を使ったゲームや暗号も作れる。",
      trivia: "JavaScriptは最初わずか10日間で作られた言語なんだよ！1995年に作られて、今では世界中の97%以上のWebサイトで使われてる。10日で作ったとは思えないくらいすごい言語になったね！",
      examples: [
        { text: "かけ算九九の表を作る", context: "forループと文字列結合で九九表が完成！Pythonと比べながら書くとJavaScriptの書き方が身につくよ。" },
        { text: "素数を探すプログラムを作る", context: "素数はインターネットの暗号技術の基本。スマホで安全にお買い物できるのも素数のおかげ！" },
        { text: "ランダムな数当てゲームを作る", context: "Math.random()で乱数を作る方法が学べる。コンピューターがランダムに数を選ぶ仕組みの基本！" },
        { text: "フィボナッチ数列を表示する", context: "ひまわりの種の並び方や巻き貝の形にも現れる不思議な数列。自然界の数学をコードで再現できるよ！" },
      ],
    },
    {
      id: "string", icon: "edit_note", title: "文字列処理", ready: true,
      description: "文字を加工したり、暗号を作ったりできる",
      detail: "JavaScriptは文字列を自由に操作できるよ。split()で文字を分割したり、join()でつなげたり、includes()で文字を検索したりできる。テンプレートリテラルも便利！",
      trivia: "GoogleやYahooの検索機能も、JavaScriptの文字列処理技術が使われてるよ。あなたが検索ボックスに文字を打つたびに、JavaScriptが文字を解析して候補を出してるんだ！",
      examples: [
        { text: "文字を逆さにして回文チェッカーを作る", context: "split('')とreverse()とjoin('')を組み合わせると1行で書けるよ！Pythonとはちょっと違う書き方が面白い。" },
        { text: "文章に含まれる単語を数える", context: "Webサービスの文字数カウント機能の仕組みと同じ。SNSの140文字制限もこの技術で実現してるよ！" },
        { text: "テンプレートで自己紹介文を作る", context: "テンプレートリテラル（バッククォート）を使うと変数を文章に埋め込める。メールの自動作成機能と同じ仕組み！" },
        { text: "シーザー暗号で文字を暗号化する", context: "charCodeAt()とString.fromCharCode()を使った暗号化。スパイ映画みたいな暗号が作れるよ！" },
      ],
    },
    {
      id: "array", icon: "list_alt", title: "配列・データ管理", ready: true,
      description: "配列を使って複数のデータを管理しよう",
      detail: "JavaScriptの配列はとても強力！map()・filter()・reduce()など、データを変換・絞り込む便利なメソッドがたくさん。Pythonのリストと比べながら学べるよ。",
      trivia: "SpotifyやNetflixの「おすすめ機能」もJavaScriptの配列処理が基本。あなたの再生履歴がJavaScriptの配列で管理されて、AIがおすすめを計算してるんだよ！",
      examples: [
        { text: "点数リストから合格者だけ抜き出す", context: "filter()メソッドでデータを絞り込む。SNSのフィルタリング機能やECサイトの絞り込み検索と同じ仕組み！" },
        { text: "名前リストをランダムにシャッフルする", context: "Math.random()を使った並べ替えアルゴリズム。くじ引きアプリや席替えアプリの基本！" },
        { text: "ショッピングカートの合計金額を計算する", context: "reduce()で配列の合計を計算。Amazonの買い物カゴと同じ仕組みが学べるよ！" },
        { text: "テストの点数を集計して統計を出す", context: "平均・最大・最小を求める方法。データ分析の基本で、ゲームのスコア管理にも使える！" },
      ],
    },
    {
      id: "logic", icon: "psychology", title: "ロジック・アルゴリズム", ready: true,
      description: "条件分岐やループで問題を解くプログラムを作ろう",
      detail: "if文・for文・while文を組み合わせると、複雑な問題も解けるよ。ゲームのルール判定、スコア計算、探索アルゴリズムなど、プログラムの「考える力」を鍛えよう！",
      trivia: "Googleマップが「最短ルート」を見つけるのも、JavaScriptで書かれたアルゴリズムが動いてるよ。出発地から目的地まで何億通りものルートを瞬時に計算してるんだ！",
      examples: [
        { text: "FizzBuzzプログラムを作る", context: "3の倍数で「Fizz」、5の倍数で「Buzz」を表示。世界中のプログラミング面接で出る有名な問題！" },
        { text: "じゃんけんの勝ち負けを判定する", context: "条件分岐の練習に最適！コンピューターとじゃんけんして結果を表示するシンプルなゲームが作れるよ。" },
        { text: "バブルソートで数字を並べ替える", context: "「比べて入れ替える」を繰り返すソートの基本アルゴリズム。コンピューター科学の重要な概念が学べる！" },
        { text: "再帰関数で階乗を計算する", context: "関数が自分自身を呼び出す「再帰」という考え方。AIや数学の計算に使われる強力なテクニック！" },
      ],
    },
    {
      id: "object", icon: "data_object", title: "オブジェクト", ready: true,
      description: "データをまとめて管理するオブジェクトを使おう",
      detail: "JavaScriptのオブジェクトはPythonの辞書（dict）に似てるよ。キーと値のペアでデータを整理できる。クラスを使えばゲームのキャラクターや商品管理など、本格的なプログラムが作れる！",
      trivia: "WebサービスのAPI（データのやり取り）はほぼすべてJSONという形式を使ってるよ。JSONはJavaScript Object Notationの略で、JavaScriptのオブジェクトから生まれた世界共通のデータ形式なんだ！",
      examples: [
        { text: "RPGキャラクターのデータを管理する", context: "名前・HP・攻撃力などをオブジェクトで管理。ゲームのキャラクターデータの仕組みがそのまま学べるよ！" },
        { text: "学生の成績管理システムを作る", context: "複数の学生データをオブジェクトの配列で管理。実際の成績管理アプリと同じデータ構造が学べる！" },
        { text: "簡単なショップの在庫管理を作る", context: "商品・価格・在庫数をオブジェクトで管理。ECサイトのデータベース設計の考え方が身につくよ！" },
        { text: "クラスを使ってコインを実装する", context: "クラスとメソッドの基本。オブジェクト指向プログラミングの入り口で、大きなプログラムを作る力が身につく！" },
      ],
    },
  ],
};

export default config;
