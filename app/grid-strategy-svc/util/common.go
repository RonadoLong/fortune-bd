package util

//
//func ChangeNeedCount(needCount float64) float64 {
//	splitArr := strings.Split(fmt.Sprint(needCount), ".")
//	if len(splitArr) == 2 {
//		var roundLen = len(splitArr[1])
//		var subLen = int32(roundLen - 1)
//		needCount, _ = decimal.NewFromFloat(needCount).Round(subLen).Float64()
//	}
//	return needCount
//}
//
//func ChangeStringNeedCount(needCount string) string {
//	splitArr := strings.Split(needCount, ".")
//	if len(splitArr) == 2 {
//		var roundLen = len(splitArr[1])
//		var subLen = int32(roundLen - 1)
//		d, _ := decimal.NewFromString(needCount)
//		needCount = d.Round(subLen).String()
//	}
//	return needCount
//}
//
//
