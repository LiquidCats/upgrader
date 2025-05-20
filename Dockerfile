FROM golang:1.24.2-alpine AS build

WORKDIR /app

ADD ./ /app

RUN go build -o main ./cmd/upgrader/main.go

FROM scratch

EXPOSE 8080

COPY --from=build /app/main /main

CMD ["/main"]



