# ベースとなるイメージを指定
FROM golang:1.17

# 作業ディレクトリを設定
WORKDIR /app

# ホストのファイルをコンテナにコピー
COPY . .

# 依存関係をダウンロード
RUN go mod download

# アプリケーションをビルド
RUN go build -o main .

# コンテナを実行する際のコマンドを指定
CMD ["./main"]


