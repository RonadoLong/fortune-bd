package handler

import (

	"shop-micro/service/home-service/proto"
)

func (repo * HomeRepository) FindHomeNavList() ([]*shop_srv_home.HomeNav, error) {
	var homeNavList []*shop_srv_home.HomeNav
	err := repo.DB.Table("home_nav").
		Where("`status` = 1").
		Order("`sort` desc").
		Find(&homeNavList).Error

	if err != nil{
		return homeNavList, err
	}
	return homeNavList, nil
}

func (repo * HomeRepository) FindHomeCarouselList() ([]*shop_srv_home.HomeCourse,error){
	var homeCarousel []*shop_srv_home.HomeCourse
	err := repo.DB.Table("home_carousel").
		Where("`status` = 1").
		Order("`sort` desc").
		Find(&homeCarousel).Error

	if err != nil{
		return homeCarousel, err
	}
	return homeCarousel, nil
}
