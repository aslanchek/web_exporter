FROM mirror.gcr.io/golang:1.23-alpine AS builder

ARG GOPROXY

WORKDIR /app
COPY . .
RUN go build -o=/bin/web_exporter /app/cmd/web_exporter

FROM mirror.grc.io/alpine
COPY --from=builder /bin/web_exporter /bin/web_exporter

CMD ["/bin/web_exporter"]