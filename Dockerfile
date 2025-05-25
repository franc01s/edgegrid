FROM golang:1.24-bullseye AS build
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends wget unzip

COPY . ./
RUN go mod download 

RUN go build -ldflags "-s -w" -o /server

FROM debian:bullseye-slim
WORKDIR /
# Install only CA certificates

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
