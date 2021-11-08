# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.17 as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.sum/go.mod and warm up the module cache (so that this
# rather long step can be cached if go.mod/go.sum don't change)
COPY go.* $D/
CMD go mod download

# Now copy the rest.
COPY . /app/

# Build the binary.
RUN cd /app/ && CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/savannahghi/mycarehub

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates
# add timezone then test
RUN apk add tzdata && cp /usr/share/zoneinfo/Africa/Nairobi /etc/localtime
RUN echo "Africa/Nairobi" >  /etc/timezone && date
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
COPY --from=builder /app/deps.yaml /deps.yaml

# Run the web service on container startup.
CMD ["/server"]
