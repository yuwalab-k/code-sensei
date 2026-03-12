CodeSensei

小学生向けプログラミング学習環境
- バックエンド：Go（コード送信用 API）
- フロントエンド：Vite / Node.js
- Runner：Python コード実行環境(言語は拡張予定)

---

前提条件

- Docker / Docker Compose がインストール済み
- Mac / Linux / Windows で動作可能
- ポート 5173（frontend）、8081（backend）、8002（runner）が空いていること

---

プロジェクト構成

code-sensei/
├─ backend/       # Go バックエンド
├─ frontend/      # フロントエンド (Vite)
└─ runner/        # Python コード実行コンテナ

---

セットアップ手順

1. リポジトリをクローン

git clone <repository_url>
cd code-sensei

2. Docker Compose でビルド・起動

docker compose up --build -d

3. 起動確認

http://localhost:5173/

# backend が起動しているか
docker logs code-sensei-backend-1

# backend API のテスト
curl http://localhost:8081/execute -X POST -H "Content-Type: application/json" -d '{"code":"print(1+1)"}'
# → {"output":"not implemented yet"}

- ログに "Server starting on 0.0.0.0:8081..." が表示されれば正常

---

Docker ポートマッピング

サービス   | コンテナポート | ホストポート
---------- | ------------- | ------------
frontend  | 5173          | 5173
backend   | 8081          | 8081
runner    | 8002          | 8002

---

使い方

- フロントエンドから /execute エンドポイントに JSON でコードを送信可能
- 例：

{
  "code": "print(1+1)"
}

- 現在は実行結果は "not implemented yet" で返却
