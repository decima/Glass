FROM node AS frontend
WORKDIR /app
COPY package.json /app/package.json
COPY yarn.lock /app/yarn.lock
COPY *.js /app/
RUN yarn

COPY templates /app/templates
COPY assets /app/assets
RUN mkdir -p /app/public
RUN yarn build


FROM golang:1.16.4 AS build
ENV GO111MODULE=on
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
COPY --from=frontend /app/public/build ./public/build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/glass
FROM alpine AS release
RUN apk add  --no-cache ffmpeg
ENV GIN_MODE=release
ARG VERSION=dev
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ /app/
WORKDIR /app
EXPOSE 8000
ENTRYPOINT ["/app/glass"]
CMD ["-dir=/data"]
