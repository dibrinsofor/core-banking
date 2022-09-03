FROM golang:1.18-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o core-banking main.go
RUN apk update && apk add git

FROM alpine
COPY --from=builder /build/core-banking .

EXPOSE 8080

ENTRYPOINT [ "./core-banking"]