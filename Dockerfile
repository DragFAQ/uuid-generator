################################################
FROM golang:1.18-alpine3.17 AS builder
################################################

RUN apk add bash
SHELL ["/bin/bash", "-c"]
ARG version=v0.0.0

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

ARG mode
RUN if [ "$mode" = "debug" ]; then gcflags=(-gcflags "all=-N -l"); fi; GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
    go build -v "${gcflags[@]}" -ldflags "-X github.com/DragFAQ/uuid-generator/config.version=$version" \
    -a -installsuffix cgo -o /uuid-generator .

FROM builder AS finaldebug

FROM alpine:3.17 AS final

################################################
# Final image
FROM final${mode}
################################################
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /uuid-generator /uuid-generator

ARG mode
RUN if [ "$mode" = "debug" ]; then apk add make gcc g++ go delve bash; fi

ARG APP_USER=appuser

RUN addgroup -g 2000 "$APP_USER" && adduser -g "" -D -u 1001 -G "$APP_USER" "$APP_USER"
USER "$APP_USER"
