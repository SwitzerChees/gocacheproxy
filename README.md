# Reverse HTTP Proxy In Go With InMemory Caching

## Features

1. Reverse Proxy the one target
2. Cache filetypes
3. Cache specific pages
4. Docker Image available ([Dockerfile](https://github.com/SwitzerChees/gocacheproxy/blob/master/Dockerfile))

## Getting Started

These instructions will cover usage information about the go application and for the docker container image. This reverse proxy only supports HTTP (no HTTPS support for target). Speed up TTFB up to 1 - 15 ms.

### Prerequisities

In order to run this container you'll need docker installed.

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)

### Usage

#### Container Parameters

Start a simple reverse proxy container to localhost:8080

```shell
docker run -p 80:80 -e PROXY_PORT=80 -e PROXY_TARGET=http://localhost:8080 -e CONFIG_FILE= /config/cache.config --name gocacheproxy switzerchees/gocacheproxy
```

Link the reverse proxy with a webapp for caching between

```shell
docker run -p 80:80 -e PROXY_TARGET=http://wordpress:80 --link wordpress:wordpress --name gocacheproxy switzerchees/gocacheproxy
```

## Enpoints

1. `/healthcheck` - Check if the webserver is running (httpstatus 200 success)
2. `/flushcache` - flush the actual cache
3. `/cacheactive` - activate or deactivate the cache functionality (POST allows toggle)

## Environment Variables

* `PROXY_PORT` - Port of the Reverse Proxy (Default: 80
* `PROXY_TARGET` - Target to redirect traffic (Default: http://localhost:8080)
* `CONFIG_FILE` - Path to the configuration files for define files to cache (Default: /config/cache.config)
* `PROXY_CACHE_HTTP_HEADER` - Header to Check if File comes from cache (Default: X-GoProxy: FromCache im file comes from the reverse proxy cache)

## Configuration File

* `{mimetype}image/jpeg` - The files with the ending mime-type jpeg are cached
* `{page}/` - With the prefix {page} you can specify a specific page or resource to cache

## Configuration File Default

* {mimetype}text/css
* {mimetype}application/javascript
* {mimetype}text/javascript
* {mimetype}image/jpeg
* {mimetype}image/gif
* {mimetype}image/png
* {mimetype}image/svg+xml
* {page}/
