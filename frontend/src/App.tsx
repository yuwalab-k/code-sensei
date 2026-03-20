import { useState, useEffect, useRef } from "react";
import Editor from "@monaco-editor/react";
import axios from "axios";
import ReactMarkdown from "react-markdown";
import "./App.css";

const API = "http://localhost:8081";

type Mode = "python";
type OS = "windows" | "mac" | "linux";
type ChatMode = "free" | "generated" | "reviewing";
type ExecMode = "browser" | "file";

interface ChatMessage {
  from: "ai" | "user";
  text: string;
}

interface FeatureExample {
  text: string;
  context: string;
}

interface Feature {
  id: string;
  icon: string;
  title: string;
  description: string;
  detail: string;
  trivia: string;
  examples: FeatureExample[];
  ready: boolean;
}

interface LanguageDef {
  id: Mode;
  label: string;
  icon: string;
  description: string;
  tag: string;
  available: boolean;
}

// ===== 言語定義 =====

const LANGUAGES: LanguageDef[] = [
  { id: "python", label: "Python", icon: "terminal", description: "計算やデータ処理が得意！AIでも使われる人気の言語", tag: "はじめての人におすすめ", available: true },
];

const COMING_SOON = [
  { label: "JavaScript", icon: "language",      description: "ウェブサイトを動かす言語" },
  { label: "Ruby",       icon: "diamond",       description: "読みやすくて書きやすい言語" },
  { label: "PHP",        icon: "public",        description: "ウェブサービスの裏側で動く言語" },
  { label: "Perl",       icon: "auto_fix_high", description: "テキスト処理が得意な言語" },
  { label: "Rust",       icon: "memory",        description: "安全で高速なシステム言語" },
  { label: "Go",         icon: "rocket_launch", description: "速くて軽い！サーバー向けの言語" },
];

// ===== 言語ごとの「できること」 =====

