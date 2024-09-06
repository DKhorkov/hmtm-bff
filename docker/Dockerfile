# https://habr.com/ru/companies/otus/articles/660301/
FROM golang:alpine AS builder
LABEL authors='dkhorkov'

WORKDIR /build

COPY . .

RUN go build -o hmtmbff ./cmd/hmtmbff/main.go

FROM alpine AS runner

WORKDIR /app

COPY .env .

COPY --from=builder /build/hmtmbff /app/hmtmbff

CMD ["./hmtmbff"]





#WORKDIR /app
#
#COPY . /app
#
#RUN go build -o hmtmbff ./cmd/hmtmbff/main.go
#CMD ["./hmtmbff"]
