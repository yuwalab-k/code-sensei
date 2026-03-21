import subprocess
import tempfile
import os
import json
from http.server import BaseHTTPRequestHandler, HTTPServer

# 言語ごとの実行設定
# 新しい言語を追加するときはここに追記するだけでOK
LANGUAGE_CONFIG = {
    "python": {
        "suffix": ".py",
        "cmd": ["python"],
    },
    "javascript": {
        "suffix": ".js",
        "cmd": ["node"],
    },
    # "go": {
    #     "suffix": ".go",
    #     "cmd": ["go", "run"],
    # },
}

DEFAULT_LANGUAGE = "python"


class RunnerHandler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        pass  # デフォルトのアクセスログを抑制

    def do_POST(self):
        if self.path != "/run":
            self.send_response(404)
            self.end_headers()
            return

        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)
        data = json.loads(body)

        code = data.get("code", "")
        language = data.get("language", DEFAULT_LANGUAGE).lower()

        output = self._run_code(code, language)

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps({"output": output}).encode())

    def _run_code(self, code: str, language: str) -> str:
        config = LANGUAGE_CONFIG.get(language)
        if config is None:
            return f"エラー: 「{language}」はまだ対応していないよ"

        with tempfile.NamedTemporaryFile(
            mode="w", suffix=config["suffix"], delete=False
        ) as f:
            f.write(code)
            tmpfile = f.name

        try:
            result = subprocess.run(
                config["cmd"] + [tmpfile],
                capture_output=True,
                text=True,
                timeout=5,
            )
            if result.returncode != 0:
                return result.stderr.strip()
            return result.stdout.strip()
        except subprocess.TimeoutExpired:
            return "タイムアウト: 実行時間が長すぎます（5秒以内にしてね）"
        except Exception as e:
            return f"エラー: {e}"
        finally:
            os.unlink(tmpfile)


if __name__ == "__main__":
    server = HTTPServer(("0.0.0.0", 8002), RunnerHandler)
    print("Runner starting on 0.0.0.0:8002...")
    server.serve_forever()
