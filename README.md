CodeSensei

小学生向けプログラミング学習環境
- バックエンド：Go（コード送信用 API）
- フロントエンド：Vite / Node.js
- Runner：Python コード実行環境(言語は拡張予定)

---

# 1. コンテナ起動
docker compose up --build

# 2. 別ターミナルでモデルをダウンロード（初回のみ・数分かかる）
docker compose exec ollama ollama pull qwen2.5-coder:7b