package mongoClient

import (
	"testing"
)

func TestInitMongo(t *testing.T) {
	_, err := InitMongo("mongodb://wquant:wqabc123@47.57.161.25:38028/ifortune")
	if err != nil {
		t.Errorf("mongodb连接失败")
	}

}
