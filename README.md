English | [中文](README_ZH.md)

# Overview

a proxy tool for HTTP/HTTPS benchmark test

- HTTP/HTTPS proxy
- benchmark for api request
- statistics spend time,success count and error count of requests
- replace data(query's param or body) by tag in proxy request
- check request api response data and status for diagnose success or fail

# Optional

self-provision certificates is optional while using HTTPS

1. `openssl genrsa -out ca.key 2048`
2. `openssl req -new -x509 -key ca.key -out ca.crt -days 1095`
3. `move ca.* resources/`
4. Update the system's trusted certificates

# Run

> benchmark-proxy --port {port} --ca-crt {file_path} --ca-key {file_path}

# Request Example

> curl -x 127.0.0.1:9900 https://www.baidu.com -H 'Benchmark-Proxy-Times:1' -H 'Benchmark-Proxy-Concurrency:1'

# CustomHeaders

| Header                              | Meaning                                                              |
|-------------------------------------|----------------------------------------------------------------------|
| Benchmark-Proxy-Times               | indicate how many times exec in each http request                    |
| Benchmark-Proxy-Duration            | indicate how much second exec in each http requests                  |
| Benchmark-Proxy-Concurrency         | concurrency in running                                               |
| Benchmark-Proxy-Check-Result-Status | indicate the response status to determine whether request is success |
| Benchmark-Proxy-Check-Result-Body   | indicate the response body to determine whether request is success   |

### Response Checker

example :

1. `Benchmark-Proxy-Check-Result-Status: 200` Indicate http request is success that response status is 200
2. `Benchmark-Proxy-Check-Result-Body: hello world` Indicate http request is success return that http response body is '
   hello world'
3. `Benchmark-Proxy-Check-Result-Body: @Reg[\w]` Indicate http request is success return that http response body is
   satisfied with provider regexp rule

# ReplaceTag

| Tag       | Example                                                     |
|-----------|-------------------------------------------------------------|
| ${uuid}   | d035581b-53a3-48e5-9461-ba24709f06c9                        |
| ${int}    | 77                                                          |
| ${float}  | 0.94                                                        |
| ${string} | 762edb6805                                                  |
| ${incr}   | 1(default:1,it will auto increment in every proxy request ) |

# Architecture

![alt 数据流图](./doc/benchmark-proxy.png)