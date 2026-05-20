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

1. `PROMPT.md` を開いて末尾に AtCoder の問題ページのテキストを貼り付ける

2. `PROMPT.md` の内容ごと Claude などの AI に渡す

3. AI が出力した内容を `problems-raw/<id>.md` として保存する
   - 命名規則: `abc300_a.md`、`typical90_a.md` など

4. ビルドコマンドを実行する

   ```bash
   cd src
   go run .
   ```

5. `git push` → GitHub Pages に反映

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

# TypeScript 解答
\```typescript
...
\```

# Ruby 解答
\```ruby
...
\```

# PHP 解答
\```php
...
\```

# Rust 解答
\```rust
...
\```

# Perl 解答
\```perl
...
\```
```

### Markdown 記法のポイント

| 書き方 | 表示 | 用途 |
|--------|------|------|
| `` `L` `` | `L`（コード） | 変数名・記号 |
| `*L*` | *L*（斜体） | 数学的な変数 |
| `**太字**` | **太字** | 強調 |
| `` ```python ... ``` `` | コードブロック | 解答コード |

問題文中の変数は `` `L` `` と書けば等幅で表示されます。  
数学的に斜体にしたい場合は `*L*` と書いてください。

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
