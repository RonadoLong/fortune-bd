package utils

// ReTryFunc 重试
// count 次数  doFunc 业务方法
// doFunc 放回值，isStop 是否停止， error 错误信息
func ReTryFunc(count int, doFunc func() (bool, error)) error {
	var err error
	var isStop bool
	for {
		if count == 0 {
			break
		}
		if isStop, err = doFunc(); err != nil && !isStop {
			count--
			continue
		}
		break
	}

	return err
}
