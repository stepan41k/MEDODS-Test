FROM golang:1.24.1

RUN go version

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./wait-for-postgres.sh

RUN go mod download
RUN go build -o medods-auth-app ./cmd/token/main.go
CMD ["./medods-auth-app"]