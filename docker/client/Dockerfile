FROM golang:alpine AS builder
WORKDIR /src
ADD cmd cmd
ADD proto proto
COPY go.* .
RUN CGO_ENABLED=0 go build -o /src/bin/client cmd/client/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /src/bin/client /app/client
EXPOSE 55000/tcp
ENTRYPOINT ["/app/client", "server", "init"]