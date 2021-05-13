FROM golang:1.15.12-alpine3.12

RUN mkdir /app

ADD . /app

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["/app/main"]