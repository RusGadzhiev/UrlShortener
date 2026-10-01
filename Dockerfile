FROM golang:1.26 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./url_shortener ./cmd/url_shortener

FROM alpine:latest AS runner

COPY --from=builder /app/url_shortener .
COPY --from=builder /app/config.yaml ./config.yaml

EXPOSE 8000 5432

CMD [ "./url_shortener" ]

