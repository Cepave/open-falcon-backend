falcon-sender
=============

alarm处理报警event可能会产生报警短信或者报警邮件，alarm不负责发送，只是把报警邮件、短信写入redis队列，sender负责读取并发
送。

各个公司有自己的短信通道，自己的邮件发送方式，sender如何调用各个公司自己的组件呢？那只能制定规范了，sender的配置文件
cfg.json中配置了api:sms和api:mail，即两个http接口，这是需要各个公司提供的。

当要发送短信的时候，sender就会调用api:sms中配置的http接口，post方式，参数是：

- tos：用逗号分隔的多个手机号
- content：短信内容

当要发送邮件的时候，sender就会调用api:mail中配置的http接口，post方式，参数是：

- tos：用逗号分隔的多个邮箱地址
- content：邮件正文
- subject：邮件标题

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/open-falcon
cd $GOPATH/src/github.com/open-falcon
git clone https://github.com/open-falcon/sender.git
cd sender
go get ./...
./control build
# vi cfg.json modify configuration
./control start
```

## Configuration

- redis: redis地址需要和alarm、judge使用同一个
- queue: 维持默认即可，需要和alarm的配置一致
- worker: 最多同时有多少个线程玩命得调用短信、邮件发送接口
- api: 短信、邮件发送的http接口，各公司自己提供

## How to debug

想知道 sender 是否可以正常運作，需要去查看 Redis 的狀態
```
redis-cli monitor
```

假設 redis 裡的 /sms 佇列是用來產生 sms 的 queue，用如下的方式可以直接在 redis 裡塞入資料。

```
redis-cli

127.0.0.1:6379> LPUSH /sms "{\"tos\":\"phone number\",\"content\":\"redis testing\"}"
```
