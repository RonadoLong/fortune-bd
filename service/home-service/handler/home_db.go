package handler

import (
	"shop-micro/service/home-service/model"
)

func (repo * HomeRepository) FindHomeNavList() ([]model.HomeNav, error) {
	var homeNavList []model.HomeNav
	err := repo.DB.Table("home_nav").
		Where("`status` = 1").
		Order("`sort` desc").
		Find(&homeNavList).Error

	if err != nil{
		return homeNavList, err
	}
	return homeNavList, nil
}

func (repo * HomeRepository) FindHomeCarouselList() ([]model.HomeCarousel,error){
	var homeCarousel []model.HomeCarousel
	err := repo.DB.Table("home_carousel").
		Where("`status` = 1").
		Order("`sort` desc").
		Find(&homeCarousel).Error

	if err != nil{
		return homeCarousel, err
	}
	return homeCarousel, nil
}
