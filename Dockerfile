#########################################
# Build stage
#########################################
FROM golang:1.16 AS builder

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=za

# Copy sources into the container
ADD . /go/src/$REPOSITORY/$ARTIFACT

# Set working directory
WORKDIR /go/src/$REPOSITORY/$ARTIFACT

# Build the binary
RUN make build

#########################################
# Distribution stage
#########################################
FROM gcr.io/distroless/base-debian10

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=za

# Install binary
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT /usr/local/bin/$ARTIFAC6

# Exposed ports
EXPOSE 8080 9213

# Define entrypoint
ENTRYPOINT [ "za" ]