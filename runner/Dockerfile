FROM python:3.11-slim

# Node.js・Ruby・PHP・Perl・Go・Rust をインストール
RUN apt-get update && apt-get install -y nodejs ruby php perl golang-go && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain stable
ENV PATH="/root/.cargo/bin:${PATH}"

WORKDIR /app

COPY run.py .

CMD ["python", "run.py"]
