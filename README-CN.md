# 接收 mqtt 消息写入 TDengine

作为消费者接入 MQTT 消息,根据配置解析 json 写入 TDengine

[English](README.md) | 简体中文

## 编译方法
1. 安装golang(1.14+) [https://golang.google.cn/doc/install](https://golang.google.cn/doc/install)
2. 设置 golang 代理 
```shell
go env -w GOPROXY=https://goproxy.cn,direct
```
3. 安装 TDengine 客户端 [https://www.taosdata.com/en/getting-started/](https://www.taosdata.com/en/getting-started/)
4. 配置 C 编译环境(GCC)
5. 在本项目目录执行
```shell
go build
```

## MQTT 配置

```json
{
  "address": "mqtt 地址",
  "clientID": "客户端ID 如果不设置则使用 uuid",
  "username": "用户名",
  "password": "密码",
  "keepAlive": 30,
  "caPath": "ca 证书路径",
  "certPath": "证书路径",
  "keyPath": "证书 key 路径"
}
```
`keepAlive` 保持存活时间,单位秒
## TDengine 配置

```json
{
  "host": "地址",
  "port": 6030,
  "user": "用户名",
  "password": "密码",
  "db": "数据库"
}
```

`port` 为 TDengine 服务端口

## 解析规则配置

```json
[
  {
    "ruleName": "规则名称",
    "topic": "主题",
    "rule": {
      "sTable": "对应 STable 名称",
      "table": {
        "defaultValue": "默认值",
        "path": "json 路径"
      },
      "tags": [
        {
          "name": "对应 TDengine 中的 tag 名",
          "valueType": "值的类型",
          "length": "值最大长度(在值类型为 string 时需要设置)",
          "defaultValue": "默认值",
          "path": "json 路径",
          "timeLayout": "时间格式化的布局(在值类型为 timeString 时需要设置)"
        }
      ],
      "columns": [
        {
          "name": "对应 TDengine 中的 column 名",
          "valueType": "值的类型",
          "length": "值最大长度(在值类型为 string 时需要设置)",
          "defaultValue": "默认值",
          "path": "json 路径",
          "timeLayout": "时间格式化的布局(在值类型为 timeString 时需要设置)"
        }
      ]
    }
  }
]
```
* 默认值: 在 json 中未找到 path 对应的值时使用默认值
* 值的类型: `"int"
  "float"
  "bool"
  "string"
  "timeString"
  "timeSecond"
  "timeMillisecond"
  "timeMicrosecond"
  "timeNanosecond"`
* json 路径  见[https://github.com/tidwall/gjson](https://github.com/tidwall/gjson)
* 当值类型为 `string` 时必须设置 `length` 参数
* `tags` 至少设置一个
* `columns` 至少设置两个,第一个参数名称一定是 `ts` 且类型为 `"timeString" "timeSecond" "timeMillisecond" "timeMicrosecond" "timeNanosecond"` 中的一个
* `timeLayout` 为 golang 时间格式化模板 [https://golang.google.cn/pkg/time/#pkg-constants](https://golang.google.cn/pkg/time/#pkg-constants)
## 参数
```
  -c string
        配置文件路径 (默认 "./config/config.json")
  -rc string
        规则配置文件路径 (默认 "./config/rule.json")
```
## 配置样例
见 [example folder](./example)