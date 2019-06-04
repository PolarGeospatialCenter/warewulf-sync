FROM golang:stretch

WORKDIR /go/src/warewulf-sync

COPY . ./
ENV GO111MODULE on
RUN go build -mod=readonly -o /bin/warewulf-sync

FROM scratch

COPY --from=0 /bin/warewulf-sync /bin/warewulf-sync
ENTRYPOINT ["/bin/warewulf-sync"]
