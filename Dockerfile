FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY vendor/ vendor/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -trimpath -ldflags="-s -w" -o /helloworld ./cmd/helloworld

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /helloworld /helloworld

USER 65532:65532

ENTRYPOINT ["/helloworld"]
