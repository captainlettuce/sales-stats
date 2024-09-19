FROM golang:1.23-alpine

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cortesi/modd/cmd/modd@latest

CMD ["modd"]
