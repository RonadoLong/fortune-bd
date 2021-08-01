package v1

import (

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/render"
	"wq-fotune-backend/service/grid-strategy-srv/model"
)

// CreateStrategyType 添加策略类型
func CreateStrategyType(c *gin.Context) {
	form := &gridType{}

	err := render.BindJSON(c, form)
	if err != nil {
		render.Err400Msg(c, "解析json错误, "+err.Error())
		return
	}

	err = form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	gt := form.toGridType()
	err = gt.Insert()
	if err != nil {
		logger.Error("insert data failed", logger.Err(err), logger.Any("form", form))
		render.Err500(c, "保存数据失败")
		return
	}

	model.SetStrategyTypeCache(gt.Type, gt)

	render.OK(c)
}

// DelStrategyType 删除策略类型
func DeleteStrategyType(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		render.Err400Msg(c, "field Name is empty")
		return
	}

	query := bson.M{"_id": bson.ObjectIdHex(id)}
	gs, _ := model.FindStrategyType(query, bson.M{})

	_, err := model.DeleteStrategyType(query)
	if err != nil {
		logger.Error("delete data failed", logger.Err(err), logger.String("id", id))
		render.Err500(c, "删除数据失败")
		return
	}

	if gs.Name != "" {
		model.DelStrategyTypeCache(gs.Type)
	}

	render.OK(c)
}

// ListStrategyTypes 获取策略类型列表
func ListStrategyTypes(c *gin.Context) {
	form := &reqListForm{
		pageStr:  c.Query("page"),
		limitStr: c.Query("limit"),
		sort:     c.Query("sort"),
	}

	err := form.valid()
	if err != nil {
		logger.Error("参数错误", logger.Err(err), logger.Any("form", form))
		render.Err400Msg(c, err.Error())
		return
	}

	gts, err := model.FindStrategyTypes(bson.M{}, bson.M{}, form.page, form.limit, "type")
	if err != nil {
		logger.Error("获取数据失败", logger.Err(err), logger.String("page", form.pageStr), logger.String("limit", form.limitStr))
		render.Err500(c, err.Error())
		return
	}

	total, err := model.CountStrategyTypes(bson.M{})
	if err != nil {
		logger.Error("获取数据失败", logger.Err(err))
		render.Err500(c, err.Error())
		return
	}

	render.OK(c, gin.H{"gridTypes": convert2Values(gts), "total": total})
}
