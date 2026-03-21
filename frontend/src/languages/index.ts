import type { LanguageId, LanguageConfig } from "./types";
import pythonConfig from "./python";
import javascriptConfig from "./javascript";
import rubyConfig from "./ruby";
import phpConfig from "./php";
import perlConfig from "./perl";
import rustConfig from "./rust";
import goConfig from "./go";

export type { LanguageId, LanguageConfig, Feature, LanguageDef, OS } from "./types";

// 言語レジストリ
// 新しい言語を追加するには、このマップにエントリを追加するだけ
export const LANGUAGE_CONFIGS: Record<LanguageId, LanguageConfig> = {
  python: pythonConfig,
  javascript: javascriptConfig,
  ruby: rubyConfig,
  php: phpConfig,
  perl: perlConfig,
  rust: rustConfig,
  go: goConfig,
};

// 有効な言語リスト（言語選択UIで使用）
export const LANGUAGES = Object.values(LANGUAGE_CONFIGS)
  .filter((c) => c.def.available)
  .map((c) => c.def);

// 準備中の言語リスト（言語選択UIで「準備中」として表示）
export const COMING_SOON: { label: string; icon: string; description: string }[] = [];
