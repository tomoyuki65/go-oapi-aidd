####################
# ビルドステージ
####################
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /go/src

COPY ./src .

# 依存関係をインストール
RUN go install

# ビルド（GOOS：OS指定、GOARCH：CPUアーキテクチャ指定）
RUN GOOS=linux GOARCH=amd64 go build -o main .

####################
# 実行ステージ
####################
FROM alpine:3.23 AS runner

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

# インストール可能なパッケージ一覧の更新
RUN apk update && \
    apk upgrade && \
    # パッケージのインストール（--no-cacheでキャッシュ削除）
    apk add --no-cache \
            tzdata

WORKDIR /app

# コンテナ用ユーザー作成
RUN addgroup --system --gid 1001 appuser && \
    adduser --system --uid 1001 appuser

# ビルドステージで作成したバイナリをコピー
COPY --from=builder --chown=appuser:appuser ./go/src/main .

# ポートを設定
EXPOSE 8080

# コンテナ起動ユーザー設定
USER appuser

# APIサーバー起動コマンド
CMD ["./main"]
