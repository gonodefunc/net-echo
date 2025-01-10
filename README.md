# Quick and dirty TCP/UDP-based echo server and client in Go

TCP/UDP 端口测试工具，基于 https://github.com/bediger4000/udp-echo-server 修改，支持多端口多协议测试

## Building

```sh
$ make build
```

## Usage

In Window 1:

    $ ./port-server :: udp UDP-PORT1,UDP-PORT2 tcp TCP-PORT1,TCP-PORT2

In Window 2:

    $ ./port-client IP udp UDP-PORT1,UDP-PORT2 tcp TCP-PORT1,TCP-PORT2

Connected:
```sh
[ACK-UDP:10000]-[STATUS: CONNECTED]-[MSG: PONG]
[ACK-UDP:10001]-[STATUS: CONNECTED]-[MSG: PONG]
[ACK-TCP:20001]-[STATUS: CONNECTED]-[MSG: PONG]
[ACK-TCP:20000]-[STATUS: CONNECTED]-[MSG: PONG]
```
Failed:
```sh
[ACK-UDP:10001]-[STATUS: FAILED]
[ACK-UDP:10000]-[STATUS: FAILED]
[ACK-TCP:20000]-[STATUS: FAILED]
[ACK-TCP:20001]-[STATUS: FAILED]
```