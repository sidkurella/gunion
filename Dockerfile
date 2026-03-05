FROM golang:1.24-alpine

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/gunion /usr/local/bin/gunion

ENTRYPOINT ["gunion"]
