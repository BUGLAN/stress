## stress

> golang实现的压力测试工具

## 安装

```golang
go get github.com/BUGLAN/stress
```

## 使用示例

```bash
stress -c 200 -n 10 -u https://www.baidu.com
```

## 指标解释

* 请求数: 当前发出的所有请求数量
* 成功数: 当前发出的所有请求数量中成功数量
* 失败数: 当前发出的所有请求数量中失败数量
* QPS: 每秒钟处理请求数量
* 最长耗时: 成功请求数量里面单个请求最长耗时
* 最短耗时: 成功请求数量里面单个请求最短耗时
* 平均耗时: 成功请求平均耗时

## 感谢

* https://github.com/link1st/go-stress-testing
* JetBrains