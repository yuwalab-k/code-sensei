# 競プロ教材 — Code Sensei

小学生向け AtCoder 解説集。静的サイトとして GitHub Pages に展開します。

---

## ディレクトリ構成

```
code-sensei/
├── data/                  # 問題データ（編集対象）
│   ├── index.json         # 問題一覧（ビルド時に自動更新）
│   └── problems/
│       └── typical90.json # 問題の本文・解説・コード
├── docs/                  # 自動生成される HTML（GitHub Pages 用、編集不要）
│   ├── index.html
│   ├── problems/
│   │   └── typical90_a.html
│   ├── print/
│   │   └── typical90_a.html
│   └── style.css
└── src/                   # ビルドツール（Go）
    ├── main.go
    └── go.mod
```

---

## 問題を追加する手順

1. `data/problems/typical90.json` の `problems` 配列に問題オブジェクトを追加する
2. `data/index.json` にエントリを追加する
3. ビルドコマンドを実行する

   ```bash
   cd src
   go run .
   ```

4. `git push` → GitHub Pages に反映

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
    "file": "typical90"
  }
]
```

### data/problems/typical90.json

問題の全データ。`file` の値（例: `typical90`）がファイル名になる。

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
