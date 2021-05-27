# ddns-go

基于golang和腾讯云api开发的ddns服务

## install

``` bash
make
sudo make install
```

## config

```
sudo vim /etc/ddns-go/config.yaml
```

## config参数

### Type

域名解析商家 目前支持腾讯云和阿里云

* qcloud 腾讯云
* aliyun 阿里云

### SecretId SecretKey
密钥信息

### Domain
解析的域名

### Record
解析的记录

### Cron
控制定时查询的参数

在线生成工具 https://www.matools.com/cron


## start

```
sudo systemctl enable ddns-go
sudo systemctl start ddns-go
sudo systemctl status ddns-go
```
