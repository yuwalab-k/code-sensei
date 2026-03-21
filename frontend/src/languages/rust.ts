import type { LanguageConfig } from "./types";

const config: LanguageConfig = {
  def: {
    id: "rust",
    label: "Rust",
    icon: "memory",
    logo: "devicon-rust-original",
    description: "安全で高速なシステム言語",
    tag: "本格的なプログラミングに挑戦！",
    available: true,
    usedFor: ["システム開発", "WebAssembly", "組み込みシステム"],
  },
  outputMode: "text",
  initialCode: "// AI先生に話しかけてみよう！\n// 例：「掛け算クイズを作って」\nfn main() {\n    println!(\"Hello, World!\");\n}\n",
  features: [
    {
      id: "basics", icon: "code", title: "変数・型・基本構文", ready: true,
      description: "Rustの変数と型の仕組みを学ぼう",
      detail: "Rustでは let で変数を作るよ。変数は最初「変更できない」(immutable)で、変えたいときは let mut と書く。数字の型（i32, f64）や文字列型（String, &str）があって、コンパイラが型の間違いを事前に教えてくれるのがRustの強み！",
      trivia: "Stack Overflow の開発者アンケートで、Rustは10年連続で「最も愛されているプログラミング言語」に選ばれてるよ！難しいけど、使えるようになったプログラマーはみんな大好きになってしまう言語なんだ！",
      examples: [
        { text: "変数を使って自己紹介を表示する", context: "let で変数を定義して println! で表示。Rustの基本的な書き方が身につくよ！" },
        { text: "整数と小数の計算をする", context: "i32 と f64 の違いを体験。型の厳格さがRustの安全さを生む仕組みが学べるよ！" },
        { text: "かけ算九九の表を作る", context: "for ループと println! の組み合わせ。Rustの基本的な繰り返し処理が学べる！" },
        { text: "FizzBuzzプログラムを作る", context: "条件分岐と繰り返しの練習。Rustの match 式を使ったエレガントな書き方も試してみよう！" },
      ],
    },
    {
      id: "function", icon: "functions", title: "関数", ready: true,
      description: "fn で処理をまとめる関数を作ろう",
      detail: "Rustでは fn で関数を定義するよ。引数の型と戻り値の型を必ず書くのがRustのルール。型を書くことでバグが減り、コードが読みやすくなる。最後の式を return なしで書くのもRustらしい書き方！",
      trivia: "LinuxのカーネルにRustのコードが採用されたよ！2022年、世界中のサーバーやスマートフォンの土台になっているLinuxのシステムに、Rustで書かれたコードが組み込まれ始めた。これは大ニュースだったんだ！",
      examples: [
        { text: "2つの数字を足す関数を作る", context: "fn の基本。引数と戻り値の型を書く Rust ならではのスタイルが学べるよ！" },
        { text: "最大値を求める関数を作る", context: "比較演算と関数の組み合わせ。Rustの型安全性が実感できるよ！" },
        { text: "フィボナッチ数列を計算する関数を作る", context: "再帰関数の書き方。Rustでも再帰が書けることが学べる！" },
        { text: "温度を摂氏から華氏に変換する関数を作る", context: "実用的な計算関数。引数と戻り値の型宣言の練習にもなるよ！" },
      ],
    },
    {
      id: "control", icon: "psychology", title: "条件分岐・ループ", ready: true,
      description: "if・match・for・loopを使いこなそう",
      detail: "Rustには if 式、match 式、for ループ、while ループ、loop（無限ループ）がある。特に match はPythonのif/elifより強力で、パターンを使った条件分岐ができるよ。Rustらしい書き方が身につく重要な機能！",
      trivia: "Discordはサービスの一部をGoからRustに書き換えたよ！ゲーマー向けチャットツールのDiscordが「メモリ使用量が大幅に減った」「速度が上がった」という理由でRustに乗り換えた事例は、Rustの実力を証明するできごとだったんだ！",
      examples: [
        { text: "match式でじゃんけんの勝ち負けを判定する", context: "match パターンマッチングの入門。if/else より読みやすく書けるRustらしい書き方が学べるよ！" },
        { text: "forループで1から100の合計を求める", context: "Rustの範囲（1..=100）を使ったforループ。シンプルな繰り返し処理の基本！" },
        { text: "whileループで数字を当てるゲームを作る", context: "while ループの使い方。Rustでシンプルなゲームが作れるよ！" },
        { text: "素数を探すプログラムを作る", context: "条件分岐とループの組み合わせ。Rustの効率的な書き方で素数判定が学べるよ！" },
      ],
    },
    {
      id: "ownership", icon: "lock", title: "所有権・借用の基本", ready: true,
      description: "Rustならではの所有権システムを体験しよう",
      detail: "Rustの最大の特徴は「所有権（ownership）」システム！変数が「データの所有者」になって、所有者がいなくなると自動的にメモリが解放される。&（参照）を使えばデータを借りることもできる。難しく聞こえるけど、これがRustを安全・高速にしてる秘密だよ！",
      trivia: "GoogleはChromeブラウザの新しい部分をRustで書き始めてるよ！セキュリティのバグの70%が「メモリの使い方ミス」が原因なんだけど、Rustの所有権システムを使えばそのバグが起きないんだ。世界中で使われるChromeをRustが守ってる！",
      examples: [
        { text: "文字列の所有権の移動を体験する", context: "Rustの所有権の基本。「移動（move）」という考え方が学べるよ！" },
        { text: "参照（&）を使って文字列の長さを返す関数を作る", context: "借用の基本。データを渡さずに「借りる」仕組みが学べるよ！" },
        { text: "clone() を使ってデータをコピーする", context: "所有権を移動させずにコピーする方法。借用と所有権の違いが実感できる！" },
        { text: "文字列スライス（&str）を使う", context: "文字列の一部を参照する。Rustの String と &str の違いが学べるよ！" },
      ],
    },
    {
      id: "collections", icon: "list_alt", title: "ベクタ・文字列", ready: true,
      description: "データをまとめて管理するVecとStringを使おう",
      detail: "RustのVec（ベクタ）はPythonのリストに似た動的配列だよ。push() でデータを追加、iter() でループできる。Stringは変更できる文字列で、&str は変更できない文字列スライス。この違いを理解するとRustが使いこなせるようになる！",
      trivia: "AWS（アマゾンのクラウドサービス）の一部のサービスもRustで作られてるよ！世界中の何百万ものサーバーを動かすAWSが「安全で高速なRust」を選んだことで、Rustが「本物のシステム言語」として認められた証拠になったんだ！",
      examples: [
        { text: "ベクタに点数を追加して平均を計算する", context: "Vec::new() と push() の基本。Pythonのリストと比較しながら学べるよ！" },
        { text: "ベクタを for で繰り返して全要素を表示する", context: "iter() を使ったループ。Rustらしいイテレータの使い方が学べる！" },
        { text: "文字列を結合して文章を作る", context: "format! マクロで文字列を組み立てる。Rustの文字列操作の基本！" },
        { text: "ベクタから最大値・最小値を求める", context: "iter().max() と iter().min() を使う。Rustの関数型スタイルが体験できるよ！" },
      ],
    },
  ],
  frameworks: [
    {
      name: "Actix-web",
      description: "世界最速クラスのウェブサーバーフレームワーク",
      detail: "ウェブフレームワークのベンチマーク（速さ比較）で常にトップを争う超高速フレームワーク。「とにかく最速の API を作りたい・大量のリクエストを少ないサーバーで捌きたい」という要件に向く。ただし学習コストが高く、Rust の所有権・ライフタイムへの理解が必要。「Rust に慣れてきて、本格的なウェブサーバーを書いてみたい」段階で挑戦する。",
    },
    {
      name: "Axum",
      description: "Tokio 公式の使いやすいウェブフレームワーク",
      detail: "Tokio チームが開発した「書きやすさと速さを両立」したウェブフレームワーク。Actix-web より学習コストが低く、Rust のエコシステムに自然に馴染む設計。現在 Rust でウェブ開発をするなら Axum が最初の選択肢になってきている。「Rust でウェブ API を作りたい・Actix-web は難しそう」というときに選ぶ。",
    },
    {
      name: "Tokio",
      description: "Rust の非同期処理の標準ランタイム",
      detail: "「同時に大量の処理を並行して動かす」ための非同期ランタイム。Axum・Actix-web など Rust のほぼすべてのウェブフレームワークが Tokio の上で動いている。フレームワークというより「土台となるインフラ」に近い存在。Rust でネットワークやサーバーを扱うコードを書くなら必ず登場するライブラリ。",
    },
    {
      name: "Tauri",
      description: "Rust でデスクトップアプリを作れるフレームワーク",
      detail: "HTML・CSS・JavaScript でデザインし、Rust でバックエンドを書くデスクトップアプリフレームワーク。同じことができる Electron と比べてインストーラーサイズが 10 分の 1 以下でメモリ使用量も大幅に少ない。「ウェブの知識を活かしてデスクトップアプリを作りたいが、Electron は重すぎる」という場面で選ばれる。",
    },
  ],
  installGuide: {
    windows: {
      steps: [
        { title: "rustupをダウンロードする", detail: "「rustup-init.exe」を公式サイト（rustup.rs）からダウンロードしよう。Rustの公式インストーラーだよ。" },
        { title: "rustupを実行する", detail: "ダウンロードした rustup-init.exe を実行して、1番（デフォルト）を選んでEnterキーを押そう。", warn: "Visual C++ Build Toolsが必要な場合があるよ。画面の指示に従ってインストールしてね。" },
        { title: "インストール確認", detail: "コマンドプロンプトで「rustc --version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    mac: {
      steps: [
        { title: "rustupでインストールする", detail: "ターミナルで「curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh」と入力してEnterキーを押そう。" },
        { title: "パスを設定する", detail: "インストール後、「source $HOME/.cargo/env」を実行しよう。または新しいターミナルを開けば自動的に設定されるよ。" },
        { title: "インストール確認", detail: "ターミナルで「rustc --version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
    linux: {
      steps: [
        { title: "rustupでインストールする", detail: "ターミナルで「curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh」と入力してEnterキーを押そう。" },
        { title: "パスを設定する", detail: "インストール後、「source $HOME/.cargo/env」を実行しよう。または新しいターミナルを開けば自動的に設定されるよ。" },
        { title: "インストール確認", detail: "ターミナルで「rustc --version」と入力してEnterキーを押そう。バージョン番号が表示されたら成功！" },
      ],
    },
  },
};

export default config;
