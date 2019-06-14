# vnt explorer 运行指南

## 下载代码
```
git clone https://github.com/vntchain/vnt-explorer.git
```

## 安装mysql服务器
1. 安装`mysql`服务器 `mysql-server-community-5.7`
2. 创建数据库`vnt`
3. 修改配置文件`conf/app.conf`

```$xslt
appname = vnt-explorer
runmode = dev
autorender = false
recoverpanic = false
copyrequestbody = true

# 配置go-vnt节点rpc服务的ip和端口
# 该节点需要提前运行起来，并且开启rpc服务
[node]
rpc_host=192.168.9.99
rpc_port=8880

# 下面是mysql的配置
[mysql]
user = "root"
pass = ""
host = "localhost"
port = "3306"
db   = "vnt"

[dev]
httpport = 8080
[prod]
httpport = 80
[test]
httpport = 8888

# 测试水龙头的配置
[hydrant]
# 这是发币账号，一般是测试网上拥有大量vnt的账户
from = "0x122369f04f32269598789998de33e3d56e2c507a"
# 账号的私钥明文，用来做发币时的签名
privateKey = "ac355731983f9ad945b642f15ed60022fa4aeb8f5c069d4f15a24c4b5100195b"
# 发币量
count = 100
# 发币间隔，同一个账号在这个间隔内只能发一次币
interval = 3600
# 测试网的chianid
chainId = 2

# 监测超级节点是否在线的配置
[supernode]
# 监测间隔，隔多少秒去检查下节点们是否在线
interval = 300
# 监测程序使用的端口
p2p_port = 5050
```

## 安装依赖
```
$ cd vnt-explorer
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
$ go get github.com/go-sql-driver/mysql
$ go get github.com/bluele/gcache
$ go get github.com/vntchain/go-vnt

# 切换一下beego依赖的分支
$ cd $GOPATH/src/github.com/astaxie/beego
$ git checkout -b v1.11.1 v1.11.1
$ cd -
```

### 编译代码
```bash
$ ./build.sh
```
编译完成后，目录下会多出一个`bin/`目录，分别包含了六个可执行文件：
* dbsync 数据库初始化命令，用于初始化表格
* feixiaohao 获取交易所数据，用于同步当前vnt价格，市值，交易量等数据
* geniuser 创世块同步工具
* nodemonitor 超级节点检测程序，用于检查超级节点及候选节点是否在线
* racer 数据同步命令，用于将区块链中的数据同步到mysql数据库中
* vnt-explorer 浏览器后端服务，用于开启后端服务

## 运行浏览器后端进程

该后端为前端提供接口

```
$ nohup ./bin/vnt-explorer > vnt-explorer.log 2>&1 &
```

## 测试
```
# 获取区块列表
GET /v1/blocks
```
## 运行数据同步模块

该模块会连接go-vnt节点，并将数据同步到mysql数据库中

```
# 初始化数据库表，只需要在初始化时运行一次即可，此命令慎用
$ ./bin/dbsync
```
> 警告：如果你是在原有的系统上升级，则不要进行该步骤

## 运行数据迁移
本次升级包含了数据库字段的变更，并且新增了一些新表，所以需要通过迁移脚本完成这些表结构的自动升级。

如果本次安装是一次全新的安装，则不用运行该步骤
```
$ bee migrate -conn="root:root@tcp(127.0.0.1:3306)/vnt" -dir=database/migrations
```
> 警告：如果你要进行一次全新的安装，则不要进行该步骤


## 开启同步数据进程
```
$ nohup ./bin/racer > racer.log 2>&1 &
```

## 开启同步交易所数据进程
```
$ nohup ./bin/feixiaohao > feixiaohao.log 2>&1 &
```

## 开启节点检测进程
```
$ nohup ./bin/nodemonitor > nodemonitor.log 2>&1 &
```

# 工具
### 1. 创世块同步工具 geniuser
> 该工具用于将创世块的allocation数据同步为交易和账户数据
> 该工具一般建议在浏览器数据同步前执行一次，也可以在后面随时执行

运行命令：
```bash
$ ./bin/geniuser
```