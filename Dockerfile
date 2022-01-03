# syntax=docker/dockerfile:1

FROM golang:1.17.5 as builder

WORKDIR /smp-server/build

COPY go.mod .
COPY go.sum .

# download dependencies
RUN go mod download

# copy sources
COPY . .

# build the program and move the executable to parent folder
RUN GOOS=linux go build -a -installsuffix cgo -o smp

#FROM alpine:latest
FROM debian:11-slim

WORKDIR /smp-server/
COPY --from=builder /smp-server/build/smp ./
RUN mkdir config
COPY --from=builder /smp-server/build/config/config.yaml  ./config/

# default port
EXPOSE 1984

ENTRYPOINT [ "./smp" ]
