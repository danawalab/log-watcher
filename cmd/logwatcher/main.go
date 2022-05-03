package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.danawa.com/fastcatx/log-scrap/logging"
	"gitlab.danawa.com/fastcatx/log-scrap/service"
	"gitlab.danawa.com/fastcatx/log-scrap/setting"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

// 수집하기위한 상태.
var Status = "Ready"

func setupRouter() http.Handler {
	e := gin.New()
	e.GET("/", func(c *gin.Context) {
		c.JSON(200, "LogWatcher Running")
	})
	return e
}

func main() {
	// 설정파일을 로드합니다.
	fmt.Println("SETTING_FILE_PATH : " + os.Args[1])
	loadSettingError := setting.LoadSetting(os.Args[1])
	if loadSettingError != nil {
		fmt.Println("[ERROR] ", loadSettingError)
		os.Exit(1)
	}
	// 로그 설정 합니다.
	logging.LoadLogging(setting.GetLogging())
	logging.Info(fmt.Sprint("INIT_PORT : ", setting.GetPort()))
	logging.Info(fmt.Sprintf("WATCH_TARGET : %d", len(setting.GetWatches())))
	for index, watch := range setting.GetWatches() {

		fmt.Println("WATCH_FILE " + fmt.Sprintf("%d", index+1) + " : " + watch.File + ", LABEL : " + watch.Label + ", PATH : " + watch.Path)

		// 로그 수집 함수를 호출합니다.
		go service.ScrapRun(watch)
	}

	// 텔레그램 메시징 함수를 호출합니다.
	go service.IntervalSend(&Status, setting.GetInterval(), setting.GetTelegramInfo())

	// 시작할때 카운트 다운을 합니다.
	// 용도. 로그 수집시 마지막 라인 이동 시 지연이 혹시 있을 수 있어 넉넉하게 5초 기다려줍니다.
	for c := 5; c > 0; c-- {
		time.Sleep(1 * time.Second)
		logging.Info(fmt.Sprintf("Start Count Down %d", c))
	}
	Status = "Running"

	// 로그 와쳐 서버를 생성합니다.
	apiServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.GetPort()),
		Handler:        setupRouter(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	gin.SetMode(gin.ReleaseMode)
	g := errgroup.Group{}
	g.Go(func() error {
		return apiServer.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
