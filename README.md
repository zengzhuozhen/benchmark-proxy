# Overview

a proxy tool for HTTP/HTTPS benchmark, use it by curl -x
benchmark provides:

- HTTP/HTTPS proxy
- benchmark for api request
- statistics spend time,success count and error count of requests

# Requirement

it's required flag ca-crt and ca-key(using RSA algorithm) before running benchmark-proxy,generate ca file could be like this
way:
https://www.cnblogs.com/lab-zj/p/15176787.html

# Run
> benchmark-proxy --port {port} --ca-crt {file_path} --ca-key {file_path}

# Request Example
> curl -x 127.0.0.1:9900 http://www.baidu.com -H 'Benchmark-Proxy-Type:'times'' -H 'Benchmark-Proxy-Times:100' -H '
> Benchmark-Proxy-Concurrency:100'

# CustomHeaders

| Header                      | Meaning                                               |
|-----------------------------|-------------------------------------------------------|
| Benchmark-Proxy-Type        | proxy type, it can only one of "times" and "duration" |
| Benchmark-Proxy-Times       | run http requests times                               |
| Benchmark-Proxy-Duration    | run http requests duration                            |
| Benchmark-Proxy-Concurrency | concurrency in running                                |

# Architecture

![alt 数据流图](./doc/benchmark-proxy.png)