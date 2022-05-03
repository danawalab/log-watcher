package setting

import (
	"encoding/json"
	"io/ioutil"

	"gitlab.danawa.com/fastcatx/log-scrap/model"
)

var (
	setting = model.Setting{}
)

func LoadSetting(filePath string) error {
	raw, err := ioutil.ReadFile(filePath)
	if err == nil {
		_ = json.Unmarshal(raw, &setting)
	}
	return err
}

func GetWatches() []model.Watch {
	return setting.Watch
}

func GetPort() int {
	return setting.Port
}

func GetLogging() model.Logging {
	return setting.Logging
}

func GetInterval() string {
	return setting.Interval
}

func GetTelegramInfo() model.Telegram {
	return setting.Telegram
}
