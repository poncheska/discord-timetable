FROM golang:alpine as builder

RUN apk update && apk upgrade && apk add --no-cache git

RUN mkdir -p /github.com/poncheska/discord-timetable
WORKDIR /github.com/poncheska/discord-timetable

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /github.com/poncheska/discord-timetable/app .

CMD ["./app"]