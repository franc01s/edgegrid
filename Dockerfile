FROM golang:1.23-bullseye AS build
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends upx wget unzip

COPY . ./
RUN go mod download 

RUN go build -ldflags "-s -w" -o /server \
    && upx /server

FROM debian:bullseye-slim
WORKDIR /
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
