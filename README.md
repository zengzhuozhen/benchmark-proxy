# Overview

a proxy tool for HTTP/HTTPS benchmark, use it by curl -x
benchmark provides:

- HTTP/HTTPS proxy
- benchmark for api request
- statistics spend time,success count and error count of requests
- replace data(query's param or body) by tag in proxy request

# Requirement

To use HTTPS proxy, you must generate CA files.
1. `openssl genrsa -out ca.key 2048`
2. `openssl req -new -x509 -key ca.key -out ca.crt -days 1095`
3. Update the system's trusted CA certificates


# Run

> benchmark-proxy --port {port} --ca-crt {file_path} --ca-key {file_path}

# Request Example

> curl -x 127.0.0.1:9900 http://www.baidu.com -H 'Benchmark-Proxy-Type:'times'' -H 'Benchmark-Proxy-Times:100' -H '
> Benchmark-Proxy-Concurrency:100'

# CustomHeaders

| Header                              | Meaning                                                              |
|-------------------------------------|----------------------------------------------------------------------|
| Benchmark-Proxy-Times               | indicate how many times exec in each http request                    |
| Benchmark-Proxy-Duration            | indicate how much second exec in each http requests                  |
| Benchmark-Proxy-Concurrency         | concurrency in running                                               |
| Benchmark-Proxy-Check-Result-Status | indicate the response status to determine whether request is success |
| Benchmark-Proxy-Check-Result-Body   | indicate the response body to determine whether request is success   |

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