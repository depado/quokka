# Build Step
FROM golang:1.25.2-alpine@sha256:06cdd34bd531b810650e47762c01e025eb9b1c7eadd191553b91c9f2d549fae8 AS builder

# Dependencies
RUN apk update && apk add --no-cache make git

# Source
WORKDIR $GOPATH/src/github.com/depado/quokka
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN make tmp

# Final Step
FROM gcr.io/distroless/static@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c
COPY --from=builder /tmp/qk /go/bin/qk
ENTRYPOINT ["/go/bin/qk"]
