# Introduction
A `eventlistener` for supervisor, it may listen and redirect events to e-mail, webhook, slack and so on.  
More details http://supervisord.org/events.html

# Features
* support e-mail, webhook, slack
* support bearychat, lack(feishu)

# Usage

1. setup supervisor-eventlistener 
    ```toml
    [mail]
    receivers = ["hello@163.com", "world@163.com"]
    server_user = "test@163.com"
    server_password = "123456"
    server_host = "smtp.163.com"
    server_port = 25

    [slack]
    url = "https://hooks.slack.com/services/xxxx/xxx/xxxx"
    channel = "exception"
    timeout = 6

    [webhook]
    url = "http://my.webhook.com"
    timeout = 6

    [bearychat]
    url = "https://hook.bearychat.com/xxx/xxxx"
    channel = "alert"
    timeout = 6

    [feishu]
    url = "https://hook.feishu.com/xxx/xxxx?signKey=it_is_optional"
    timeout = 6
    ```

2. setup supervisor
    ```ini
    [eventlistener:supervisor-event-listener]
    command=/usr/local/bin/supervisor-event-listener -c /etc/supervisor-event-listener.toml
    user=nobody
    group=nobody
    events=
        PROCESS_STATE_EXITED,
        PROCESS_STATE_FATAL,
        PROCESS_STATE_STOPPED,
        PROCESS_STATE_RUNNING
    ```

3. start supervisor-eventlistener
    ```bash
    supervisorctl start supervisor-event-listener
    ```

That's all.