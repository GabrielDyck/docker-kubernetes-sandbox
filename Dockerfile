
FROM golang:1.12-alpine

COPY . .

WORKDIR main/

RUN go build main.go

EXPOSE 9290

CMD ["./main"]



