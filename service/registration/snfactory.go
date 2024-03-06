package registration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"time"
)

type SNFactory interface {
	GetSN(key string) string
}

type SNItem struct {
	Key  string `json:"key"`
	Date string `json:"date"`
	SN   int    `json:"sn"`
}

type DefaultSNFactory struct {
	SNPool map[string]SNItem
}

func (factory *DefaultSNFactory) Init(repo Repository) {
	//查询所有的号源配置信息
	snItemList, _ := repo.GetCurrentSN()
	if snItemList == nil {
		return
	}

	//初始化SNPool
	factory.SNPool = make(map[string]SNItem)
	for _, item := range *snItemList {
		factory.SNPool[item.Key] = item
	}
}

func (factory *DefaultSNFactory) SavePool(fileName string) {
	//convert SNPool to json and save to file
	jsonStr, _ := json.Marshal(factory.SNPool)
	ioutil.WriteFile(fileName, jsonStr, 0644)
}

func (factory *DefaultSNFactory) LoadPool(fileName string) {
	//load json from file and convert to SNPool
	jsonStr, err := ioutil.ReadFile(fileName)
	if err != nil {
		slog.Error("read file failed", "file", fileName, "error", err.Error())
		return
	}

	err = json.Unmarshal(jsonStr, &factory.SNPool)
	if err != nil {
		slog.Error("json file decode failed.", "error", err.Error())
	}
}

func (factory *DefaultSNFactory) GetSN(key string) string {
	nowDate := time.Now().Format("2006-01-02")
	item, ok := factory.SNPool[key]
	if !ok {
		item = SNItem{
			Key:  key,
			Date: nowDate,
			SN:   0,
		}
	} else {
		if item.Date != nowDate {
			item.Date = nowDate
			item.SN = 0
		}
	}

	item.SN++
	factory.SNPool[key] = item

	return fmt.Sprintf("%04d", item.SN)
}
