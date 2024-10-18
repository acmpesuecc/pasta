FROM golang:1.23.2-alpine3.20

WORKDIR /pasta

COPY go.mod ./
COPY main.go ./

RUN go build -o pasta

EXPOSE 8080
EXPOSE 80

CMD ["./pasta"]
