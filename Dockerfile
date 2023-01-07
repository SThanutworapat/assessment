FROM golang:1.19.1-alpine as build-base

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .



#RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/go-expenses .

# =====================================

FROM alpine:3.16.2

COPY --from=build-base /app/out/go-expenses /app/go-expenses

ENV DATABASE_URL=postgres://gfoydwkf:Xc1f1ENOVnCNy-_OxdWG48Kq4oxii6x9@tiny.db.elephantsql.com/gfoydwkf

ENV PORT=2565

CMD ["/app/go-expenses"]