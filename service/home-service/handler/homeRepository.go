package handler

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"shop-micro/helper"
	pb "shop-micro/service/home-service/proto"
	"strings"
)

const (
	HomeHeaderKey = "HOME:HEADER"
	HomeListKey   = "HOME:LIST"
)

type HomeRepository struct {
	DB        *gorm.DB
	RedisPool *redis.Pool
}

func (repo *HomeRepository) FindHomeNav(req *pb.HomeHeaderReq, resp *pb.HomeHeadersResp) error {
	redisConn := repo.RedisPool.Get()
	defer redisConn.Close()

	exists, err := redis.Int(redisConn.Do("exists", HomeHeaderKey))
	if err != nil {
		fmt.Printf("redis db err %v", err)
		return err
	}

	if exists == 0 {
		navList, err := repo.FindHomeNavList()
		if err != nil {
			fmt.Printf("FindHomeNavList err %v", err)
			return err
		}

		for idx := range navList {
			if strings.Index(navList[idx].ImgUrl, "http") == -1 {
				navList[idx].ImgUrl = helper.Img_prefix + navList[idx].ImgUrl
			}
		}

		carousels, err := repo.FindHomeCarouselList()
		if err != nil {
			fmt.Printf("FindHomeCarouselList err %v", err)
			return err
		}
		for idx := range navList {
			if strings.Index(carousels[idx].ImgUrl, "http") == -1 {
				carousels[idx].ImgUrl = helper.Img_prefix + carousels[idx].ImgUrl
			}
		}

		resp.HomeNavList = navList
		resp.HomeCourseList = carousels

		bytes, _ := helper.MarshalToByte(resp)
		_, _ = redisConn.Do("set", HomeHeaderKey, bytes)

	} else {
		reply, _ := redisConn.Do("get", HomeHeaderKey)
		_ = helper.UnMarshal(resp, reply.([]byte))
	}
	return nil
}
