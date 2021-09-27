FROM golang:1.17 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN go mod tidy -compat=1.17

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o dhcpserver

FROM alpine:3 as alpine
RUN apk add --no-cache ca-certificates

FROM scratch
# Copy ca-certs from alpine
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/dhcpserver /dhcpserver

COPY config /etc/coredhcp/config.yml

# Run the web service on container startup.
ENTRYPOINT ["/dhcpserver"]
