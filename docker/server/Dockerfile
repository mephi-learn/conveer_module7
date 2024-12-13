FROM golang:alpine AS builder
WORKDIR /src
ADD cmd cmd
ADD proto proto
COPY go.* .
RUN CGO_ENABLED=0 go build -o /src/bin/server cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /src/bin/server /app/server
EXPOSE 53000/tcp 53001/tcp
ENTRYPOINT ["/app/server"]