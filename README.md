#  mall-micro 基于单体应用重构更新中...
#  申明: 个人作品 仅供学习
 * 基于go-micro微服务实战，涉及视频资讯，电商支付，广告发布等需求
 * [配套mall-app 基于React native 开发Android,iOS APP](https://github.com/TorettoLong/mall-app)
 * 基于Vue和springboot构建的后台管理系统 [前端代码](https://github.com/TorettoLong/mall-admin) [API代码](https://github.com/TorettoLong/mall-admin-java)


## 待加功能
  * 使用rabbitmq中间件, 实现延迟队列
  * 添加grpc重试，限流，熔断
  * .....


## 运行步骤
 * sh build.sh
 * docker-compose build
 * docker-compose up
 
##功能点:
  * 登陆: 短信 微信 facebook 
  * 视频：视频播放
  * 新闻: 新闻
  * 发布广告
  * 上传图片七牛
  * 支付宝支付
  * 银联支付....

## 前端部分页面展示
<img src="https://image.showm.xin/phone/test/01.png" width="375px" height="667px">

<img src="https://image.showm.xin//test/04.png" width="375px">

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

