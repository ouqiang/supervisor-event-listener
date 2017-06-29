# supervisor-event-listener
Supervisor事件通知, 支持邮件, Slack, WebHook

## 简介
Supervisor是*nix环境下的进程管理工具, 可以把前台进程转换为守护进程, 当进程异常退出时自动重启.  
supervisor-event-listener监听进程异常退出事件, 并发送通知.
  
## 下载
[v1.0](https://github.com/ouqiang/supervisor-event-listener/releases)

### 源码安装
* `go get -u github.com/ouqiang/supervisor-event-listener`

## Supervisor配置
```ini
[eventlistener:supervisor-event-listener]
; 默认读取配置文件/etc/supervisor-event-listener.ini
command=/path/to/supervisor-event-listener
; 指定配置文件路径
;command=/path/to/supervisor-event-listener -c /path/to/supervisor-event-listener.ini
events=PROCESS_STATE_EXITED
......
```

## 配置文件, 默认读取`/etc/supervisor-event-listener.ini`

```ini 
[default]
# 通知类型 mail,slack,webhook 只能选择一种
notify_type = mail

# 邮件服务器配置
mail.server.user = test@163.com
mail.server.password = 123456
mail.server.host = smtp.163.com
mail.server.port = 25

# 邮件收件人配置, 多个收件人, 逗号分隔
mail.user = hello@163.com

# Slack配置
slack.webhook_url = https://hooks.slack.com/services/xxxx/xxx/xxxx
slack.channel = exception

# WebHook通知URL配置 
webhook_url = http://my.webhook.com

```

## 通知内容
邮件、Slack
```shell
Host: ip(hostname)
Process: process-name
PID: 6152
EXITED FROM state: RUNNING
```
WebHook, Post请求, 字段含义查看Supervisor文档
```json
{
  "Header": {
    "Ver": "3.0",
    "Server": "supervisor",
    "Serial": 11,
    "Pool": "supervisor-listener",
    "PoolSerial": 11,
    "EventName": "PROCESS_STATE_EXITED",
    "Len": 84
  },
  "Payload": {
    "Ip": "ip(hostname)",
    "ProcessName": "process-name",
    "GroupName": "group-name",
    "FromState": "RUNNING",
    "Expected": 0,
    "Pid": 6371
  }
}
```
