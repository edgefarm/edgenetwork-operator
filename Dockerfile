FROM golang:1.19 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /controller -gcflags "all=-N -l" -ldflags '-extldflags "-static"' cmd/controller/main.go

FROM alpine:3.20
COPY --from=builder /controller /
CMD ["/controller"]
