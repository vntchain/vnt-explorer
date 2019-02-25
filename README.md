# vnt explorer 运行指南

## 下载代码
```
git clone https://github.com/vntchain/vnt-explorer.git
```

## 安装mysql服务器
1. 安装mysql服务器
2. 创建数据库vnt
3. 修改conf/app.conf中的mysql相关配置信息

## 安装依赖
```
cd vnt-explorer
go get github.com/astaxie/beego
go get github.com/beego/bee
go get github.com/go-sql-driver/mysql
```

## 运行
```
bee run
```

## 测试
```
# 创建一个Test
POST /v1/test -d {"Name": "TestName"}

# 获取一个Test
GET /v1/test/1
```