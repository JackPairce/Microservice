FROM golang:1.22.3-alpine3.20

WORKDIR /app

RUN apk add --no-cache make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .

CMD ["make", "run-server"]