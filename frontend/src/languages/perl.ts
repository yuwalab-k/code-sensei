import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "perl",
    label: "Perl",
    icon: "auto_fix_high",
    logo: "devicon-perl-plain colored",
    description: "テキスト処理が得意な言語",
    tag: "文字や文章を自在に操りたい人に！",
    available: true,
  },
  outputMode: "text",
  initialCode: "# AI先生に話しかけてみよう！\n# 例：「掛け算クイズを作って」\nuse strict;\nuse warnings;\n",
  features: [
    {
      id: "calc", icon: "calculate", title: "計算・数学", ready: true,
      description: "数字を使った計算やゲームが作れる",
      detail: "Perlは数学計算もシンプルに書けるよ。int(rand(N)) で乱数を作ったり、** で累乗を計算したりできる。print で結果を表示するのが基本！数値も文字列も自動的に変換してくれる柔軟さがPerlの特徴。",
      trivia: "IMDbという映画データベースサイトは最初Perlで作られたよ！1990年代、世界最大の映画情報サイトの土台にPerlが使われてた。あの膨大な映画データの処理をPerlがこなしてたんだ！",
      examples: [
        { text: "かけ算九九の表を作る", context: "for ループと print で九九表が完成！Perlの基本的な繰り返し処理が学べるよ。" },
        { text: "素数を探すプログラムを作る", context: "Perlのシンプルな条件分岐で素数判定。数学とプログラミングの融合が楽しめる！" },
        { text: "ランダムな数当てゲームを作る", context: "rand() で乱数を作る方法が学べる。Perlで簡単なゲームを作ってみよう！" },
        { text: "フィボナッチ数列を表示する", context: "自然界の不思議な数列をPerlで再現。繰り返し処理の基本が身につくよ！" },
      ],
    },
    {
      id: "string", icon: "edit_note", title: "文字列処理", ready: true,
      description: "文字を自在に操る強力な処理ができる",
      detail: "Perlは文字列処理の王様！ . でつなげたり、length() で長さを測ったり、uc() で大文字にしたり、reverse で逆さにしたりできるよ。特にテキストを加工する処理はPerlが世界一得意な領域！",
      trivia: "Perlを作ったラリー・ウォール（Larry Wall）は言語学者でもあったんだよ！プログラミング言語を設計するのに、人間の言語の研究が役立ったんだ。「プログラマーの怠惰は美徳」という有名な言葉もラリーが言ったんだよ！",
      examples: [
        { text: "文字を逆さにして回文チェッカーを作る", context: "reverse で文字列を逆さにするのが超簡単！Perlらしいシンプルさが実感できるよ。" },
        { text: "文章に含まれる単語を数える", context: "split で単語に分割して scalar で数える。Perlの得意分野が体験できるよ！" },
        { text: "名前を大文字・小文字に変換するツールを作る", context: "uc() / lc() / ucfirst() の使い分けを練習しよう！" },
        { text: "ROT13でテキストを変換する", context: "アルファベットを13文字ずらす古典的な暗号。Perlではtr///演算子で一発変換できちゃう！" },
      ],
    },
    {
      id: "regex", icon: "search", title: "正規表現", ready: true,
      description: "パターンで文字列を自在に検索・置換しよう",
      detail: "Perlは正規表現の生みの親ともいえる言語！m// でパターン検索、s/// で置換が直感的に書ける。正規表現はメールアドレスのチェックや、文章からのデータ抽出など、多くのプログラムで使われる強力な技術だよ。",
      trivia: "今使われている多くのプログラミング言語の正規表現は、Perlの正規表現を参考にして作られてるよ！Python、Java、JavaScript…みんなPerlから学んだんだ。Perlが「正規表現の標準」を作ったといえるくらい影響力があるんだよ！",
      examples: [
        { text: "メールアドレスの形式をチェックする", context: "m// でパターンマッチング。フォームのバリデーション（入力チェック）と同じ技術が学べる！" },
        { text: "文章から数字だけを取り出す", context: "\\d+ で数字のパターンを検索。ログファイルからデータを取り出す技術の基礎！" },
        { text: "文章中の特定の単語を別の言葉に置き換える", context: "s/// で置換。テキストエディタの「検索と置換」と同じ仕組みが学べるよ！" },
        { text: "HTMLタグを取り除いてテキストだけを取り出す", context: "ウェブスクレイピングの基礎。Perlが得意とする実用的なテキスト処理が体験できる！" },
      ],
    },
    {
      id: "array", icon: "list_alt", title: "配列・ハッシュ", ready: true,
      description: "リストや辞書型でデータを管理しよう",
      detail: "Perlの配列（@）とハッシュ（%）はとても使いやすい！push/pop でデータを追加・削除したり、sort で並べ替えたり、keys/values でハッシュの内容を取り出したりできるよ。データ管理の基本が学べる！",
      trivia: "バイオインフォマティクス（生命情報科学）の世界ではPerlが長らく主役だったよ！人間のゲノム（遺伝子情報）を解析するプログラムも多くがPerlで書かれてた。あなたの体の設計図の解読にも使われた言語なんだ！",
      examples: [
        { text: "点数リストから合格者だけ抜き出す", context: "grep でデータを絞り込む。PerlのgrepはLinuxコマンドのgrepと同じ考え方！" },
        { text: "名前リストをランダムにシャッフルする", context: "List::Util の shuffle でシャッフル。くじ引きアプリの基本が学べるよ！" },
        { text: "ハッシュで英和辞書を作る", context: "キーと値のペアでデータを管理。プログラムの辞書機能の基礎が学べるよ！" },
        { text: "テストの点数を集計して統計を出す", context: "配列のデータを集計・分析。Perlのデータ処理力が実感できるよ！" },
      ],
    },
    {
      id: "sub", icon: "functions", title: "サブルーチン", ready: true,
      description: "処理をまとめるサブルーチンを作ろう",
      detail: "Perlでは sub（サブルーチン）を使って処理をまとめられるよ。@_ で引数を受け取り、return で値を返す。繰り返し使う処理をサブルーチンにまとめるのがプログラミングの基本スキル！",
      trivia: "Perlには「やり方は一つじゃない（There's More Than One Way To Do It）」という哲学があるよ！略してTMTOWTDI（タームトウイッティ）と呼ばれるこの考え方は、プログラマーに自由な発想を与えてくれるんだ。PythonとはちょっとちがうPerlのユニークな哲学だよ！",
      examples: [
        { text: "挨拶するサブルーチンを作る", context: "sub の基本。処理をまとめて再利用する「プログラミングの基本」が学べるよ！" },
        { text: "最大値を求めるサブルーチンを作る", context: "@_ で引数を受け取る方法。Perlならではのサブルーチンの書き方が学べる！" },
        { text: "リストを合計するサブルーチンを作る", context: "reduce の考え方。データを集約する処理の基礎が学べるよ！" },
        { text: "名前と年齢を受け取って自己紹介するサブルーチン", context: "引数と文字列処理の組み合わせ。実用的なサブルーチンの作り方が学べる！" },
      ],
    },
  ],
  installGuide: {
    windows: {
      steps: [
        { title: "Strawberry Perlをダウンロードする", detail: "「Strawberry Perl Windows」と検索して、公式サイトからダウンロードしよう。Windowsで一番使いやすいPerlだよ。" },
        { title: "インストールする", detail: "ダウンロードしたファイルを実行して、画面の指示に従ってインストールしよう。" },
        { title: "インストール確認", detail: "コマンドプロンプトで「perl -v」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    mac: {
      steps: [
        { title: "最初から入ってるよ！", detail: "Macには最初からPerlが入ってることが多いよ。ターミナルで「perl -v」と入力して確認してみよう。" },
        { title: "入っていない場合はHomebrewでインストール", detail: "ターミナルで「brew install perl」と入力してEnterキーを押そう。" },
      ],
    },
    linux: {
      steps: [
        { title: "最初から入ってることが多い", detail: "多くのLinuxにはPerlが最初から入ってるよ。ターミナルで「perl -v」と入力して確認しよう。" },
        { title: "入っていない場合はaptでインストール", detail: "ターミナルで「sudo apt-get install perl」と入力してEnterキーを押そう。" },
      ],
    },
  },
};

export default config;
