# Step 1: Modules caching
FROM golang:1.21.5-alpine as modules

LABEL maintainer="Suraj Bobade"

# git is required to fetch go dependencies
RUN apk add --no-cache ca-certificates git

# Create the user and group files that will be used in the running
# container to run the process as an unprivileged user.
# RUN mkdir /user && \
#    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
#    echo 'nobody:x:65534:' > /user/group

COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21.5-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./src/cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/myApp ./cmd/kafka

# Step 3: Final
FROM scratch
COPY --from=builder /bin/myApp /myApp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/myApp"]