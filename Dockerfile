ARG ARCH=amd64
# build stage
FROM golang:1.19 AS builder
RUN mkdir -p /go/src/meteo
WORKDIR /go/src/meteo
COPY . ./
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build -a -o /app .


# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app ./
RUN chmod +x ./app
ENTRYPOINT ["./app"]
EXPOSE 80