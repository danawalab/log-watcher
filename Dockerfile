FROM golang

WORKDIR /home/danawa/apps/logwatcher

COPY . .

RUN go build -o logwatcher cmd/logwatcher/main.go

RUN chmod +x logwatcher

CMD ["./logwatcher", "./setting.json"]