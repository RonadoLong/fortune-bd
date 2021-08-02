package dao

//func (d *Dao) GetStrategyList(pageNum, pageSize int32) (strategyList []*model.WqStrategy) {
//	if err := d.db.Table(TABLE_WQ_STRATEGY).Where("state = 1").Order("id desc").Limit(pageSize).Offset((pageNum - 1) * pageSize).
//		Find(&strategyList).Error; err != nil {
//		logger.Errorf("GetStrategyList has err %v", err)
//		return
//	}
//	return
//}

//func (d *Dao) GetStrategy(id int64) (*model.WqStrategy, error) {
//	strategy := &model.WqStrategy{}
//	if err := d.db.Table(TABLE_WQ_STRATEGY).Where("state =1 and id = ?", id).First(strategy).Error; err != nil {
//		logger.Errorf("GetStrategy has err", err)
//		return nil, err
//	}
//	return strategy, nil
//}
