# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY .env ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /main
EXPOSE 8000

CMD [ "/main" ]

ENTRYPOINT [ "/main" ]