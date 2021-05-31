# ddns-go

基于golang开发的动态域名解析服务，适用于动态公网解析

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

密钥获取的地址如下
* 腾讯云 https://console.cloud.tencent.com/cam/capi
* 阿里云 https://ram.console.aliyun.com/manage/ak

请创建密钥或者使用已有的密钥

> 请妥善管理密钥信息

SecretId SecretKey对应的信息如下

|服务商|SecretId|SecretKey|
| - | - | - |
|腾讯云|SecretId|SecretKey|
|阿里云|AccessKey ID|AccessKey Secret|

### Domain Record
解析的域名和记录名

假设需要动态修改IP的域名如为

```
sub.test.com
|_| |______|
|          |
Record     Domain
```

Domain为`test.com`

Record为`sub`

> 该条记录请先自行创建 `ddns-go` 只进行解析记录的修改 不会创建记录 

### Cron
控制定时查询的参数

在线生成工具 https://www.matools.com/cron

## start

```
sudo systemctl enable ddns-go
sudo systemctl start ddns-go
sudo systemctl status ddns-go
```
