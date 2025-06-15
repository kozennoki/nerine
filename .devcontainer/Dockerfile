# 最新LTSのNode.jsを使用
FROM node:22

# Claude Codeをインストール
RUN npm install -g @anthropic-ai/claude-code

# Go 1.24をインストール
RUN curl -fsSL https://go.dev/dl/go1.24.0.linux-amd64.tar.gz | tar -xzC /usr/local
ENV PATH="/usr/local/go/bin:${PATH}"

# Taskをインストール
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# 作業ディレクトリを設定
WORKDIR /workspace
