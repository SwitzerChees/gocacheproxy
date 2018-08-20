# Reverse HTTP Proxy In Go With In Memory Caching

## Features

1. Reverse Proxy the one target
2. Cache filetypes
3. Cache specific pages
4. Docker Image compatible

## Getting Started

These instructions will cover usage information about the go application and for the docker container image. This reverse proxy only supports HTTP (no HTTPS support for target)

### Prerequisities

In order to run this container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

### Usage

#### Container Parameters

Start a simple reverse proxy container to localhost:8080

```shell
docker run -e PROXY_PORT=80 -e PROXY_TARGET=http://localhost:8080 -e CONFIG_FILE= /config/cache.config switzerchees/reverse-cache-proxy
```

## Enpoints

1. `/healthcheck` - Check if the webserver is running (httpstatus 200 success)
2. `/flushcache` - flush the actual cache
3. `/flushcache` - activate or deactivate the cache functionality (POST allows toggle)

## Environment Variables

* `PROXY_PORT` - Port of the Reverse Proxy (Default: 80
* `PROXY_TARGET` - Target to redirect traffic (Default: http://localhost:8080)
* `CONFIG_FILE` - Path to the configuration files for define files to cache (Default: /config/cache.config)
* `PROXY_CACHE_HTTP_HEADER` - Header to Check if File comes from cache (Default: X-GoProxy: FromCache im file comes from the reverse proxy cache)