# 競プロ教材 — Code Sensei

小学生向け AtCoder 解説集。静的サイトとして GitHub Pages に展開します。

---

## ディレクトリ構成

```
code-sensei/
├── problems-raw/          # 問題ファイル（Markdown）
│   └── typical90_a.md
├── data/                  # 自動生成される JSON（コミット対象）
│   ├── index.json
│   └── problems/
│       └── typical90.json
├── docs/                  # 自動生成される HTML（GitHub Pages 用）
│   ├── index.html
│   └── problems/
│       └── typical90_a.html
└── src/                   # ビルドツール（Go）
    ├── main.go
    └── go.mod
```

---

## 問題を追加する手順

1. AtCoder の問題ページのテキストをコピーして `problems-raw/<id>.md` に保存する
   - 命名規則: `abc300_a.md`、`typical90_a.md` など
   - フォーマット例は既存ファイルを参照

2. AI（Claude など）に `.md` を貼り付けて解説と解答コードを書いてもらう

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
go run .              # convert + build（通常はこれだけ）
go run . convert      # md → JSON のみ
go run . build        # JSON → HTML（既存ファイルはスキップ）
go run . build -f     # JSON → HTML（全件強制再生成）
go run . serve        # ローカル確認 http://localhost:8080
```

---

## Markdown ファイルの書き方

```
---
contest: 典型90問
atcoder_url: https://atcoder.jp/contests/typical90/tasks/typical90_a
difficulty: 4
tags: [二分探索, 貪欲]
added_at: 2026-05-20
---

# 001 - Yokan Party（★4）

## 問題文
...

## 制約
...

## 入力例 1
\```text
...
\```

## 出力例 1
\```text
...
\```

# 解説
## 方針
...

# Python 解答
\```python
...
\```

# C++ 解答
\```cpp
...
\```
```

### difficulty の基準

| 値 | 意味 |
|----|------|
| 1  | A問題レベル |
| 2  | B問題レベル |
| 3  | C問題レベル |
| 4  | D問題レベル |

---

## GitHub Pages の設定（初回のみ）

1. GitHub リポジトリ → Settings → Pages
2. Source: **Deploy from a branch**
3. Branch: `main` / `docs`
