package htx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/laoliu6668/esharp_htx_utils/util"
)

func (c *ApiConfigModel) Get(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("GET", gateway, path, data, time.Second*5)
}

func (c *ApiConfigModel) Post(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("POST", gateway, path, data, time.Second*3)
}

func (c *ApiConfigModel) GetTimeout(gateway, path string, data map[string]any, timeout time.Duration) (body []byte, resp *http.Response, err error) {
	return c.Request("GET", gateway, path, data, timeout)
}
func (c *ApiConfigModel) PostTimeout(gateway, path string, data map[string]any, timeout time.Duration) (body []byte, resp *http.Response, err error) {
	return c.Request("POST", gateway, path, data, timeout)
}

// 获取TRONSCAN API数据
func (c *ApiConfigModel) Request(method, gateway, path string, data map[string]any, timeout time.Duration) (body []byte, resp *http.Response, err error) {

	if timeout == 0 {
		timeout = time.Second * 5
	}
	// 创建http client
	client := &http.Client{
		Timeout: timeout,
	}
	if UseProxy {
		uri, _ := url.Parse(fmt.Sprintf("http://%s", ProxyUrl))
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	// 声明 querypMap
	queryMap := map[string]any{
		"AccessKeyId":      c.AccessKey,
		"Timestamp":        UTCTimeNow(),
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
	}
	// 声明 body
	var reqBody io.Reader
	if method == "POST" {
		// 添加body
		buf, _ := json.Marshal(data)
		reqBody = strings.NewReader(string(buf))
	} else if method == "GET" {
		// 融合query 参数
		for k, v := range data {
			queryMap[k] = v
		}
	} else {
		err = errors.New("不支持的http方法")
		return
	}

	// 签名
	queryMap["Signature"] = Signature(method, gateway, path, queryMap, c.SecretKey)
	// 构造query
	url := GetQueryUrl("https://", gateway, path, queryMap)

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("http响应%v", resp.StatusCode)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func GetQueryUrl(proto, gateway, path string, queryMap map[string]any) string {
	return fmt.Sprintf("%s%s%s?%s", proto, gateway, path, util.HttpBuildQuery(queryMap))

}

func Signature(method, gateway, path string, args map[string]any, key string) string {
	str1 := strings.ToUpper(method) + "\n"
	str2 := gateway + "\n"
	str3 := path + "\n"
	str4 := util.HttpBuildQuery(args)
	str5 := str1 + str2 + str3 + str4
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str5))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func UTCTimeNow() string {
	return time.Now().In(time.UTC).Format("2006-01-02T15:04:05")
}
