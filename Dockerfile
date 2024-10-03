FROM golang:1.20.10-alpine AS build
LABEL authors="t-stark"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pay_system .

# Create new image for the runtime
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/pay_system .
COPY --from=build /app/.env .

EXPOSE 3534

CMD ["./pay_system"]