# benchmark-proxy
a proxy tool for HTTP/HTTPS benchmark,  use it by curl -x 
# example 
> curl -x 127.0.0.1:9900 http://www.baidu.com -H 'Benchmark-Proxy-Type:'times'' -H 'Benchmark-Proxy-Times:100' -H 'Benchmark-Proxy-Concurrency:100'
# CustomHeaders

| Header                      | Meaning                                               |
|-----------------------------|-------------------------------------------------------|
| Benchmark-Proxy-Type        | proxy type, it can only one of "times" and "duration" |
| Benchmark-Proxy-Times       | run http requests times                               |
| Benchmark-Proxy-Duration    | run http requests duration                            |
| Benchmark-Proxy-Concurrency | concurrency in running                                |
