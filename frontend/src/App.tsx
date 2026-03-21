import { useState, useEffect, useRef } from "react";
import Editor from "@monaco-editor/react";
import axios from "axios";
import ReactMarkdown from "react-markdown";
import "./App.css";
import { LANGUAGE_CONFIGS, LANGUAGES, COMING_SOON } from "./languages";
import type { LanguageId, Feature, LanguageDef } from "./languages";

const API = "http://localhost:8081";

// ロゴが設定されていれば devicon、なければ Material Symbol を表示
function LangIcon({ def, className }: { def: Pick<LanguageDef, "icon" | "logo">; className?: string }) {
  if (def.logo) return <i className={`${def.logo} ${className ?? ""}`} style={{ fontSize: "inherit" }} />;
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{def.icon}</span>;
}

type OS = "windows" | "mac" | "linux";
type ChatMode = "free" | "generated" | "reviewing";
type ExecMode = "browser" | "file";

interface ChatMessage {
  from: "ai" | "user";
  text: string;
}

// ===== 初期チャットメッセージ =====

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
  const [language, setLanguage] = useState<LanguageId>(() => {
    const saved = localStorage.getItem("language");
    return (saved === "python" || saved === "javascript") ? saved : "python";
  });
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

  // 言語×モードごとに独立した状態
  const [codes, setCodes] = useState<Record<LanguageId, Record<ExecMode, string>>>(() => ({
    python: {
      browser: localStorage.getItem("codePythonBrowser") ?? LANGUAGE_CONFIGS.python.initialCode,
      file: localStorage.getItem("codePythonFile") ?? LANGUAGE_CONFIGS.python.initialCode,
    },
    javascript: {
      browser: localStorage.getItem("codeJsBrowser") ?? LANGUAGE_CONFIGS.javascript.initialCode,
      file: localStorage.getItem("codeJsFile") ?? LANGUAGE_CONFIGS.javascript.initialCode,
    },
  }));

  const [chats, setChats] = useState<Record<LanguageId, Record<ExecMode, ChatMessage[]>>>(() => ({
    python: {
      browser: (() => { try { const s = localStorage.getItem("chatPythonBrowser"); if (s) return JSON.parse(s); } catch {} return INITIAL_CHAT_BROWSER; })(),
      file: (() => { try { const s = localStorage.getItem("chatPythonFile"); if (s) return JSON.parse(s); } catch {} return INITIAL_CHAT_FILE; })(),
    },
    javascript: {
      browser: (() => { try { const s = localStorage.getItem("chatJsBrowser"); if (s) return JSON.parse(s); } catch {} return INITIAL_CHAT_BROWSER; })(),
      file: (() => { try { const s = localStorage.getItem("chatJsFile"); if (s) return JSON.parse(s); } catch {} return INITIAL_CHAT_FILE; })(),
    },
  }));

  const [chatModes, setChatModes] = useState<Record<LanguageId, Record<ExecMode, ChatMode>>>({
    python: { browser: "free", file: "free" },
    javascript: { browser: "free", file: "free" },
  });

  const [outputs, setOutputs] = useState<Record<LanguageId, Record<ExecMode, string>>>({
    python: { browser: "", file: "" },
    javascript: { browser: "", file: "" },
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
    localStorage.setItem("codePythonBrowser", codes.python.browser);
    localStorage.setItem("codePythonFile", codes.python.file);
    localStorage.setItem("codeJsBrowser", codes.javascript.browser);
    localStorage.setItem("codeJsFile", codes.javascript.file);
  }, [codes]);

  // チャット保存 & スクロール
  useEffect(() => {
    localStorage.setItem("chatPythonBrowser", JSON.stringify(chats.python.browser));
    localStorage.setItem("chatPythonFile", JSON.stringify(chats.python.file));
    localStorage.setItem("chatJsBrowser", JSON.stringify(chats.javascript.browser));
    localStorage.setItem("chatJsFile", JSON.stringify(chats.javascript.file));
    chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chats]);

  // 言語保存
  useEffect(() => {
    localStorage.setItem("language", language);
  }, [language]);

  // モード切り替え時に穴抜けをリセット
  const setExecMode = (next: ExecMode) => {
    setExecModeState(next);
    setBlankedCode(null);
    setCleared(false);
  };

  // 言語切り替え
  const switchLanguage = (lang: LanguageId) => {
    setLanguage(lang);
    setShowLangModal(false);
    setBlankedCode(null);
    setCleared(false);
  };

  // ===== ヘルパー =====

  const currentConfig = LANGUAGE_CONFIGS[language];
  const chat = chats[language][execMode];
  const chatMode = chatModes[language][execMode];
  const output = outputs[language][execMode];
  const editorCode = chatMode === "reviewing" && blankedCode !== null
    ? blankedCode
    : codes[language][execMode];

  const handleEditorChange = (v: string | undefined) => {
    if (chatMode === "reviewing" && blankedCode !== null) {
      setBlankedCode(v ?? "");
    } else {
      setCodes(prev => ({
        ...prev,
        [language]: { ...prev[language], [execMode]: v ?? "" },
      }));
    }
  };

  const addChat = (msg: ChatMessage) =>
    setChats(prev => ({
      ...prev,
      [language]: {
        ...prev[language],
        [execMode]: [...prev[language][execMode], msg],
      },
    }));

  const setChatMode = (m: ChatMode) =>
    setChatModes(prev => ({
      ...prev,
      [language]: { ...prev[language], [execMode]: m },
    }));

  const setOutput = (o: string) =>
    setOutputs(prev => ({
      ...prev,
      [language]: { ...prev[language], [execMode]: o },
    }));

  const isErrorOutput = (out: string) =>
    out.includes("Traceback") || out.includes("Error") || out.includes("error");

  // ===== 実行 =====

  const runCode = async () => {
    const codeToRun = chatMode === "reviewing" && blankedCode !== null
      ? blankedCode
      : codes[language][execMode];

    setRunning(true);
    setOutput("実行中...");
    try {
      const res = await axios.post<{ output: string }>(`${API}/execute`, { code: codeToRun, language });
      const out = res.data.output || "（出力なし）";
      setOutput(out);
      if (chatMode === "reviewing" && blankedCode !== null) {
        if (isErrorOutput(out)) {
          addChat({ from: "ai", text: `惜しい！エラーが出たよ。「___」の部分をもう一度考えてみよう！\n\nエラー:\n${out}` });
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
    setChats(prev => ({
      ...prev,
      [language]: {
        ...prev[language],
        [execMode]: prev[language][execMode].filter((_, i) => i !== index),
      },
    }));

  const resetChat = () => {
    const initial = execMode === "browser" ? INITIAL_CHAT_BROWSER : INITIAL_CHAT_FILE;
    setChats(prev => ({
      ...prev,
      [language]: {
        ...prev[language],
        [execMode]: [{ from: "ai", text: "チャットをリセットしたよ！また何でも聞いてね！" }],
      },
    }));
    localStorage.setItem(`chat${language === "python" ? "Python" : "Js"}${execMode === "browser" ? "Browser" : "File"}`, JSON.stringify(initial));
    setChatMode("free");
    setBlankedCode(null);
    setCleared(false);
  };

  // ===== ダウンロード =====

  const downloadFile = (src: string) => {
    const ext = language === "javascript" ? "js" : "py";
    const blob = new Blob([src], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `program.${ext}`;
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
        language,
        exec_mode: execMode,
        birth_year: birthYear,
      }, { signal });
      const { code: generated, explanation } = res.data;
      if (!generated) {
        addChat({ from: "ai", text: "ごめんね、うまく作れなかったよ。もう一度試してみて！" });
        return;
      }
      setCodes(prev => ({ ...prev, [language]: { ...prev[language], [execMode]: generated } }));
      setBlankedCode(null);
      setCleared(false);

      if (execMode === "file") {
        addChat({ from: "ai", text: `できたよ！エディタにコードを入れたよ。\n\n${explanation}\n\n「ファイルをダウンロード」を押して、パソコンに保存してね！` });
        setChatMode("generated");
        return;
      }

      addChat({ from: "ai", text: `できたよ！エディタにコードを入れたよ。\n\n${explanation}` });
      setOutput("実行中...");
      try {
        const runRes = await axios.post<{ output: string }>(`${API}/execute`, { code: generated, language });
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
        language,
        exec_mode: execMode,
        birth_year: birthYear,
      }, { signal });
      const { code: fixed, explanation } = res.data;
      if (!fixed) return;
      setCodes(prev => ({ ...prev, [language]: { ...prev[language], [execMode]: fixed } }));
      setOutput("実行中...");
      const runRes = await axios.post<{ output: string }>(`${API}/execute`, { code: fixed, language });
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

    const currentCode = codes[language][execMode];
    try {
      const [reviewRes, blankRes] = await Promise.all([
        axios.post<{ questions: string }>(`${API}/ai/review`, { code: currentCode, language, birth_year: birthYear }, { signal }),
        axios.post<{ blanked_code: string }>(`${API}/ai/blank`, { code: currentCode, language, birth_year: birthYear }, { signal }),
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
      const hasCode = codes[language][execMode].trim() !== "" && codes[language][execMode].trim() !== currentConfig.initialCode.trim();
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
        code: codes[language][execMode],
        question: userMsg,
        mission: "",
        language,
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

  const installGuide = currentConfig.installGuide;

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
                <button
                  key={lang.id}
                  className={`lang-card ${lang.id === language ? "selected" : ""}`}
                  onClick={() => switchLanguage(lang.id)}
                >
                  <div className="lang-card-icon"><LangIcon def={lang} /></div>
                  <div className="lang-card-label">{lang.label}</div>
                  <div className="lang-card-desc">{lang.description}</div>
                  <div className="lang-card-tag">{lang.tag}</div>
                  {lang.id === language && <div className="lang-card-check">✓ 使用中</div>}
                </button>
              ))}
              {COMING_SOON.map((lang) => (
                <div key={lang.label} className="lang-card disabled">
                  <div className="lang-card-icon"><LangIcon def={lang} /></div>
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
      {showInstallModal && installGuide && (
        <div className="modal-overlay" onClick={() => setShowInstallModal(false)}>
          <div className="modal modal-install" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <span className="modal-title">{currentConfig.def.label} のインストール方法</span>
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
              {installGuide[installOS].steps.map((step, i) => (
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
            <div className="modal-review-body">今のコードはどうする？</div>
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
            <LangIcon def={currentConfig.def} />
            <span>{currentConfig.def.label}</span>
            <span className="material-symbols-outlined">expand_more</span>
          </button>
        </div>
      </header>

      {/* ===== 言語特徴パネル ===== */}
      <div className="features-bar">
        <div className="features-bar-left">
          <span className="features-bar-title">
            <LangIcon def={currentConfig.def} /> {currentConfig.def.label} でできること
          </span>
        </div>
        <div className="features-list">
          {currentConfig.features.map((f) => (
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
              <LangIcon def={currentConfig.def} /> {currentConfig.def.label}
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
              language={language}
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
                  downloadFile(codes[language][execMode]);
                  addChat({
                    from: "ai",
                    text: language === "javascript"
                      ? "ダウンロードしたよ！実行するには：\n① ターミナルを開く\n② `node program.js` と入力してEnterキーを押す"
                      : "ダウンロードしたよ！実行するには：\n① ターミナル（Windowsはコマンドプロンプト）を開く\n② `python program.py`（Macは `python3 program.py`）と入力してEnterキーを押す",
                  });
                }}>
                  <span className="material-symbols-outlined">download</span>
                  ファイルをダウンロード
                </button>
                {installGuide && (
                  <button className="btn-install-inline" onClick={() => setShowInstallModal(true)}>
                    <span className="material-symbols-outlined">build</span>
                    インストール方法
                  </button>
                )}
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
