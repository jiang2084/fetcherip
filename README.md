## 说明
抓取免费的代理ip，并保存到数据库中
## 安装及使用
```
git clone git@github.com:jiang2084/fetcherip.git

cd fetcherip

进入config,配置好mysql等配置文件

mkdir bin
go build -o bin/fetchip
cd bin
./fetchip start -c ../conf/config.yml
```
