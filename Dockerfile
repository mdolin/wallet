# Golang base image
FROM golang:alpine as builder

# Working dir
WORKDIR /go/src/github.com/mdolin/wallet

# Copy mod and sum files
COPY go.mod go.sum /go/src/github.com/mdolin/wallet/

# Download dependencies
RUN go mod download

# Copy the source from the PWD to the working directory inside the container
COPY . /go/src/github.com/mdolin/wallet

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/wallet github.com/mdolin/wallet

# Start new stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Copy the pre-build binary file from the previous stage
COPY --from=builder /go/src/github.com/mdolin/wallet/build/wallet /usr/bin/wallet

# Expose port to the outside world
EXPOSE 8000 8000

# Set executables that will run when the container is initiated
ENTRYPOINT ["/usr/bin/wallet"]