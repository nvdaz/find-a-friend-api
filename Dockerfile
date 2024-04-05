FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o ./find-a-friend

EXPOSE 8080

CMD [ "./find-a-friend" ]
