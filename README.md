## 소개

프로그램에서 생성되는 로그파일을 감시하고 알람 기능을 제공합니다.

로그파일을 감시하고 있다가 설정된 레벨의 로그가 찍히면 감지하여 Telegram으로 보내주는 기능을 수행합니다.

## 설정파일

setting.json SAMPLE )
```text

"watch": [
  {
    "label": "ES NODE-1"   // 텔레그램에 "[ES NODE-1] {로그내용}" 이렇게 표시됨.
    "path": "/data/es-7.8.1/log",
    "file": "es-cluster.log",
    "level" : ["ERROR", "WARN"]
  }
  {
    "label": "ES NODE-2"
    "path": "/data/es-7.8.1/log",
    "file": "es-cluster.log",
    "level" : ["ERROR", "WARN"],
    "contains" : ["overhead", "outofmemory"] // 로그레벨과 상관없이 문자열이 포함되어 있다면 전송해준다.
  }
],
 
"interval": "1m",
 
"telegram": {
   "bot_token": "999783271:AAH-RWxxxxxxxxxxxxxrYLXgYM-uJrY" // 봇 토큰
   "chat_id": 1111111111 // chat_id(channel)
},
 
"logging": { // 로그와쳐 자체 로그
    "filename": "/logs/application.log",
    "level": "info",
    "max_size": 500,
    "max_backups": 3,
    "max_age": 14,
    "compress": true
}
```