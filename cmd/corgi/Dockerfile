FROM golang:1.17.6 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /corgi /app/cmd/corgi/main.go

EXPOSE 8081

CMD [ "/corgi" ]
