FROM golang:1.21.6

ENV GOPATH=/

WORKDIR /go/src/movie-lib
COPY . .

RUN go mod download
RUN go build -o movie-lib-app cmd/server/main.go

CMD ["./movie-lib-app --docker"]
