FROM golang:1.20-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./bin/app cmd/q-quote-svc/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /
COPY .env /.env

EXPOSE 8091 8091

CMD ["/app"]





