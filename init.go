package htx

import "time"

var ApiConfig = ApiConfigModel{}

// var ApiConfig ApiConfigModel

const (
	ExchangeName = "htx"
)

var (
	ProxyUrl = ""
	UseProxy = false
)

func InitConfig(config ApiConfigModel) {
	ApiConfig = config
}

func SetProxy(proxyUrl string) {
	UseProxy = true
	ProxyUrl = proxyUrl
}

func ClearProxy() {
	UseProxy = false
	ProxyUrl = ""
}

func GetTimeFloat() float64 {
	return float64(time.Now().UnixNano()) / 1000000000
}
