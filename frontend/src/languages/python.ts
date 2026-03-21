import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "python",
    label: "Python",
    icon: "terminal",
    logo: "devicon-python-plain colored",
    description: "計算やデータ処理が得意！AIでも使われる人気の言語",
    tag: "はじめての人におすすめ",
    available: true,
    usedFor: ["ウェブ開発", "AI・機械学習", "データ分析", "自動化・スクリプト"],
  },
  outputMode: "text",
  initialCode: "# AI先生に話しかけてみよう！\n# 例：「掛け算の問題を出してくれるプログラムを作って」\n",
  features: [
    {
      id: "calc", icon: "calculate", title: "計算・数学", ready: true,
      description: "素数探しや暗号など、数字で遊ぼう",
      detail: "Pythonは計算がとても得意！数学の公式をそのままコードにできるよ。数字を使ったゲームや、秘密の暗号を作ることもできる。",
      trivia: "Pythonの名前はヘビじゃなくて、イギリスのコメディグループ「Monty Python」から来てるよ！作者がファンだったんだって。プログラマーはユーモアが好きなんだ。",
      examples: [
        { text: "かけ算九九の表を作る", context: "電卓アプリや計算ゲームの仕組みが学べる。コンピューターが計算する基本がわかるよ！" },
        { text: "素数を探すプログラムを作る", context: "素数はインターネットのパスワード保護（暗号）に使われてる技術。スマホで安全にお買い物できるのも素数のおかげ！" },
        { text: "シーザー暗号で文字を暗号にする", context: "古代ローマの将軍ユリウス・カエサルが実際に使ってた暗号。2000年前のスパイ技術がコードで再現できる！" },
        { text: "フィボナッチ数列を表示する", context: "ひまわりの種の並び方や巻き貝の形にも現れる不思議な数のパターン。自然界の数学が体感できるよ！" },
      ],
    },
    {
      id: "text", icon: "edit_note", title: "テキスト処理", ready: true,
      description: "文字を加工したり、暗号を作ったりできる",
      detail: "文字を自由に加工できるのもPythonの強み。文字を逆さにしたり、単語を数えたり、名前からイニシャルを取り出したりできるよ。",
      trivia: "Googleが検索結果を表示するときも、Pythonのテキスト処理技術が使われてるよ。1秒間に何十億もの文字を処理してるんだ！",
      examples: [
        { text: "入力した文字を逆さにして回文チェッカーを作る", context: "「しんぶんし」「たけやぶやけた」みたいな回文を自動で見つけられるよ！" },
        { text: "文章の中に同じ言葉が何回出るか数える", context: "ニュース記事の分析やSNSのトレンド調査でも使われる技術。AIが文章を理解する仕組みの基本でもあるよ！" },
        { text: "名前からイニシャルを取り出すツールを作る", context: "スポーツのユニフォームや名刺づくりアプリと同じ仕組み！名前を入力するだけで自動生成できる。" },
        { text: "ランダムな名言ジェネレーターを作る", context: "朝起きたらランダムな名言が表示されるアプリの仕組み。好きな言葉リストから自動で選ぶよ！" },
      ],
    },
    {
      id: "game", icon: "sports_esports", title: "テキストゲーム", ready: true,
      description: "文字だけで遊べるゲームが作れる",
      detail: "画面がなくてもゲームは作れる！数当てゲームや、じゃんけん、クイズゲームなど、printとifとrangeだけで本格的なゲームが完成するよ。",
      trivia: "世界で初めて広まったコンピューターゲーム「テキストアドベンチャー」は1970年代に文字だけで作られた。ポケモンも最初は文字だけのゲームから発想されたんだよ！",
      examples: [
        { text: "1〜100の数当てゲームを作る", context: "「もっと大きい/小さい」のヒントで絞り込む。二分探索というアルゴリズムの考え方が自然と身につく！" },
        { text: "じゃんけんゲームを作る", context: "randomとifを使うだけで完成！コンピューターがランダムに手を選ぶ仕組みが学べるよ。" },
        { text: "クイズゲームを作る", context: "問題と答えをリストで管理する方法が学べる。問題数を増やすだけでどんどん発展できる！" },
        { text: "タイピングスピードテストを作る", context: "入力した時間を計測する仕組み。ゲームのタイムアタック機能の基本が学べる！" },
      ],
    },
    {
      id: "datetime", icon: "schedule", title: "日付・時間", ready: true,
      description: "今の時間や日付を使ったプログラムが作れる",
      detail: "Pythonのdatetimeライブラリを使うと、今の時刻・曜日・誕生日までの日数など、時間に関するあらゆる処理ができるよ。",
      trivia: "国際宇宙ステーション（ISS）の軌道計算にもPythonが使われてるよ。宇宙飛行士の帰還日時を計算するのにも同じ技術！時間の計算って宇宙でも大切なんだね。",
      examples: [
        { text: "今日の曜日と日付を表示する", context: "アプリのホーム画面によく表示されてるあれ！datetimeライブラリで1行で取得できるよ。" },
        { text: "誕生日までの日数カウントダウンを作る", context: "イベントアプリのカウントダウン機能と同じ仕組み。「あと○日！」が自動で計算できる！" },
        { text: "今が朝・昼・夜かによってメッセージを変える", context: "スマホのおはようございます通知の仕組みと同じ！時間帯によって違う動作をするプログラムが作れる。" },
        { text: "何秒でボタンを押せるか計測するゲーム", context: "タイマーアプリやゲームの反応速度テストの基本。時間の計測方法が学べるよ！" },
      ],
    },
    {
      id: "list", icon: "list_alt", title: "リスト・データ管理", ready: true,
      description: "リストを使って複数のデータを管理しよう",
      detail: "リスト（list）を使うと、たくさんのデータをまとめて管理できるよ。ランキング、名簿、ショッピングリストなど、現実のアプリでも必ず使われる基本技術！",
      trivia: "SpotifyやNetflixの「おすすめ機能」もリストとデータ管理が基本。あなたの再生リストや視聴履歴がリストで管理されて、AIがおすすめを計算してるんだよ！",
      examples: [
        { text: "ショッピングリストアプリを作る", context: "追加・削除・表示の操作が学べる。これがデータベースアプリの基本になるよ！" },
        { text: "クラスのテストの点数から平均点を計算する", context: "成績管理システムや家計簿アプリと同じ仕組み。たくさんの数字をまとめて処理できる！" },
        { text: "名前リストをランダムに並べ替える", context: "体育の班分けや席替えアプリに使える！random.shuffle()という便利な関数が学べるよ。" },
        { text: "好きなものランキングを集計する", context: "投票アプリやアンケートシステムの仕組み。辞書（dict）という強力なデータ構造が学べる！" },
      ],
    },
    {
      id: "graph", icon: "bar_chart", title: "グラフ作成", ready: true,
      description: "データを棒グラフや円グラフで見える化",
      detail: "matplotlibというライブラリを使うと、数字のデータをきれいなグラフにできるよ。好きな食べ物アンケートの結果や、気温の変化なども視覚化できる。",
      trivia: "NASAもmatplotlibを使ってデータを分析してるよ！火星探査機から送られてくる膨大なデータをグラフで見える化して、科学者たちが研究に使ってるんだ。",
      examples: [
        { text: "好きな食べ物ランキングを棒グラフにする", context: "学校のアンケート結果を発表するときに使える！グラフで見ると一目でわかるよね。" },
        { text: "1週間の気温変化を折れ線グラフにする", context: "天気予報アプリが気温を表示する仕組みと同じ。データを「見える化」する力が身につく！" },
        { text: "クラスの誕生月を円グラフにする", context: "新聞や会社の発表でよく使われる表現方法。データを整理して伝える力が学べるよ！" },
      ],
    },
    {
      id: "image", icon: "image", title: "画像処理", ready: true,
      description: "写真にフィルターをかけたり加工できる",
      detail: "Pillowというライブラリを使うと、写真をグレースケールにしたり、明るさを変えたり、モザイクをかけたりできるよ。自分だけのフィルターが作れる！",
      trivia: "Instagramのフィルター機能もPythonで作られてるよ！毎日1億枚以上の写真が処理されてる。スマホで何気なく使ってるフィルターの裏側にPythonがあるんだ。",
      examples: [
        { text: "写真を白黒に変換するフィルターを作る", context: "スマホのカメラアプリにある「モノクロフィルター」と同じ仕組み。自分だけのフィルターアプリが作れるよ！" },
        { text: "画像を回転・拡大縮小する", context: "ゲームのキャラクター画像やスタンプを動かす仕組みと同じ技術！" },
        { text: "写真にテキストを書き込むツールを作る", context: "ミームジェネレーターや、写真に日付を入れるアプリの仕組みと同じだよ！" },
      ],
    },
    {
      id: "file", icon: "folder_open", title: "ファイル操作", ready: true,
      description: "テキストファイルを作ったり読んだりできる",
      detail: "Pythonはファイルの読み書きが得意。日記をファイルに保存したり、メモ帳アプリを作ったりできるよ。保存したデータを次に読み込むことも簡単！",
      trivia: "世界中で毎秒何百万もの.txtや.csvファイルがPythonで処理されてるよ。銀行の取引記録や病院のカルテ管理にもPythonのファイル操作が使われてるんだ！",
      examples: [
        { text: "日記をテキストファイルに保存する", context: "ゲームのセーブ機能と同じ仕組み！データを保存して次回に読み込むことが学べるよ。" },
        { text: "ファイルから文章を読み込んで表示する", context: "電子書籍リーダーやゲームのシナリオ表示の仕組みと同じ。読み込んだ内容を加工することもできる！" },
        { text: "名前リストをファイルに書き出す", context: "運動会の参加者リストや、ゲームのランキングボードを保存する仕組みと同じだよ！" },
      ],
    },
  ],
  frameworks: [
    {
      name: "Django",
      logo: "devicon-django-plain colored",
      description: "フルスタックウェブ開発の定番フレームワーク",
      detail: "「ログイン機能・管理画面・データベース・セキュリティ対策が全部最初から入ってる」フルスタックフレームワーク。一からウェブサービスを作るなら Django が一番早い。Instagram・Pinterest・Disqus も最初は Django で作られた。一方で「Djangoのやり方」に沿う必要があるため自由度は低め。「とにかく速く本格的なサービスを作りたい」ときに選ぶ。",
    },
    {
      name: "Flask",
      logo: "devicon-flask-original colored",
      description: "シンプルで軽量なウェブ API 開発向けフレームワーク",
      detail: "最小限の機能だけを持つ軽量フレームワーク。「ちょっとした API を作りたい」「機械学習モデルをウェブで公開したい」ときに向いている。Django より学習コストが低く、コードの自由度も高い。ただし認証やデータベースなどは自分で追加する必要があるため、大規模サービスには向かない。「シンプルに始めたい・AIモデルをAPIとして公開したい」ときに選ぶ。",
    },
    {
      name: "FastAPI",
      logo: "devicon-fastapi-plain colored",
      description: "高速な API サーバー開発フレームワーク",
      detail: "Flask の書きやすさを残しつつ、速度を Node.js や Go レベルまで引き上げたフレームワーク。コードを書くだけで API ドキュメントが自動生成されるため、フロントエンドエンジニアとの連携がしやすい。近年 AI サービスのバックエンドとして急速に普及中。「API の速度を重視したい・ドキュメントを自動生成したい」なら FastAPI が現在のベストチョイス。",
    },
    {
      name: "NumPy",
      logo: "devicon-numpy-original colored",
      description: "科学技術計算・数値処理の標準ライブラリ",
      detail: "行列やベクトルの計算を C 言語並みの速度で処理できるライブラリ。pandas・TensorFlow・PyTorch など Python の数値計算ライブラリはほぼすべて NumPy の上に作られている。「データ分析・AI・物理シミュレーション・信号処理」など数値を扱う用途なら必ずお世話になる。Python で数学や科学の計算をするなら避けられない基礎ライブラリ。",
    },
    {
      name: "pandas",
      logo: "devicon-pandas-original colored",
      description: "データ分析・データ処理の王道ライブラリ",
      detail: "Excel のような表形式のデータを Python で自由に扱えるライブラリ。CSV の読み込み・集計・グラフ化・データのクレンジングが得意で、データサイエンティストの日常業務には欠かせない。「大量のデータを分析したい・Excelで限界を感じてきた」なら pandas が最初のステップ。NumPy と組み合わせて使うことが多い。",
    },
    {
      name: "TensorFlow",
      logo: "devicon-tensorflow-original colored",
      description: "Google が作った AI・機械学習フレームワーク",
      detail: "Google が開発した業界標準の機械学習フレームワーク。スマホ向け（TensorFlow Lite）やブラウザ向け（TensorFlow.js）も揃っていて、研究から本番環境まで一気通貫で使える。企業での AI 導入実績が多く「仕事で AI を使うプロジェクトに関わりたい」なら覚えておきたい。PyTorch に比べると本番デプロイのサポートが手厚い。",
    },
    {
      name: "PyTorch",
      logo: "devicon-pytorch-original colored",
      description: "AI 研究・ディープラーニングで世界シェア No.1",
      detail: "Meta が作った AI 研究向けフレームワーク。「書いたコードがそのまま直感的に動く」設計が研究者に大人気で、AI 論文の 70% 以上が PyTorch を使用。ChatGPT を作った OpenAI もPyTorch を使っている。TensorFlow と比べると研究・プロトタイプ開発向きで、「最新の AI 技術を試したい・研究がしたい」なら PyTorch が第一選択。",
    },
  ],
  installGuide: {
    windows: {
      steps: [
        { title: "python.org を開く", detail: "ブラウザで「python.org」と検索して開こう" },
        { title: "ダウンロードページへ", detail: "「Downloads」→「Windows」をクリック" },
        { title: "インストーラーを実行", detail: "ダウンロードした .exe ファイルをダブルクリック", warn: "「Add Python to PATH」に必ずチェックを入れてね！" },
        { title: "インストール完了", detail: "「Install Now」をクリックして待つ" },
        { title: "動作確認", detail: "コマンドプロンプトを開いて「python --version」と入力してみよう" },
      ],
    },
    mac: {
      steps: [
        { title: "python.org を開く", detail: "ブラウザで「python.org」と検索して開こう" },
        { title: "ダウンロードページへ", detail: "「Downloads」→「macOS」をクリック" },
        { title: "pkg ファイルを実行", detail: "ダウンロードした .pkg ファイルをダブルクリックしてインストール" },
        { title: "動作確認", detail: "Terminal（ターミナル）を開いて「python3 --version」と入力してみよう" },
        { title: ".py ファイルを実行するには", detail: "Terminal で「python3 ファイル名.py」と打てば動くよ！" },
      ],
    },
    linux: {
      steps: [
        { title: "Terminal を開く", detail: "アプリ一覧から「Terminal」を探して開こう" },
        { title: "パッケージリストを更新", detail: "「sudo apt update」と入力して Enter" },
        { title: "Python をインストール", detail: "「sudo apt install python3 python3-pip」と入力して Enter", warn: "パスワードを聞かれたら自分のパソコンのパスワードを入力してね" },
        { title: "動作確認", detail: "「python3 --version」と入力してバージョンが表示されればOK！" },
        { title: ".py ファイルを実行するには", detail: "「python3 ファイル名.py」と打てば動くよ！" },
      ],
    },
  },
};

export default config;
