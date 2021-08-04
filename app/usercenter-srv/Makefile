#PKG_LIST := $(shell /usr/local/go/bin/go list ./... | grep -v /vendor/)
#PC_PKG_LIST := $(shell /usr/local/go/bin/go list ./... | grep -v /vendor/)
#GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

test:
	go test

lint:
	golint

# 编译程序
build:
	bash ./scripts/build-app.sh

# 构建镜像
build-image:
	bash ./scripts/build-image.sh $(REGISTRY_HOST) $(TAG)

# 推送镜像到仓库
push-image:
	bash ./scripts/push-image.sh $(REGISTRY_HOST) $(TAG)

# 清除本地镜像
delete-image:
	bash ./scripts/delete-images.sh

# 部署到k8s集群
deploy:
	bash ./scripts/deploy.sh

# 重命名项目名称
rename:
	bash ./scripts/rename.sh ./scripts $(OLDNAME) $(NEWNAME)
