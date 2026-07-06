# 競プロ教材 — Code Sensei

小学生向け AtCoder 解説集。静的サイトとして GitHub Pages に展開します。

---

## ディレクトリ構成

```
code-sensei/
├── data/                  # 問題データ（編集対象・参照用ソース）
│   ├── index.json         # 問題一覧（ビルド時に自動更新）
│   └── problems/
│       └── typical90_a.json など # 問題ごとの本文・制約・解答コード
├── docs/                  # GitHub Pages に展開される HTML
│   ├── index.html         # 自動生成（編集不要）
│   ├── glossary/          # 自動生成（編集不要）
│   ├── code_reading.html  # 自動生成（編集不要）
│   ├── style.css          # 自動生成（編集不要）
│   └── problems/
│       └── typical90_a.html など # 手動作成のスライドショーページ（要コミット・要保管）
└── src/                   # ビルドツール（Go）
    ├── main.go
    └── go.mod
```

> `docs/problems/*.html` はゲーム風UI・マスコット・実行環境つきの手作りページです。
> `go run . -f`（強制ビルド）で消えるので、必ず Git にコミットして保管してください。
> `data/problems/*.json` は問題文・サンプル・模範解答などの**参照用ソース**で、
> 実際のページ制作時は [PROMPT.md](PROMPT.md) の構成に沿って `docs/problems/*.html` を直接作成・編集します。

---

## 問題を追加する手順

1. `data/problems/typical90_x.json`（問題ごとに1ファイル）に問題オブジェクトを追加する
2. `data/index.json` に `"file": "typical90_x"` でエントリを追加する
3. [PROMPT.md](PROMPT.md) の構成に沿って、`docs/problems/typical90_x.html` をゲーム風スライドショーとして直接作成する（マスコット・ロボット含む）
4. ビルドコマンドを実行する（`docs/index.html` や `docs/glossary/*` など自動生成ページのみ更新される。手作りした `docs/problems/*.html` は上書きされない）

   ```bash
   cd src
   go run .
   ```

5. `git push` → GitHub Pages に反映

---

## ビルドコマンド

```bash
cd src
go run .           # JSON → HTML（既存ファイルはスキップ）
go run . -f        # JSON → HTML（全件強制再生成）
go run . serve     # ローカル確認 http://localhost:8080
```

---

## JSON の構造

### data/index.json

問題一覧ページ用の軽量インデックス。

```json
[
  {
    "id": "typical90_a",
    "contest": "典型90問",
    "problem": "001",
    "title": "Yokan Party",
    "difficulty": 4,
    "tags": ["二分探索", "貪欲"],
    "file": "typical90_a"
  }
]
```

### data/problems/typical90_x.json

問題の参照用データ（本文・制約・サンプル・模範解答）。`file` の値（例: `typical90_a`）がファイル名になる。
実際に表示されるゲーム風ページ（`docs/problems/typical90_x.html`）は、このデータを元に
[PROMPT.md](PROMPT.md) の構成に沿って別途手作りする。

```json
{
  "contest": "典型90問",
  "problems": [
    {
      "id": "typical90_a",
      "problem": "001",
      "title": "Yokan Party",
      "atcoder_url": "https://atcoder.jp/contests/typical90/tasks/typical90_a",
      "difficulty": 4,
      "tags": ["二分探索", "貪欲"],
      "statement": "問題文（Markdown）",
      "constraints": "制約（Markdown）",
      "samples": [
        {
          "input": "3 34\n1\n8 13 26",
          "output": "13",
          "explanation": "解説（Markdown）"
        }
      ],
      "statement_note": "入力例1を使った説明（Markdown）",
      "easy_explanation": "やさしい解説（Markdown）",
      "explanation": "くわしい解説（Markdown）",
      "solutions": {
        "python": { "code": "..." },
        "cpp":    { "code": "..." }
      },
      "bad_solutions": {
        "python": {
          "code": "...",
          "label": "O(L×N) 線形探索 — 二分探索を使わず全部試す"
        }
      },
      "added_at": "2026-05-20"
    }
  ]
}
```

### difficulty の基準

| 値 | 表示 |
|----|------|
| 1  | A問題レベル |
| 2  | B問題レベル |
| 3  | C問題レベル |
| 4  | D問題レベル |

### solutions / bad_solutions で使えるキー

`python` / `cpp` / `typescript` / `ruby` / `php` / `rust` / `perl`

---

## GitHub Pages の設定（初回のみ）

1. GitHub リポジトリ → Settings → Pages
2. Source: **Deploy from a branch**
3. Branch: `main` / `docs`
