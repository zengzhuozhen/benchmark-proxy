中文 | [English](README.md)

# 简介

benchmark-proxy 是一个用于 HTTP/HTTPS 接口基准性能测试的代理工具。更具体的，它是:

- 一个 HTTP/HTTPS 代理服务
- 用于API接口基准测试/并发测试
- 能够统计测试接口的请求成功数，失败数，耗时
- 能够采用标签替换方式传递变量数据
- 检查接口返回状态码和数据来判断成功或失败

# 可选

 使用 HTTPS 时可以选择生成自己的CA文件

1. `openssl genrsa -out ca.key 2048`
2. `openssl req -new -x509 -key ca.key -out ca.crt -days 1095`
3. `move ca.* resources/`
4. 更新系统信任证书

# 运行

> benchmark-proxy --port {port}

# 请求示例

> curl -x 127.0.0.1:9900 https://www.baidu.com -H 'Benchmark-Proxy-Times:1' -H 'Benchmark-Proxy-Concurrency:1'

# 自定义请求头

| Header                              | Meaning                 |
|-------------------------------------|-------------------------|
| Benchmark-Proxy-Times               | 声明被测接口的请求次数             |
| Benchmark-Proxy-Duration            | 声明被测接口的请求时间（秒）          |
| Benchmark-Proxy-Concurrency         | 声明请求被测接口的并发度（协程数）       |
| Benchmark-Proxy-Check-Result-Status | 声明请求被测接口成功时返回的 HTTP 状态码 |
| Benchmark-Proxy-Check-Result-Body   | 声明请求被测接口成功时返回的 HTTP 数据  |

### 响应校验

示例：

1. `Benchmark-Proxy-Check-Result-Status: 200` —— 响应状态码为 200 时判定为成功
2. `Benchmark-Proxy-Check-Result-Status: 200,201` —— 响应状态码为 200 或 201 时判定为成功
3. `Benchmark-Proxy-Check-Result-Status: 200-299` —— 响应状态码在 [200,299] 区间时判定为成功
4. `Benchmark-Proxy-Check-Result-Body: hello world` —— 响应 body 等于 'hello world' 时判定为成功
5. `Benchmark-Proxy-Check-Result-Body: @Contains[success]` —— 响应 body 包含 'success' 时判定为成功
6. `Benchmark-Proxy-Check-Result-Body: @Reg[\w+]` —— 响应 body 满足正则 `\w+` 时判定为成功

**支持的状态码校验表达式：**
- `200`（等值）
- `200,201`（多值）
- `200-299`（区间）

**支持的 body 校验表达式：**
- `hello world`（等值）
- `@Contains[success]`（包含）
- `@Reg[pattern]`（正则）

# CustomizeTag
在请求数据（body、header、URL参数）中可使用如下标签，程序会自动替换为真实值。

**支持的标签及用法示例：**

| 标签             | 示例 / 用法                        | 说明                         |
|------------------|------------------------------------|------------------------------|
| ${uuid}          | d035581b-53a3-48e5-9461-ba24709f06c9 | 随机 UUID                   |
| ${int}           | 77                                 | 随机整数                     |
| ${int:10,20}     | 15                                 | 10~20 间随机整数             |
| ${float}         | 0.94                               | 随机浮点数                   |
| ${float:1.5,2.5} | 2.01                               | 1.5~2.5 间随机浮点数         |
| ${string}        | 762edb6805                         | 随机字符串（默认10位）        |
| ${string:8}      | 1a2b3c4d                           | 指定长度的随机字符串         |
| ${bool}          | true                               | 随机布尔值                   |
| ${date:2006-01-02}| 2024-05-01                        | 当前日期，指定格式           |
| ${timestamp}     | 1714550400                         | 当前时间戳                   |
| ${incr}          | 1                                  | 自增整数（默认从1开始）       |
| ${incr:100,2}    | 100,102,104...                     | 从100开始，步长为2自增       |
| ${range:1,5}     | 1,2,3,4,5                          | 区间自增 [1,5]               |
| ${list:[a,b,c]}  | b                                  | 从列表中随机选取一个          |
| ${const:hello}   | hello                              | 常量                         |

这些标签可用于请求体、请求头、URL参数等任意位置，每次请求都会动态生成。

# 架构图

![alt 数据流图](./doc/sequence-diagram.png)