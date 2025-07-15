FROM golang:1.24-alpine

RUN apk add --no-cache git bash netcat-openbsd

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy
RUN go mod download

COPY . .

COPY wait-for-db.sh /app/wait-for-db.sh
RUN chmod +x /app/wait-for-db.sh

EXPOSE 8080

CMD ["/app/wait-for-db.sh", "db", "go", "run", "main.go"]
