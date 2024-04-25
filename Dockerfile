FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/go-expense-tracker ./cmd/api/main.go

EXPOSE 80

CMD ["./bin/go-expense-tracker", "--port", "80"]
