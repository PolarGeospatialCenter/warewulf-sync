FROM golang:stretch

WORKDIR /go/src/warewulf-sync

COPY . ./
ENV GO111MODULE on
RUN go build -mod=readonly -o /bin/warewulf-sync

FROM quay.io/polargeospatialcenter/warewulf3:2019.06.17.141300

COPY --from=0 /bin/warewulf-sync /bin/warewulf-sync
COPY scripts/* /usr/local/bin/
