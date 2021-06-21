FROM golang:1.16.4 AS build
ENV GO111MODULE=on
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/glass
FROM scratch AS release
ENV GIN_MODE=release
ARG VERSION=dev
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ /app/
WORKDIR /app
EXPOSE 8000
ENTRYPOINT ["/app/glass"]
CMD ["-dir=/data"]
