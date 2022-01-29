FROM golang:1.17.6 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /corgi /app/cmd/corgi/main.go

FROM scratch

COPY --from=builder /corgi /corgi
COPY rbac_model.conf /
COPY rbac_policy.csv /

EXPOSE 8081

CMD [ "/corgi" ]