const FEATURES: Record<Mode, Feature[]> = {
  python: [
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
};

// ===== インストールガイド =====

const INSTALL_GUIDE: Record<OS, { steps: { title: string; detail: string; warn?: string }[] }> = {
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
};

// ===== 初期コード =====

const INITIAL_CODE = `# AI先生に話しかけてみよう！\n# 例：「掛け算の問題を出してくれるプログラムを作って」\n`;

const INITIAL_CHAT_BROWSER: ChatMessage[] = [
  { from: "ai", text: "こんにちは！わたしはAI先生だよ。「○○を作って」って話しかけると、コードを作ってあげるよ！" },
];

const INITIAL_CHAT_FILE: ChatMessage[] = [
  { from: "ai", text: "ファイル生成モードだよ！input()なども使えるプログラムを作れるよ。「○○を作って」って話しかけてみよう！" },
];

const isGenerateIntent = (msg: string) =>
  /作って|書いて|作ってほしい|書いてほしい|プログラムを|コードを書|教えて|見せて|やって|して|コードは|どう書|どうやって/.test(msg);

// ===== コンポーネント =====

export default function App() {
  const [execMode, setExecModeState] = useState<ExecMode>("browser");
  const [birthYear, setBirthYearState] = useState<number>(() => {
    const saved = localStorage.getItem("birthYear");
    return saved ? Number(saved) : 0;
  });
  const setBirthYear = (y: number) => {
    setBirthYearState(y);
    localStorage.setItem("birthYear", String(y));
  };
  const [showLangModal, setShowLangModal] = useState(false);
  const [showInstallModal, setShowInstallModal] = useState(false);
  const [showReviewModal, setShowReviewModal] = useState(false);
  const [showOverwriteModal, setShowOverwriteModal] = useState(false);
  const [pendingDescription, setPendingDescription] = useState("");
  const [showFeatureModal, setShowFeatureModal] = useState(false);
  const [modalFeature, setModalFeature] = useState<Feature | null>(null);
  const [installOS, setInstallOS] = useState<OS>("windows");

  // モードごとに独立した状態
  const [codes, setCodes] = useState<Record<ExecMode, string>>(() => ({
    browser: localStorage.getItem("codeBrowser") ?? INITIAL_CODE,
    file: localStorage.getItem("codeFile") ?? INITIAL_CODE,
  }));
  const [chats, setChats] = useState<Record<ExecMode, ChatMessage[]>>(() => ({
    browser: (() => {
      try {
        const s = localStorage.getItem("chatBrowser");
        if (s) return JSON.parse(s) as ChatMessage[];
      } catch {}
      return INITIAL_CHAT_BROWSER;
    })(),
    file: (() => {
      try {
        const s = localStorage.getItem("chatFile");
        if (s) return JSON.parse(s) as ChatMessage[];
      } catch {}
      return INITIAL_CHAT_FILE;
    })(),
  }));
  const [chatModes, setChatModes] = useState<Record<ExecMode, ChatMode>>({
    browser: "free",
    file: "free",
  });
  const [outputs, setOutputs] = useState<Record<ExecMode, string>>({
    browser: "",
    file: "",
  });

  // 穴抜け問題
  const [blankedCode, setBlankedCode] = useState<string | null>(null);
  const [cleared, setCleared] = useState(false);

  const [running, setRunning] = useState(false);
  const [question, setQuestion] = useState("");
  const [askingHint, setAskingHint] = useState(false);
  const [thinkingSec, setThinkingSec] = useState(0);
  const thinkingTimer = useRef<ReturnType<typeof setInterval> | null>(null);
  const [panelWidth, setPanelWidth] = useState(() => Math.round(window.innerWidth * 0.25));
  const chatEndRef = useRef<HTMLDivElement>(null);
  const abortRef = useRef<AbortController | null>(null);
  const dragging = useRef(false);
  const dragStartX = useRef(0);
  const dragStartWidth = useRef(0);

  // リサイズハンドル
  useEffect(() => {
    const onMouseMove = (e: MouseEvent) => {
      if (!dragging.current) return;
      const delta = e.clientX - dragStartX.current;
      setPanelWidth(Math.max(200, Math.min(520, dragStartWidth.current + delta)));
    };
    const onMouseUp = () => { dragging.current = false; };
    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    return () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };
  }, []);

  // コード保存
  useEffect(() => {
    localStorage.setItem("codeBrowser", codes.browser);
    localStorage.setItem("codeFile", codes.file);
  }, [codes]);

  // チャット保存 & スクロール
  useEffect(() => {
    localStorage.setItem("chatBrowser", JSON.stringify(chats.browser));
    localStorage.setItem("chatFile", JSON.stringify(chats.file));
    chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chats]);

  // モード切り替え時に穴抜けをリセット
  const setExecMode = (next: ExecMode) => {
    setExecModeState(next);
    setBlankedCode(null);
    setCleared(false);
  };

  // ===== ヘルパー =====

  const chat = chats[execMode];
  const chatMode = chatModes[execMode];
  const output = outputs[execMode];

  // エディタに表示するコード（復習時は穴抜きコード）
  const editorCode = chatMode === "reviewing" && blankedCode !== null
    ? blankedCode
    : codes[execMode];

  const handleEditorChange = (v: string | undefined) => {
    if (chatMode === "reviewing" && blankedCode !== null) {
      setBlankedCode(v ?? "");
    } else {
      setCodes(prev => ({ ...prev, [execMode]: v ?? "" }));
    }
  };

  const addChat = (msg: ChatMessage) =>
    setChats(prev => ({ ...prev, [execMode]: [...prev[execMode], msg] }));


  const setChatMode = (m: ChatMode) =>
    setChatModes(prev => ({ ...prev, [execMode]: m }));

  const setOutput = (o: string) =>
    setOutputs(prev => ({ ...prev, [execMode]: o }));

  const isErrorOutput = (out: string) =>
    out.includes("Traceback") || out.includes("Error") || out.includes("error");

  // ===== 実行 =====

  const runCode = async () => {
    const codeToRun = chatMode === "reviewing" && blankedCode !== null
      ? blankedCode
      : codes[execMode];
    setRunning(true);
    setOutput("実行中...");
    try {
      const res = await axios.post<{ output: string }>(`${API}/execute`, { code: codeToRun, language: "python" });
      const out = res.data.output || "（出力なし）";
      setOutput(out);

      // 復習モード（穴抜け）の判定
      if (chatMode === "reviewing" && blankedCode !== null) {
        if (isErrorOutput(out)) {
          addChat({ from: "ai", text: `惜しい！エラーが出たよ。「___」の部分をもう一度考えてみよう！わからなかったらヒントを聞いてね。\n\nエラー:\n${out}` });
        } else {
          setCleared(true);
          addChat({ from: "ai", text: "やったー！クリア！穴埋めが全部できたね！すごいよ！" });
        }
      }
    } catch {
      setOutput("エラー: バックエンドに接続できませんでした");
    } finally {
      setRunning(false);
    }
  };

  // ===== チャット操作 =====

  const deleteMessage = (index: number) =>
    setChats(prev => ({ ...prev, [execMode]: prev[execMode].filter((_, i) => i !== index) }));

  const resetChat = () => {
    const initial = execMode === "browser" ? INITIAL_CHAT_BROWSER : INITIAL_CHAT_FILE;
    setChats(prev => ({ ...prev, [execMode]: [{ from: "ai", text: "チャットをリセットしたよ！また何でも聞いてね！" }] }));
    localStorage.setItem(execMode === "browser" ? "chatBrowser" : "chatFile", JSON.stringify(initial));
    setChatMode("free");
    setBlankedCode(null);
    setCleared(false);
  };

  // ===== ダウンロード =====

  const downloadPyFile = (src: string, filename = "program.py") => {
    const blob = new Blob([src], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  // ===== Thinking タイマー =====

  const startThinking = () => {
    setThinkingSec(0);
    setAskingHint(true);
    thinkingTimer.current = setInterval(() => setThinkingSec(s => s + 1), 1000);
  };

  const stopThinking = () => {
    setAskingHint(false);
    if (thinkingTimer.current) {
      clearInterval(thinkingTimer.current);
      thinkingTimer.current = null;
    }
  };

  // ===== 中止 =====

  const cancelRequest = () => {
    abortRef.current?.abort();
    abortRef.current = null;
    stopThinking();
    addChat({ from: "ai", text: "中止したよ。また話しかけてね！" });
  };

  const newAbort = () => {
    abortRef.current?.abort();
    const ctrl = new AbortController();
    abortRef.current = ctrl;
    return ctrl.signal;
  };

  // ===== AI 生成 =====

  const generateCode = async (description: string) => {
    const signal = newAbort();
    startThinking();
    try {
      const res = await axios.post<{ code: string; explanation: string }>(`${API}/ai/generate`, {
        description,
        language: "python",
        exec_mode: execMode,
        birth_year: birthYear,
      }, { signal });
      const { code: generated, explanation } = res.data;
      if (!generated) {
        addChat({ from: "ai", text: "ごめんね、うまく作れなかったよ。もう一度試してみて！" });
        return;
      }
      setCodes(prev => ({ ...prev, [execMode]: generated }));
      setBlankedCode(null);
      setCleared(false);

      if (execMode === "file") {
        addChat({
          from: "ai",
          text: `できたよ！エディタにコードを入れたよ。\n\n${explanation}\n\n「ファイルをダウンロード」を押して、パソコンに保存してね！`,
        });
        setChatMode("generated");
        return;
      }

      // ブラウザ実行モード：自動実行してエラーチェック
      addChat({ from: "ai", text: `できたよ！エディタにコードを入れたよ。\n\n${explanation}` });
      setOutput("実行中...");
      try {
        const runRes = await axios.post<{ output: string }>(`${API}/execute`, { code: generated, language: "python" });
        const out = runRes.data.output || "（出力なし）";
        setOutput(out);
        if (isErrorOutput(out)) {
          addChat({ from: "ai", text: `実行したらエラーが出たよ。もう一度作り直すね！\n\nエラー:\n${out}` });
          await regenerateWithError(description, generated, out, signal);
          return;
        }
      } catch {
        setOutput("");
      }
      setChatMode("generated");
    } catch (e) {
      if (axios.isCancel(e)) return;
      addChat({ from: "ai", text: "ごめんね、うまく作れなかったよ。もう一度試してみて！" });
    } finally {
      stopThinking();
    }
  };

  const regenerateWithError = async (description: string, prevCode: string, errorMsg: string, signal: AbortSignal) => {
    try {
      const res = await axios.post<{ code: string; explanation: string }>(`${API}/ai/generate`, {
        description: `${description}\n\n前回のコードでエラーが出たので修正してください。\nエラー内容:\n${errorMsg}\n前回のコード:\n${prevCode}`,
        language: "python",
        exec_mode: execMode,
        birth_year: birthYear,
      }, { signal });
      const { code: fixed, explanation } = res.data;
      if (!fixed) return;
      setCodes(prev => ({ ...prev, [execMode]: fixed }));
      setOutput("実行中...");
      const runRes = await axios.post<{ output: string }>(`${API}/execute`, { code: fixed, language: "python" });
      const out = runRes.data.output || "（出力なし）";
      setOutput(out);
      if (isErrorOutput(out)) {
        addChat({ from: "ai", text: "もう一度試したけどまだエラーが出ちゃった。内容を変えて再挑戦してみてね！" });
      } else {
        addChat({ from: "ai", text: `修正できたよ！\n\n${explanation}` });
        setChatMode("generated");
      }
    } catch (e) {
      if (!axios.isCancel(e)) throw e;
    }
  };

  // ===== 復習（穴抜け問題） =====

  const startReview = async () => {
    const signal = newAbort();
    setShowReviewModal(false);
    setChatMode("reviewing");
    setBlankedCode(null);
    setCleared(false);
    startThinking();

    const currentCode = codes[execMode];
    try {
      const [reviewRes, blankRes] = await Promise.all([
        axios.post<{ questions: string }>(`${API}/ai/review`, { code: currentCode, language: "python", birth_year: birthYear }, { signal }),
        axios.post<{ blanked_code: string }>(`${API}/ai/blank`, { code: currentCode, language: "python", birth_year: birthYear }, { signal }),
      ]);

      addChat({
        from: "ai",
        text: `復習タイム！\n\n${reviewRes.data.questions}\n\nエディタのコードが穴抜き問題になってるよ。「___」を埋めて実行してみよう！正解したらクリア！`,
      });
      setBlankedCode(blankRes.data.blanked_code || currentCode);
    } catch (e) {
      if (axios.isCancel(e)) return;
      addChat({ from: "ai", text: "ごめんね、問題を作れなかったよ。もう一度試してみて！" });
      setChatMode("generated");
    } finally {
      stopThinking();
    }
  };

  // ===== ヒント =====

  const askHint = async () => {
    if (!question.trim() || askingHint) return;
    const userMsg = question;
    addChat({ from: "user", text: userMsg });
    setQuestion("");

    if (chatMode !== "reviewing" && isGenerateIntent(userMsg)) {
      const hasCode = codes[execMode].trim() !== "" && codes[execMode].trim() !== INITIAL_CODE.trim();
      if (hasCode) {
        setPendingDescription(userMsg);
        setShowOverwriteModal(true);
        return;
      }
      await generateCode(userMsg);
      return;
    }

    const signal = newAbort();
    startThinking();
    try {
      const res = await axios.post<{ hint: string }>(`${API}/ai/hint`, {
        code: codes[execMode],
        question: userMsg,
        mission: "",
        language: "python",
        review_mode: chatMode === "reviewing",
        birth_year: birthYear,
      }, { signal });
      addChat({ from: "ai", text: res.data.hint });
    } catch (e) {
      if (axios.isCancel(e)) return;
      addChat({ from: "ai", text: "ごめんね、今ちょっと答えられないや。もう一度試してみて！" });
    } finally {
      stopThinking();
    }
  };

  const selectFeature = (f: Feature) => {
    setModalFeature(f);
    setShowFeatureModal(true);
  };

  const currentLang = LANGUAGES[0];
  const features = FEATURES["python"];

  return (
    <div className="app">

      {/* ===== 言語選択モーダル ===== */}
      {showLangModal && (
        <div className="modal-overlay" onClick={() => setShowLangModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <span className="modal-title">言語を選ぼう</span>
              <button className="modal-close" onClick={() => setShowLangModal(false)}>✕</button>
            </div>
            <div className="lang-grid">
              {LANGUAGES.map((lang) => (
                <button key={lang.id} className="lang-card selected" onClick={() => setShowLangModal(false)}>
                  <div className="lang-card-icon"><span className="material-symbols-outlined">{lang.icon}</span></div>
                  <div className="lang-card-label">{lang.label}</div>
                  <div className="lang-card-desc">{lang.description}</div>
                  <div className="lang-card-tag">{lang.tag}</div>
                  <div className="lang-card-check">✓ 使用中</div>
                </button>
              ))}
              {COMING_SOON.map((lang) => (
                <div key={lang.label} className="lang-card disabled">
                  <div className="lang-card-icon"><span className="material-symbols-outlined">{lang.icon}</span></div>
                  <div className="lang-card-label">{lang.label}</div>
                  <div className="lang-card-desc">{lang.description}</div>
                  <div className="lang-card-soon">準備中</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* ===== インストールガイドモーダル ===== */}
      {showInstallModal && (
        <div className="modal-overlay" onClick={() => setShowInstallModal(false)}>
          <div className="modal modal-install" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <span className="modal-title">Python のインストール方法</span>
              <button className="modal-close" onClick={() => setShowInstallModal(false)}>✕</button>
            </div>
            <div className="os-tabs">
              {(["windows", "mac", "linux"] as OS[]).map((os) => (
                <button key={os} className={`os-tab ${installOS === os ? "active" : ""}`} onClick={() => setInstallOS(os)}>
                  {os === "windows" ? "Windows" : os === "mac" ? "Mac" : "Linux"}
                </button>
              ))}
            </div>
            <div className="install-steps">
              {INSTALL_GUIDE[installOS].steps.map((step, i) => (
                <div key={i} className="install-step">
                  <div className="install-step-num">{i + 1}</div>
                  <div className="install-step-body">
                    <div className="install-step-title">{step.title}</div>
                    <div className="install-step-detail">{step.detail}</div>
                    {step.warn && <div className="install-step-warn">{step.warn}</div>}
                  </div>
                </div>
              ))}
            </div>
            <div className="install-footer">
              インストールできたら、このアプリで学んだコードをパソコンでも動かしてみよう！
            </div>
          </div>
        </div>
      )}

      {/* ===== 特徴説明モーダル ===== */}
      {showFeatureModal && modalFeature && (
        <div className="modal-overlay" onClick={() => setShowFeatureModal(false)}>
          <div className="modal modal-feature" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <span className="modal-title"><span className="material-symbols-outlined">{modalFeature.icon}</span> {modalFeature.title}</span>
              <button className="modal-close" onClick={() => setShowFeatureModal(false)}>✕</button>
            </div>
            <p className="feature-modal-detail">{modalFeature.detail}</p>
            {modalFeature.trivia && (
              <div className="feature-modal-trivia">
                <span className="material-symbols-outlined">lightbulb</span>
                <span>{modalFeature.trivia}</span>
              </div>
            )}
            <div className="feature-modal-examples-title">こんなものが作れるよ</div>
            <ul className="feature-modal-examples">
              {modalFeature.examples.map((ex, i) => (
                <li key={i}>
                  <div className="feature-example-text">{ex.text}</div>
                  <div className="feature-example-context">{ex.context}</div>
                </li>
              ))}
            </ul>
            {!modalFeature.ready && (
              <div className="feature-modal-soon">このカテゴリは現在準備中です</div>
            )}
          </div>
        </div>
      )}

      {/* ===== 復習確認モーダル ===== */}
      {showReviewModal && (
        <div className="modal-overlay" onClick={() => setShowReviewModal(false)}>
          <div className="modal modal-review" onClick={(e) => e.stopPropagation()}>
            <div className="modal-review-icon">📝</div>
            <div className="modal-review-title">復習に挑戦しよう！</div>
            <div className="modal-review-body">
              コードが完成したね！<br />
              AI先生が復習問題と穴抜け問題を作るよ。<br />
              エディタの「___」を埋めて実行しよう！
            </div>
            <div className="modal-review-actions">
              <button className="btn-review-ok" onClick={startReview}>はい！挑戦する！</button>
              <button className="btn-review-cancel" onClick={() => setShowReviewModal(false)}>もう少し見る</button>
            </div>
          </div>
        </div>
      )}

      {/* ===== 上書き確認モーダル ===== */}
      {showOverwriteModal && (
        <div className="modal-overlay" onClick={() => setShowOverwriteModal(false)}>
          <div className="modal modal-review" onClick={(e) => e.stopPropagation()}>
            <div className="modal-review-icon">
              <span className="material-symbols-outlined" style={{ fontSize: 40, color: "#f6ad55" }}>edit_note</span>
            </div>
            <div className="modal-review-title">エディタにコードがあるよ</div>
            <div className="modal-review-body">
              今のコードはどうする？
            </div>
            <div className="modal-review-actions">
              <button className="btn-review-ok" onClick={() => {
                setShowOverwriteModal(false);
                generateCode(pendingDescription);
              }}>
                上書きして新しく作る
              </button>
              <button className="btn-review-cancel" onClick={() => setShowOverwriteModal(false)}>
                今のコードを残す
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ===== ヘッダー ===== */}
      <header className="header">
        <div className="header-logo">
          <span className="material-symbols-outlined logo-icon-sym">smart_toy</span>
          <span className="logo-text">Code Sensei</span>
        </div>
        <span className="header-sub">プログラミング学習</span>
        <div className="header-right">
          <div className="birth-year-selector">
            <span className="material-symbols-outlined birth-year-icon">person</span>
            <input
              className="birth-year-input"
              type="number"
              min={1990}
              max={new Date().getFullYear()}
              placeholder="生まれ年（例: 2015）"
              value={birthYear === 0 ? "" : birthYear}
              onChange={(e) => setBirthYear(e.target.value === "" ? 0 : Number(e.target.value))}
            />
            {birthYear > 0 && (
              <span className="birth-year-age">{new Date().getFullYear() - birthYear}歳</span>
            )}
          </div>
          <button className="lang-switcher-btn" onClick={() => setShowLangModal(true)}>
            <span className="material-symbols-outlined">{currentLang.icon}</span>
            <span>{currentLang.label}</span>
            <span className="material-symbols-outlined">expand_more</span>
          </button>
        </div>
      </header>

      {/* ===== 言語特徴パネル ===== */}
      <div className="features-bar">
        <div className="features-bar-left">
          <span className="features-bar-title">
            <span className="material-symbols-outlined">{currentLang.icon}</span> {currentLang.label} でできること
          </span>
        </div>
        <div className="features-list">
          {features.map((f) => (
            <button
              key={f.id}
              className={`feature-chip ${!f.ready ? "disabled" : ""}`}
              onClick={() => selectFeature(f)}
              title={!f.ready ? "準備中" : f.description}
            >
              <span className="material-symbols-outlined">{f.icon}</span>
              <span>{f.title}</span>
              {!f.ready && <span className="feature-soon">準備中</span>}
            </button>
          ))}
        </div>
      </div>

      {/* ===== メインエリア ===== */}
      <div className="main">

        {/* AI チャット */}
        <aside className="ai-panel" style={{ width: panelWidth, minWidth: panelWidth }}>
          <div className="ai-panel-header">
            <span className="material-symbols-outlined">smart_toy</span>
            <span>AI 先生</span>
            {chatMode === "reviewing" && <span className="review-badge">復習モード</span>}
            <button className="btn-chat-reset" onClick={resetChat} title="チャットをリセット">
              <span className="material-symbols-outlined">delete</span>
            </button>
          </div>

          <div className="chat-messages">
            {chat.map((msg, i) => (
              <div key={i} className={`chat-bubble ${msg.from}`}>
                {msg.from === "ai" && <span className="material-symbols-outlined chat-avatar-icon">smart_toy</span>}
                <div className="chat-text"><ReactMarkdown>{msg.text}</ReactMarkdown></div>
                <button className="chat-delete" onClick={() => deleteMessage(i)} title="削除">
                  <span className="material-symbols-outlined">close</span>
                </button>
              </div>
            ))}
            {askingHint && (
              <div className="chat-bubble ai">
                <span className="material-symbols-outlined chat-avatar-icon">smart_toy</span>
                <div className="chat-text chat-thinking">考えてるよ... <span className="thinking-sec">{thinkingSec}秒</span></div>
                <button className="btn-cancel" onClick={cancelRequest} title="中止する">
                  <span className="material-symbols-outlined">stop_circle</span>
                </button>
              </div>
            )}
            <div ref={chatEndRef} />
          </div>

          {chatMode === "generated" && (
            <button className="btn-review-start" onClick={() => setShowReviewModal(true)}>
              <span className="material-symbols-outlined">school</span>
              復習問題に進む
            </button>
          )}

          <div className="chat-input-area">
            <input
              className="chat-input"
              placeholder={chatMode === "reviewing" ? "ヒントを聞いてね..." : "「○○を作って」と話しかけよう！"}
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && !e.nativeEvent.isComposing && askHint()}
              disabled={askingHint}
            />
            <button className="btn-hint" onClick={askHint} disabled={askingHint || !question.trim()}>
              <span className="material-symbols-outlined">send</span>
            </button>
          </div>
        </aside>

        {/* リサイズハンドル */}
        <div
          className="panel-resize-handle"
          onMouseDown={(e) => {
            dragging.current = true;
            dragStartX.current = e.clientX;
            dragStartWidth.current = panelWidth;
            e.preventDefault();
          }}
        />

        {/* エディタ＋出力 */}
        <div className="editor-area">
          <div className="editor-header">
            <span className="material-symbols-outlined">code</span>
            <span className="editor-title">コードエディタ</span>
            <span className="lang-badge">
              <span className="material-symbols-outlined">{currentLang.icon}</span> {currentLang.label}
            </span>
            {chatMode === "reviewing" && blankedCode !== null && (
              <span className="blank-badge">
                <span className="material-symbols-outlined">edit</span>
                穴抜け問題：「___」を埋めよう
              </span>
            )}
            <div className="exec-mode-toggle">
              <button
                className={`exec-mode-btn ${execMode === "browser" ? "active" : ""}`}
                onClick={() => setExecMode("browser")}
                title="アプリ内でそのまま実行"
              >
                <span className="material-symbols-outlined">web</span>
                ブラウザ実行
              </button>
              <button
                className={`exec-mode-btn ${execMode === "file" ? "active" : ""}`}
                onClick={() => setExecMode("file")}
                title="ファイルをダウンロードしてPCで実行"
              >
                <span className="material-symbols-outlined">download</span>
                ファイル生成
              </button>
            </div>
          </div>

          <div className="editor-wrapper">
            <Editor
              height="300px"
              language="python"
              value={editorCode}
              onChange={handleEditorChange}
              theme="vs-dark"
              options={{ fontSize: 14, minimap: { enabled: false }, scrollBeyondLastLine: false, lineNumbers: "on", fontFamily: "'JetBrains Mono', 'Fira Code', monospace" }}
            />
          </div>

          {/* クリアバナー */}
          {cleared && (
            <div className="cleared-banner">
              <span className="material-symbols-outlined">emoji_events</span>
              <span>クリア！すごい！穴埋めが全部できたよ！</span>
              <button className="cleared-close" onClick={() => setCleared(false)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>
          )}

          <div className="run-area">
            {execMode === "file" ? (
              <>
                <button className="btn-download" onClick={() => {
                  downloadPyFile(codes[execMode]);
                  addChat({
                    from: "ai",
                    text: "ダウンロードしたよ！実行するには：\n① ターミナル（Windowsはコマンドプロンプト）を開く\n② このコマンドを入力してEnterキーを押す\n   python program.py\n   （Macは python3 program.py）",
                  });
                }}>
                  <span className="material-symbols-outlined">download</span>
                  ファイルをダウンロード
                </button>
                <button className="btn-install-inline" onClick={() => setShowInstallModal(true)}>
                  <span className="material-symbols-outlined">build</span>
                  インストール方法
                </button>
              </>
            ) : (
              <button className="btn-run" onClick={runCode} disabled={running}>
                <span className="material-symbols-outlined">{running ? "hourglass_empty" : "play_arrow"}</span>
                {running ? "実行中..." : "実行する"}
              </button>
            )}
          </div>

          {execMode !== "file" && (
            <div className="output-panel">
              <div className="output-header">
                <span className="material-symbols-outlined">terminal</span>
                <span>出力</span>
              </div>
              <pre className="output-content">{output || "「実行する」を押すと、ここに結果が出るよ！"}</pre>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
