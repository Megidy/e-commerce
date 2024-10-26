FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod","go.sum" ,"./"]

RUN go mod download

COPY . . 

RUN go build -o ./api ./cmd/main.go

FROM alpine

WORKDIR /

COPY --from=builder /app/api /api
COPY --from=builder /app/frontend /frontend 

CMD ["/api"]