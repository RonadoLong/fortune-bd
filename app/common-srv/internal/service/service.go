package service

import (
	"wq-fotune-backend/app/common-srv/cache"
	"wq-fotune-backend/app/common-srv/internal/dao"
)

func NewCommonService() *CommonService {
	handler := &CommonService{dao: dao.New(), cacheSrv: cache.NewService()}
	return handler
}