import type { LanguageId, LanguageConfig } from "./types";
import pythonConfig from "./python";
import javascriptConfig from "./javascript";
import rubyConfig from "./ruby";

export type { LanguageId, LanguageConfig, Feature, LanguageDef, OS } from "./types";

// 言語レジストリ
// 新しい言語を追加するには、このマップにエントリを追加するだけ
export const LANGUAGE_CONFIGS: Record<LanguageId, LanguageConfig> = {
  python: pythonConfig,
  javascript: javascriptConfig,
  ruby: rubyConfig,
};

// 有効な言語リスト（言語選択UIで使用）
export const LANGUAGES = Object.values(LANGUAGE_CONFIGS)
  .filter((c) => c.def.available)
  .map((c) => c.def);

// 準備中の言語リスト（言語選択UIで「準備中」として表示）
export const COMING_SOON = [
  { label: "PHP",        icon: "public",        description: "ウェブサービスの裏側で動く言語" },
  { label: "Perl",       icon: "auto_fix_high", description: "テキスト処理が得意な言語" },
  { label: "Rust",       icon: "memory",        description: "安全で高速なシステム言語" },
  { label: "Go",         icon: "rocket_launch", description: "速くて軽い！サーバー向けの言語" },
];
