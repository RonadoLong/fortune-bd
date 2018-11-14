package handler

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	pb "shop-micro/service/home-service/proto"
	"shop-web/common/cost"
	"strings"
)

const (
	HomeHeaderKey = "HOME:HEADER"
	HomeListKey = "HOME:LIST"
)

type HomeRepository struct {
	DB *gorm.DB
	RedisPool *redis.Pool
}

func (repo *HomeRepository) FindHomeNav(req *pb.HomeNavListReq, resp *pb.HomeNavListResp) (error) {
	redisConn := repo.RedisPool.Get()
	defer redisConn.Close()

	exists, err := redis.Int(redisConn.Do("exists", HomeHeaderKey))
	if err != nil || exists == 0 {
		fmt.Printf("redis db err %v", err)
		return err
	}

	navList, err := repo.FindHomeNavList()
	if err != nil{
		fmt.Printf("FindHomeNavList err %v", err)
		return err
	}

	for idx := range navList {
		if strings.Index(navList[idx].ImgUrl, "http") == -1 {
			navList[idx].ImgUrl = cost.Img_prefix + navList[idx].ImgUrl
		}
	}

	resp.HomeNavLists = nil
	fmt.Printf("nav %v", navList)
	return nil
}

