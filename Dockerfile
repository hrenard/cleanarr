FROM golang as builder
WORKDIR /app
COPY . .
RUN  go get -d -v \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine
COPY --from=builder /app/cleanarr .
ENTRYPOINT ["./cleanarr"]