FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix csgo -o main .

FROM scratch

COPY --from=builder /build/main /go/main

EXPOSE 8084

ENTRYPOINT [ "/go/main" ]