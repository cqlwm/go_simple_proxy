# HTTP Shift

## quick start

## info

This tool forwards HTTP requests from one source to another while retaining all request and response information and data. It is commonly used for debugging and testing HTTP services, as well as forwarding HTTP requests between different environments.

## build 

```shell
docker build -t go_simple_proxy:0.0.1 .
```

## start

```shell
docker run --name go_simple_proxy -p 80:80 -d go_simple_proxy:0.0.1
```

## mock http request

```shell
curl 'http://localhost/api/test?id=123' \
  -H 'Hf-Access-Key: xxxx' \
  -H 'Re-Domain: https://www.hahaha.com' \
  --compressed
```

## request header

```properties
Hf-Access-Key=xxxx
Re-Domain=https://www.hahaha.com,https://dev.hahaha.com,https://www.hahaha2.com
Forward-Strategy=random || sequence || all
# Headers that do not need to be forwarded.
Filter-Headers=Hf-Access-Key,Re-Domain
Keep-Headers=Hf-Access-Key,Re-Domain
```