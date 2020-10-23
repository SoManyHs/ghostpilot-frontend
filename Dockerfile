FROM golang:1.15.2 as builder

# Copy all the source files for the frontend service.
RUN mkdir /svc
ADD . /svc

# We specify that we now wish to execute any further commands inside the /svc directory.
WORKDIR /svc

# Build the binary
ENV GOPROXY=direct
RUN go build -o frontend ./cmd/frontend


# For the real image, we'll only copy the binaryso that the image size is small.

EXPOSE 8080
CMD ["./frontend"]