FROM golang:1.26-alpine AS builder

LABEL author="canureal"

WORKDIR /app

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src .

RUN go build -o meMessages .

# this image has the CA certificates that's why we're using that
FROM gcr.io/distroless/static

COPY --from=builder /app/meMessages /meMessages

EXPOSE 4444

ENTRYPOINT ["/meMessages"]
