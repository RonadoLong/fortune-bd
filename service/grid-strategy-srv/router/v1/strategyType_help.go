package v1

import (
	"errors"
	"wq-fotune-backend/service/grid-strategy-srv/model"

	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/logger"
)

type gridType struct {
	ID       string   `json:"id"`       // 策略id
	Type     int      `json:"type"`     // 0:网格交易，1:趋势网格，2:无限网格，3:反向网格
	Name     string   `json:"name"`     // 网格名称
	Labels   []string `json:"labels"`   // 标签
	Describe string   `json:"describe"` // 说明
	RunCount int      `json:"runCount"` // 运行数量
}

func (t *gridType) valid() error {
	if t.Name == "" {
		return errors.New("field Name is empty")
	}
	if t.Type > 5 {
		return errors.New("field Type is illegality")
	}
	if len(t.Labels) == 0 {
		return errors.New("field Labels is empty")
	}
	return nil
}

func (t *gridType) toGridType() *model.StrategyType {
	return &model.StrategyType{
		Type:     t.Type,
		Name:     t.Name,
		Labels:   t.Labels,
		Describe: t.Describe,
	}
}

func convert2Values(sts []*model.StrategyType) []*gridType {
	var gts []*gridType
	for _, v := range sts {
		gts = append(gts, &gridType{
			ID:       v.ID.Hex(),
			Type:     v.Type,
			Name:     v.Name,
			Labels:   v.Labels,
			Describe: v.Describe,
			RunCount: getStrategyTypeCount(v.Type),
		})
	}
	return gts
}

func getStrategyTypeCount(strategyType int) int {
	query := bson.M{"type": strategyType, "isRun": true}

	n, err := model.CountGridStrategies(query)
	if err != nil {
		logger.Warn("model.CountGridStrategies error", logger.Err(err), logger.Any("query", query))
	}

	return n
}
