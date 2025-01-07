# Build Step
FROM golang:1.23-alpine@sha256:d37127f39271451047bcd91fc53ee014829603c96b91d02ff65ab3a7d1fb3c5e AS builder

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
FROM gcr.io/distroless/static@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58
COPY --from=builder /tmp/qk /go/bin/qk
ENTRYPOINT ["/go/bin/qk"]
