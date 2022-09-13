# benchmark-proxy
a proxy tool for HTTP/HTTPS benchmark,  use it by curl -x 
# CA
it's required flag ca-crt and ca-key before running the service,generate ca file could be like this way:
https://www.cnblogs.com/lab-zj/p/15176787.html
# example 
> curl -x 127.0.0.1:9900 http://www.baidu.com -H 'Benchmark-Proxy-Type:'times'' -H 'Benchmark-Proxy-Times:100' -H 'Benchmark-Proxy-Concurrency:100'
# CustomHeaders

| Header                      | Meaning                                               |
|-----------------------------|-------------------------------------------------------|
| Benchmark-Proxy-Type        | proxy type, it can only one of "times" and "duration" |
| Benchmark-Proxy-Times       | run http requests times                               |
| Benchmark-Proxy-Duration    | run http requests duration                            |
| Benchmark-Proxy-Concurrency | concurrency in running                                |

# data sequence
![alt 数据流图](./doc/benchmark-proxy.png)