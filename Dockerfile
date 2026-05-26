ARG GO_VERSION=1.26.2
ARG ALPINE_VERSION=3.22.1


# ...
FROM golang:${GO_VERSION}-bookworm AS base
ENV GOCACHE=/root/.cache/go-build
WORKDIR /widget
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        "musl-dev=1.2.3-1" \
        "musl-tools=1.2.3-1" \
    && rm -rf /var/lib/apt/lists/*


# ...
FROM base AS dependencies
COPY go.* ./
RUN go mod download


# ...
FROM dependencies AS build
COPY . .
WORKDIR /widget/cmd/lastfm-now-playing
RUN --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=1 CC=musl-gcc \
    go build \
        --ldflags '-linkmode=external -extldflags="-static"' \
        -tags exec_local,save_s3 \
        -o ../../widget .


# ...
FROM public.ecr.aws/lambda/provided:al2.2024.05.24.17-x86_64 AS lambda-runtime-base
WORKDIR /var/task
RUN curl https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm -O \
    && yum install -y google-chrome-stable_current_x86_64.rpm \
    && rm google-chrome-stable_current_x86_64.rpm \
    && yum clean all
COPY assets ./assets
COPY entry_script.sh /entry_script.sh
ENTRYPOINT ["/entry_script.sh"]


# ...
FROM lambda-runtime-base AS lambda-runtime-local
RUN curl https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/download/v1.19/aws-lambda-rie-x86_64 -O \
    && chmod +x aws-lambda-rie-x86_64 \
    && mv aws-lambda-rie-x86_64 /usr/local/bin/rie
COPY --from=build /widget/widget ./widget


# ...
FROM lambda-runtime-base AS lambda-runtime-aws
COPY --from=build /widget/widget ./widget


# ...
FROM alpine:${ALPINE_VERSION} AS local-runtime-base
WORKDIR /
RUN apk update \
    && apk add --no-cache chromium \
    && rm -rf /var/cache/apk/*


# ...
FROM local-runtime-base AS local-runtime
COPY assets ./assets
COPY --from=build /widget/widget ./widget
ENV CHROMIUM_BROWSER_BINARY_PATH=/usr/bin/chromium
ENTRYPOINT ["/widget"]


# ...
FROM local-runtime-base AS goreleaser-local-runtime
COPY assets ./assets
COPY lastfm-now-playing ./widget
ENV CHROMIUM_BROWSER_BINARY_PATH=/usr/bin/chromium
ENTRYPOINT ["/widget"]
