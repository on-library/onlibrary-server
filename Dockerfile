FROM golang:alpine


WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o binary

EXPOSE 8080

ENTRYPOINT ["/app/binary"]