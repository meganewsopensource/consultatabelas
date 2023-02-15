FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN apk update && apk add git

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

ENV TZ="America/Bahia"
ENV CONNHTTP = https://portalunico.siscomex.gov.br/classif/api/publico/nomenclatura/download/json
ENV CONNSTRING=postgres://admin:admin@localhost:5432/tabelas

CMD ["./main"]
