FROM golang:1.7.5 as build

RUN mkdir -p /go/src/github.com/openfaas-incubator/faas-rancher/

WORKDIR /go/src/github.com/openfaas-incubator/faas-rancher

COPY vendor     vendor
COPY handlers	handlers
COPY types      types
COPY rancher     rancher
COPY mocks mocks
COPY server.go  .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" \
  && go test $(go list ./... | grep -v integration | grep -v /vendor/ | grep -v /template/) -cover \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faas-rancher .

FROM alpine:3.5
RUN apk --no-cache add ca-certificates
WORKDIR /root/

EXPOSE 8080
ENV http_proxy      ""
ENV https_proxy     ""

COPY --from=build /go/src/github.com/openfaas-incubator/faas-rancher/faas-rancher    .

CMD ["./faas-rancher"]
