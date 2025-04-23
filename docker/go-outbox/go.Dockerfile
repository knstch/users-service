FROM golang:1.23 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build ./cmd/cargo

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/cargo /build/.env ./
COPY --from=builder /build/cargo ./
COPY certs ./certs/

EXPOSE ${HTTP_PORT}
CMD ["/app/cargo"]