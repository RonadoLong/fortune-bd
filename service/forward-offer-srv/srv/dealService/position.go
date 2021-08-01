package dealService

import (
	"errors"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"wq-fotune-backend/libs/helper"
	"wq-fotune-backend/libs/logger"
	"wq-fotune-backend/service/forward-offer-srv/global"
	"wq-fotune-backend/service/forward-offer-srv/srv/model"
)

const updateVolumeLua = `
local positionKey = KEYS[1]
local position_val = cjson.decode(ARGV[1])
local fieldID = position_val['symbol']
local is_exists = redis.call('HEXISTS', positionKey, fieldID)
local pNumber = tonumber(ARGV[2])
if is_exists == 1 then
    local position = cjson.decode(redis.call('hget', positionKey, fieldID))
    position_val['volume'] = tonumber(position['volume']) + pNumber
	redis.call('hset', positionKey, fieldID, cjson.encode(position_val))
	return position_val['volume']
else
	position_val['volume'] = pNumber
	redis.call('hset', positionKey, fieldID, cjson.encode(position_val))
	return position_val['volume']
end`

// CalculationPosition 计算交易回调持仓
// @param changeVolume 交易持仓
// @param position 缓存持仓对象
func CalculationPosition(changeVolume float64, position *model.RedisCachePosition) (*model.RedisCachePosition, error) {
	posKey := helper.StringJoinString(global.PositionKey, position.StrategyID)
	val, _ := jsoniter.MarshalToString(position)
	script := redis.NewScript(updateVolumeLua)
	var lastVolume = float64(0)
	lastVolume, err := script.Run(global.RedisClient, []string{posKey}, val, changeVolume).Float64()
	if err != nil {
		position.Volume = changeVolume
		logStr := helper.StringJoinString("【 接收到成交回报 】保存持仓到redis失败, 会影响到策略发单，请求人工恢复持仓: ", helper.StructToJsonStr(position))
		logger.Errorf(logStr)
		return nil, errors.New(logStr)
	}
	currentVolume := lastVolume - changeVolume
	position.Volume = lastVolume
	logger.Warnf("【 接收到成交回报 】结束计算持仓量, 交易持仓为: %f 当前的持仓: %f 计算后持仓为：%f 策略ID: %s", changeVolume, currentVolume, lastVolume, position.StrategyID)
	return position, nil
}

// GetCachePosition 获取缓存持仓
func GetCachePosition(strategyID string) string {
	key := global.StringJoinString(global.PositionKey, strategyID)
	val := global.RedisClient.HGetAll(key).Val()
	if val != nil {
		return global.StructToJsonStr(val)
	}
	return ""
}
