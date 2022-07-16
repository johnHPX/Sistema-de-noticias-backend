FROM golang:1.18.3 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY configs/ ./
COPY internal/ ./
COPY cmd/webapi/main.go ./

RUN go mod tidy

RUN go build -o /server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server /server
COPY --from=build /configs/ /server

EXPOSE 4083

USER nonroot:nonroot

ENTRYPOINT [ "/server" ]