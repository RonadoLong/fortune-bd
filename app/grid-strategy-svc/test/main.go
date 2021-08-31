package main

import (
	"flag"
	"fmt"
	"fortune-bd/app/grid-strategy-svc/model"
	"io/ioutil"
	"time"

	"github.com/json-iterator/go"
	"github.com/zhufuyi/pkg/mongo"
)

// 配置初始化
var (
	saveType   string // 数据保存目标，csv或mgo
	mgoAddr    string // mongodb地址
	csvFile    string // csv文件绝对路径
	configFile string // 配置文件路径
)

func main() {
	flag.StringVar(&saveType, "t", "mgo", "csv or mgo")
	flag.StringVar(&csvFile, "csv", "", "csv absolute path")
	flag.StringVar(&mgoAddr, "mgo", "mongodb://wq:abc123@192.168.5.5:38888/ifortune", "mongodb addr")
	flag.StringVar(&configFile, "c", "./config.json", "config param")
	flag.Parse()

	if saveType != "csv" && saveType != "mgo" {
		fmt.Printf("%s is not supported\n", saveType)
		return
	}

	if saveType == "mgo" && mgoAddr == "" {
		fmt.Printf("mgoAddr is empty\n")
		return
	}
	if configFile == "" {
		fmt.Printf("configFile is empty\n")
		return
	}

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	gf := &model.CalculateGrid{}
	err = jsoniter.Unmarshal(content, gf)
	if err != nil {
		panic(err)
	}

	data := make(chan *model.GridFilter)

	var file string
	if saveType == "csv" {
		if csvFile != "" {
			file = csvFile
		} else {
			file = fmt.Sprintf("%s-%s-%d.csv", gf.Exchange, gf.Symbol, time.Now().Unix())
		}
		go model.Save2CSV(file, data)
	} else if saveType == "mgo" {
		err := mongo.InitializeMongodb(mgoAddr)
		if err != nil {
			panic(err)
		}
		go model.Save2Mgo(data)
	}

	gf.DoneAndSave(data)
	close(data)
	time.Sleep(2 * time.Second)
}
