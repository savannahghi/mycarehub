# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.14 as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.sum/go.mod and warm up the module cache (so that this
# rather long step can be cached if go.mod/go.sum don't change)
COPY go.* $D/
CMD go mod download

# Now copy the rest.
COPY . /app/

# Set up the credentials needed to fetch private code
ARG CI_JOB_TOKEN
RUN git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.slade360emr.com".insteadOf "https://gitlab.slade360emr.com"

# Retrieve application dependencies.
RUN GOPRIVATE="gitlab.slade360emr.com/go/*" go mod download

# Build the binary.
RUN cd /app/ && CGO_ENABLED=0 GOOS=linux GOPRIVATE="gitlab.slade360emr.com/go/*" go build -v -o server gitlab.slade360emr.com/go/profile

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
COPY --from=builder /app/deps.yaml /deps.yaml

# Run the web service on container startup.
CMD ["/server"]
