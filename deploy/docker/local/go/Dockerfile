FROM golang:1.26.3-alpine3.23

# インストール可能なパッケージ一覧の更新
RUN apk update && \
    apk upgrade && \
    # パッケージのインストール（--no-cacheでキャッシュ削除）
    apk add --no-cache \
            curl \
            postgresql-client

WORKDIR /go/src

COPY ./src .

# go.modがあれば依存関係をインストール
RUN if [ -f ./go.mod ]; then \
      go install; \
    fi

# 開発用のライブラリをインストール
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.0
RUN go install github.com/air-verse/air@v1.65.3
RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.12.2
RUN go install go.uber.org/mock/mockgen@v0.6.0

EXPOSE 8080
