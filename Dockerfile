# Build Step
FROM golang:1.27.3-alpine@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS builder

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
FROM gcr.io/distroless/static@sha256:448810357cc8e67b7208988abdba0f5ec3d3aebee0a5b13eb624979792cffc0f
COPY --from=builder /tmp/qk /go/bin/qk
ENTRYPOINT ["/go/bin/qk"]
