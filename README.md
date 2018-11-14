# micro-web 更新中...
基于go-micro 微服务实战 配套Android iOS 前端

#运行步骤
```
 1. make build
 2. docker-compose build
 3. docker-compose up

```
#部分页面展示
<img src="https://image.showm.xin//test/01.png" width="375px">

##govendor命令	功能

```
init	初始化 vendor 目录
list	列出所有的依赖包
add	添加包到 vendor 目录，如 govendor add +external 添加所有外部包
add PKG_PATH	添加指定的依赖包到 vendor 目录
update	从 $GOPATH 更新依赖包到 vendor 目录
remove	从 vendor 管理中删除依赖
status	列出所有缺失、过期和修改过的包
fetch	添加或更新包到本地 vendor 目录
sync	本地存在 vendor.json 时候拉去依赖包，匹配所记录的版本
get	类似 go get 目录，拉取依赖包到 vendor 目录

```

##组件
```$xslt
gin web         github.com/gin-gonic/gin
gorm msyql      github.com/jinzhu/gorm
redis           github.com/go-redis/redis
file setting    github.com/BurntSushi/toml
uuid            github.com/HaroldHoo/id_generator
gin pprof       github.com/DeanThompson/ginpprof
auth            gopkg.in/dgrijalva/jwt-go.v3
tail            github.com/hpcloud/tail
logs            github.com/astaxie/beego/logs
kafka           github.com/Shopify/sarama
elasticSearch   github.com/olivere/elastic
json            github.com/json-iterator/go
go-torch        github.com/uber/go-torch
```
