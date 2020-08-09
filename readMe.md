多渠道通道通知，
redis 0库 logQueue 发布订阅模式
如果没有传type 默认为钉钉推送 或者解析json错误，或者类型为钉钉 会将相关数据发送到钉钉，
publish logQueue '{"name":"hello","level":1,"time":"daads","message":"程序deb--ug"}'
name string 名称
level int  可选值 1 2 3 1、程序错误，需立即解决 2、程序错误，稍后解决！ 3、提示信息 
time string 日期时间格式
message string 相关信息
type 目前可选值 dingding