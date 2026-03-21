import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "php",
    label: "PHP",
    icon: "public",
    logo: "devicon-php-plain colored",
    description: "ウェブサービスの裏側で動く言語",
    tag: "ウェブ開発に挑戦したい人に！",
    available: true,
    usedFor: ["ウェブ開発", "CMS開発", "サーバーサイド"],
  },
  outputMode: "text",
  initialCode: "<?php\n// AI先生に話しかけてみよう！\n// 例：「掛け算クイズを作って」\n",
  features: [
    {
      id: "calc", icon: "calculate", title: "計算・数学", ready: true,
      description: "数字を使った計算やゲームが作れる",
      detail: "PHPは数学計算がシンプルに書けるよ。rand() で乱数を作ったり、round() で四捨五入したり、abs() で絶対値を求めたりできる。echo で結果を表示するのが基本！",
      trivia: "WikipediaのシステムはPHPで作られてるよ！世界中で毎日何億人もの人が使うWikipedia。あの膨大な百科事典のページを表示している仕組みがPHPで動いてるんだ！",
      examples: [
        { text: "かけ算九九の表を作る", context: "for ループと echo で九九表が完成！PHPの基本的な繰り返し処理が学べるよ。" },
        { text: "素数を探すプログラムを作る", context: "素数はインターネットのSSL暗号の基礎。PHPで実装してみよう！" },
        { text: "ランダムな数当てゲームを作る", context: "rand() で乱数を作る方法が学べる。PHPでのゲーム作りの基本！" },
        { text: "フィボナッチ数列を表示する", context: "ひまわりの種や巻き貝の形にも現れる数列をPHPで再現しよう！" },
      ],
    },
    {
      id: "string", icon: "edit_note", title: "文字列処理", ready: true,
      description: "文字を加工したり変換したりできる",
      detail: "PHPは文字列処理がとても強力！strlen() で文字数を数えたり、str_replace() で文字を置き換えたり、strtoupper() で大文字にしたりできるよ。ウェブでユーザーの入力を処理するときによく使われる機能！",
      trivia: "Facebookは最初PHPで作られたよ！マーク・ザッカーバーグが大学の寮でPHPを使って作り始めたんだ。今は独自の言語（Hack）を使ってるけど、PHPがなければFacebookは生まれなかったかもしれない！",
      examples: [
        { text: "文字を逆さにして回文チェッカーを作る", context: "strrev() で一発逆転！PHPのシンプルな文字列関数が学べるよ。" },
        { text: "文章に含まれる単語を数える", context: "str_word_count() や explode() を使った単語カウント。SNSの文字数制限と同じ仕組みが学べる！" },
        { text: "名前を大文字・小文字に変換するツールを作る", context: "strtoupper() / strtolower() / ucfirst() の使い分けを練習しよう！" },
        { text: "シーザー暗号で文字を暗号化する", context: "ord() と chr() を使った暗号化。PHPでスパイ暗号が作れるよ！" },
      ],
    },
    {
      id: "array", icon: "list_alt", title: "配列処理", ready: true,
      description: "配列でデータをまとめて管理しよう",
      detail: "PHPの配列はとても便利！array_push() でデータを追加したり、array_filter() で絞り込んだり、sort() で並べ替えたりできるよ。数値キーだけでなく、文字列キーの「連想配列」もよく使われる！",
      trivia: "WordPressは世界中のウェブサイトの43%以上で使われてるよ！あなたが見るブログやニュースサイトの多くがWordPressで動いてて、その中心にあるのがPHPの配列処理なんだ。世界のウェブの土台を支えてる！",
      examples: [
        { text: "点数リストから合格者だけ抜き出す", context: "array_filter() でデータを絞り込む。ウェブサービスの検索機能と同じ仕組みが学べるよ！" },
        { text: "名前リストをランダムにシャッフルする", context: "shuffle() で一発シャッフル！PHPのシンプルさが実感できるよ。" },
        { text: "ショッピングカートの合計金額を計算する", context: "array_sum() で配列の合計を一発計算。ECサイトのカゴ機能と同じ仕組みだよ！" },
        { text: "テストの点数を集計して統計を出す", context: "max() / min() / array_sum() でデータ分析。ウェブアプリのダッシュボード機能の基礎！" },
      ],
    },
    {
      id: "function", icon: "functions", title: "関数", ready: true,
      description: "処理をまとめて再利用できる関数を作ろう",
      detail: "PHPの function でオリジナル関数が作れるよ！同じ処理を何度も書かなくてよくなる「再利用」の考え方が学べる。引数でデータを渡して、return で結果を返す仕組みがウェブ開発の基本！",
      trivia: "Slackは最初PHPで作られたよ！最初はゲーム会社だったSlackのチームが、社内コミュニケーションツールをPHPで作ったのが今のSlackの始まり。あの使いやすいチャットツールの原型がPHPだったんだ！",
      examples: [
        { text: "挨拶する関数を作る", context: "function の基本。同じコードを何度も書かなくていい「再利用」の考え方が学べるよ！" },
        { text: "BMIを計算する関数を作る", context: "引数と return を使った実用的な計算関数。健康管理アプリと同じ仕組みが学べる！" },
        { text: "西暦を和暦に変換する関数を作る", context: "条件分岐と関数の組み合わせ。日本のウェブサービスでよく使われる変換機能が作れるよ！" },
        { text: "パスワードの強度をチェックする関数を作る", context: "文字列処理と条件分岐の総合練習。セキュリティの基礎が学べる実用的な関数！" },
      ],
    },
    {
      id: "output", icon: "web", title: "出力・制御フロー", ready: true,
      description: "echoで出力し、条件や繰り返しを使いこなそう",
      detail: "PHPの echo で文字や数字を表示できるよ。if / elseif / else で条件分岐、for / foreach / while で繰り返し処理ができる。これらを組み合わせるとどんな複雑な処理も作れるようになる！",
      trivia: "PHPはPersonal Home Pageの略だったよ！1994年にラスマス・ラードフが自分のホームページ管理のために作った小さなプログラムが始まり。今では世界中のサーバーで動く大きな言語になったんだ！",
      examples: [
        { text: "FizzBuzzプログラムを作る", context: "3の倍数でFizz、5の倍数でBuzz。プログラミング面接の定番問題をPHPで書いてみよう！" },
        { text: "じゃんけんの勝ち負けを判定する", context: "if / elseif を使った条件分岐の練習。コンピューターとじゃんけんするゲームが作れるよ！" },
        { text: "星のピラミッドを描く", context: "for ループのネスト（入れ子）の練習。繰り返しの中に繰り返しを書く方法が身につくよ！" },
        { text: "数字の逆順カウントダウンを表示する", context: "for と echo の組み合わせ。条件付き繰り返しの基本が学べるよ！" },
      ],
    },
  ],
  frameworks: [
    {
      name: "Laravel",
      logo: "devicon-laravel-plain colored",
      description: "モダン PHP 開発の王道フレームワーク",
      detail: "現在 PHP フレームワークで圧倒的人気 No.1。「美しいコードが書けるフレームワーク」として設計され、認証・データベース・メール・キューなどウェブアプリに必要な機能が全部揃っている。日本の受託開発会社での採用率も高く、「PHP で仕事を取りたい」なら Laravel を覚えるのが一番の近道。ただし機能が多い分、覚えることも多い。",
    },
    {
      name: "Symfony",
      logo: "devicon-symfony-original colored",
      description: "大規模ウェブアプリ向けの堅牢なフレームワーク",
      detail: "Laravel 自身も Symfony の部品（コンポーネント）を使って作られているほど、PHP エコシステムの土台になっている老舗フレームワーク。大企業の基幹システムや複雑なビジネスロジックを持つウェブサービスに向く。ヨーロッパ企業での採用が特に多い。「長く運用する大きなシステムを作る・企業の基幹システムに関わりたい」なら選ばれる。",
    },
    {
      name: "CakePHP",
      logo: "devicon-cakephp-plain colored",
      description: "日本でも長年使われてきた老舗フレームワーク",
      detail: "Rails にインスパイアされ 2005 年に登場した老舗フレームワーク。日本の PHP 開発現場で長年採用されてきた実績があり、既存システムのメンテナンス案件などで今でも登場する。「設定より規約」で素早く開発できる点は Laravel と似ているが、現在の新規開発では Laravel が選ばれることが多い。既存の CakePHP プロジェクトを引き継ぐ機会がある人は把握しておきたい。",
    },
  ],
  installGuide: {
    windows: {
      steps: [
        { title: "XAMPPをダウンロードする", detail: "「XAMPP Windows」と検索して、公式サイトからダウンロードしよう。PHPとApacheが一緒に入ってる便利なパックだよ。" },
        { title: "XAMPPをインストールする", detail: "ダウンロードしたファイルを実行して、画面の指示に従ってインストールしよう。インストール場所は C:\\xampp がおすすめ。" },
        { title: "インストール確認", detail: "コマンドプロンプトで「php -v」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    mac: {
      steps: [
        { title: "Homebrewでインストールする", detail: "ターミナルで「brew install php」と入力してEnterキーを押そう。Homebrewが自動的にインストールしてくれるよ。" },
        { title: "インストール確認", detail: "ターミナルで「php -v」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    linux: {
      steps: [
        { title: "aptでインストールする", detail: "ターミナルで「sudo apt-get install php」と入力してEnterキーを押そう。パスワードを聞かれたら入力してね。" },
        { title: "インストール確認", detail: "ターミナルで「php -v」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
  },
};

export default config;
