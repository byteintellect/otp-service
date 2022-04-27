# build the server binary
FROM golang:1.17.1-alpine3.14 AS builder

ARG GIT_USERNAME
ARG GIT_PASSWORD
LABEL stage=server-intermediate
RUN apk add --no-cache ca-certificates git openssh alpine-sdk openssl
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group


# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine github.com\n\
    login ${GIT_USERNAME}\n\
    password ${GIT_PASSWORD}" >> /root/.netrc \
RUN chmod 600 /root/.netrc
RUN cat /root/.netrc

RUN mkdir -p ~/.ssh

WORKDIR /go/src/otp_svc
ADD . /go/src/otp_svc
RUN go env -w GOPRIVATE=github.com/byteintellect/\*
RUN go mod tidy
RUN go build -o bin/server ./cmd/server

# copy the server binary from builder stage; run the server binary
FROM alpine:latest AS runner
WORKDIR /bin

# Go programs require libc
RUN mkdir -p /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /go/src/otp_svc/bin/server .
COPY --from=builder /go/src/otp_svc/config/* /bin/
ENTRYPOINT ["server"]
