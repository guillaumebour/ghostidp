FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN env GOOS=linux GOARCH=amd64 go build -o /go/bin/ghostidp cmd/idp/main.go


FROM gcr.io/distroless/base-debian12

COPY --from=builder /go/bin/ghostidp /

EXPOSE 8080/tcp

CMD ["/ghostidp"]
