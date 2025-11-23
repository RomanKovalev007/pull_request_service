FROM golang:1.24.3-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY migrations /migrations

RUN go build -o main ./cmd/app

EXPOSE 8080

CMD ["./main"]