# quick start


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