# Overview

a proxy tool for HTTP/HTTPS benchmark, use it by curl -x
benchmark provides:

- HTTP/HTTPS proxy
- benchmark for api request
- statistics spend time,success count and error count of requests
- replace data(query's param or body) by tag in proxy request

# Requirement

it's required flag ca-crt and ca-key(using RSA algorithm) before running benchmark-proxy,generate ca file could be like
this
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

# ReplaceTag

| Tag       | Example                                                     |
|-----------|-------------------------------------------------------------|
| ${uuid}   | d035581b-53a3-48e5-9461-ba24709f06c9                        |
| ${int}    | 6331615752200874333                                         |
| ${float}  | 0.681078                                                    |
| ${string} | 295dfd92fcd9cd9e43cfa5b2b87e806dda83eb3d7dfd97d5ef          |
| ${incr}   | 1(default:1,it will auto increment in every proxy request ) |

# Architecture

![alt 数据流图](./doc/benchmark-proxy.png)