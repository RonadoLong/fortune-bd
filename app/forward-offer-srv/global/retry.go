package global

import "time"

var timeOut = time.Millisecond * 100

// RunRetry 重试多次方法
func RunRetry(count int, todo func() error) error {
	count--
	err := todo()
	if err != nil {
		<-time.After(timeOut)
		if count == 0 {
			return err
		}
		err = RunRetry(count, todo)
	}
	return err
}
