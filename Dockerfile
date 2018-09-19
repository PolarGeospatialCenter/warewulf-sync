FROM golang:stretch

WORKDIR /go/src/github.com/PolarGeospatialCenter/warewulf-sync

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . ./
RUN dep ensure -vendor-only
RUN go build -o /bin/warewulf-sync

FROM scratch

COPY --from=0 /bin/warewulf-sync /bin/warewulf-sync
ENTRYPOINT ["/bin/warewulf-sync"]
