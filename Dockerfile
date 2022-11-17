FROM golang:alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN apk --no-cache -U add libc-dev build-base ca-certificates

RUN go build -ldflags "-linkmode external -extldflags -static" -o cfo .

RUN if [ ! -f .env ]; then cp .env.example .env


FROM scratch

COPY --from=builder /app/cfo ./cfo
COPY --from=builder /app/.env ./.env
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./cfo"]
