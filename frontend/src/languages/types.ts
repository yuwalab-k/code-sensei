export type LanguageId = "python" | "javascript" | "ruby" | "php" | "perl" | "rust" | "go";
export type OS = "windows" | "mac" | "linux";

export interface FeatureExample {
  text: string;
  context: string;
}

export interface Feature {
  id: string;
  icon: string;
  title: string;
  description: string;
  detail: string;
  trivia: string;
  examples: FeatureExample[];
  ready: boolean;
}

export interface LanguageDef {
  id: LanguageId;
  label: string;
  icon: string;        // Material Symbols（フォールバック用）
  logo?: string;       // Devicon クラス名（例: "devicon-python-plain colored"）
  description: string;
  tag: string;
  available: boolean;
}

export interface InstallStep {
  title: string;
  detail: string;
  warn?: string;
}

export interface LanguageConfig {
  def: LanguageDef;
  features: Feature[];
  initialCode: string;
  /** "text": print出力 / "canvas": p5.js iframe */
  outputMode: "text" | "canvas";
  installGuide?: Record<OS, { steps: InstallStep[] }>;
}
