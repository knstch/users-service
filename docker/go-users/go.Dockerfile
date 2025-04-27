FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build ./cmd/users

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/users /build/.env ./
COPY --from=builder /build/users ./

EXPOSE ${PUBLIC_HTTP_ADDR}
CMD ["/app/users"]