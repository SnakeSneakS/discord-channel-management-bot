FROM golang:1.19 as development
WORKDIR /work
CMD ["go","run","main.go"]

FROM golang:1.19 as builder
WORKDIR /work
COPY . /work
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /work/app
ENTRYPOINT [ "/work/app" ]

FROM alpine:3 as runner
WORKDIR /bin
COPY --from=builder /work/app /bin/app
ENTRYPOINT ["/bin/app"]

