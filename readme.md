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
填写对应信息

其中id可以使用 `ddns-go -l` 查询(需要先填写完其他信息)

## start

```
sudo systemctrl enable ddns-go
sudo systemctrl start ddns-go
```
