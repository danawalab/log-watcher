package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"gitlab.danawa.com/fastcatx/log-scrap/logging"
	"gitlab.danawa.com/fastcatx/log-scrap/model"
)

// 이원화된 메시지 저장용 변수
var registry_A = make(map[string]model.Analytic)
var registry_B = make(map[string]model.Analytic)

// 로그 파일을 수집하는 함수입니다.
func ScrapRun(l model.Watch) {
	isFirst := true
	seekInfo := tail.SeekInfo{}
	var t *tail.Tail
	for {
		// 파일이 새로 열릴때 마다 0
		n := 0
		// 파일 경로와 파일 이름을 연결하여 파일 전체 경로를 만듭니다.
		fullPath := fmt.Sprintf("%s/%s", l.Path, l.File)
		logging.Info(fmt.Sprintf("수집 로그 파일 : %s", fullPath))

		// 파일의 마지막 라인을 찾습니다.
		location := tail.SeekInfo{}
		if isFirst {
			if seekInfo.Offset > 0 {
				isFirst = false
				location.Offset = seekInfo.Offset
			} else if info, err := os.Stat(fullPath); err == nil {
				isFirst = false
				location.Offset = info.Size()
			} else {
				// 마지막 없으면 0으로 처음부터 수집합니다.
				location.Offset = 0
			}
		} else {
			time.Sleep(5 * time.Second)
			location.Offset = 0
		}
		seekInfo.Offset = 0

		// 파일에 테일을 추가합니다.
		var err error
		t, err = tail.TailFile(fullPath, tail.Config{Follow: true, ReOpen: true, MustExist: true, Poll: true, Location: &location})

		if err != nil {
			// 파일이 없을 경우 재시도 해봅니다.
			logging.Error(fmt.Sprintf("TailFile %s", err.Error()))
			time.Sleep(5 * time.Second)
			continue
		}

		// TAIL에서 라인이 추가되면 반복이 시작됩니다.
		for line := range t.Lines {
			if line.Err != nil {
				// 라인에 에러가 있으면 로그찍고 다시 기다립니다.
				logging.Error(fmt.Sprintf("Line error message : %s", line.Err.Error()))
				continue
			}

			// 설정된 로그 레벨과 키워드를 추출
			if isTargetLogContents(line.Text, l.Level, l.Contains) {
				var analytic = model.Analytic{}
				analytic.Message = line.Text
				analytic.Label = l.Label
				analytic.Count = 1

				_, isExist := registry_A[l.Label]

				if isExist {
					// 라벨이 존재할 경우 카운트만 증가
					adjust_analytic := registry_A[l.Label]
					adjust_analytic.Count++
					registry_A[l.Label] = adjust_analytic
				} else {
					// 라벨이 없을 경우 맵에 추가
					registry_A[l.Label] = analytic
				}

				n += 1
				if n%100000 == 0 {
					n = 0
				}
				logging.Debug(fmt.Sprintf("%s, fullPath: %s", line.Text, fullPath))
			}
		}
		t.Cleanup()
		t = nil
	}
}

func isTargetLogContents(message string, levels []string, keywords []string) bool {
	isTarget := false

	// 지정한 로그 레벨 검출
	for _, level := range levels {
		if strings.Contains(message, level) {
			isTarget = true
		}
	}

	// 지정한 로그 키워드 검출
	for _, keyword := range keywords {
		if strings.Contains(message, keyword) {
			isTarget = true
		}
	}

	return isTarget
}

/**
* 설정 시간이 지나면 쌓인 버퍼의 메시지를 텔레그램으로 보내줍니다.
 */
func IntervalSend(status *string, interval string, telegram model.Telegram) {
	parseT, _ := time.ParseDuration(interval)
	logging.Info(fmt.Sprintf("Time Duration : %s", parseT))

	// 지정된 스케줄마다 텔레그램 호출
	for {
		if *status == "Running" {
			time.Sleep(parseT)

			// 쌓인 내용을 B로 이동
			registry_B = registry_A

			// 이동하고 비워준다
			registry_A = make(map[string]model.Analytic)

			for _, item := range registry_B {
				// 전송할 메시지
				message := ""

				// 텔레그램 전송
				if item.Count > 1 {
					message = "<" + item.Label + " - 총 " + strconv.Itoa(item.Count) + "건의 로그 감지>\n"
				} else if item.Count == 1 {
					message = "<" + item.Label + ">\n"
				}

				message += item.Message
				telegramRequest(message, telegram)
			}
		}
	}
}

func telegramRequest(message string, telegram model.Telegram) {
	telegramApiUrl := "https://api.telegram.org/bot" + telegram.BotToken + "/sendMessage"
	response, err := http.PostForm(
		telegramApiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(telegram.ChatId)},
			"text":    {message},
		})
	if err != nil {
		logging.Error(fmt.Sprintf("Telegram Send Error : %s", err))
	}

	defer response.Body.Close()

	// 결과 출력
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Sprintf("Telegram Response Error : %s", err))
	}

	fmt.Printf("%s\n", string(data))
}
