package logClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/duzhenlin/skittle/src/config"
	"github.com/duzhenlin/skittle/src/user"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/syyongx/php2go"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type LogClient struct {
	ctx         context.Context
	config      *config.Config
	redisClient *redis.Client
}

type general struct {
	RequestUrl    string `json:"request_url"`
	RequestMethod string `json:"request_method"`
	RemoteAddress string `json:"remote_address"`
}

type User struct {
	Id           string `json:"id"`
	Nickname     string `json:"nickname"`
	Organization string `json:"organization"`
	Phone        string `json:"phone"`
}
type Http struct {
	General        general `json:"general"`
	RequestHeaders string  `json:"request_headers"`
	RequestPayload string  `json:"request_payload"`
	FormData       string  `json:"form_data"`
}

type Extra struct {
	RuntimeUuid string `json:"runtime_uuid"`
	GoVersion   string `json:"go_version"`
	User        User   `json:"user"`
	Http        Http   `json:"http"`
}

type LogEntry struct {
	Message   string    `json:"message"`
	Context   string    `json:"context"`
	Level     int       `json:"level"`
	LevelName string    `json:"level_name"`
	Channel   string    `json:"channel"`
	Datetime  time.Time `json:"datetime"`
	Extra     Extra     `json:"extra"`
}

// 提取 IP 地址
func extractIPAddress(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 如果解析失败，返回原始地址
		return remoteAddr
	}
	return host
}

// NewLog 创建log服务对象
func NewLog(ctx context.Context, cfg *config.Config, rdb *redis.Client) *LogClient {
	return &LogClient{
		ctx:         ctx,
		config:      cfg,
		redisClient: rdb,
	}
}

func (l *LogClient) Log(c *http.Request) {

	elasticsearchConfig := elasticsearch.Config{
		Addresses: []string{l.config.EsConfig.Address},
		Username:  l.config.EsConfig.Username,
		Password:  l.config.EsConfig.Password,
		// 需要忽略证书
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	elasticsearchClient, _ := elasticsearch.NewClient(elasticsearchConfig)

	res, err := elasticsearchClient.Info()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	//盗用user里的GetUserInfoNew方法
	userinfo, _ := user.NewUser(l.ctx, l.config).SetRedisClient(l.redisClient).GetUserInfoNew(c)

	if userinfo.ID == "" {
		log.Printf("无效的id")
	} else {
		//首先创建用户信息
		userInfo := User{
			Id:           userinfo.ID,
			Nickname:     userinfo.Nickname,
			Organization: userinfo.Orgname,
			Phone:        userinfo.Username,
		}

		// 创建 HTTP 请求的通用信息
		generalInfo := general{
			RequestUrl:    c.URL.Path,
			RequestMethod: c.Method,
			RemoteAddress: extractIPAddress(c.RemoteAddr),
		}
		// 创建 HTTP 请求的详细信息
		httpInfo := Http{
			General: generalInfo,
		}
		requestHeaders := make(map[string]string)
		// 填充请求头
		for name, values := range c.Header {
			if len(values) > 0 {
				requestHeaders[name] = values[0]
			}
		}
		headersJSON, err := json.Marshal(httpInfo.RequestHeaders)
		if err != nil {
			log.Fatalf("Error marshaling request headers: %s", err)
		}
		httpInfo.RequestHeaders = string(headersJSON)

		// 填充请求体
		if c.Body != nil {
			var bodyBytes []byte
			bodyBytes, err = ioutil.ReadAll(c.Body)
			c.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			if err != nil {
				log.Fatalf("Error reading request body: %s", err)
			}
			httpInfo.RequestPayload = string(bodyBytes)
		} else {
			httpInfo.RequestPayload = "{}" // 如果没有请求体，设置为空对象
		}

		// 填充表单数据并过滤文件
		if c.Method == http.MethodPost || c.Method == http.MethodPut {
			if c.Form != nil {
				formData := make(map[string]interface{})
				for key, values := range c.Form {
					if len(values) > 0 {
						formData[key] = values[0]
					}
				}
				formDataJson, _ := json.Marshal(formData)
				httpInfo.FormData = string(formDataJson)
			}
		} else {
			httpInfo.FormData = "{}" // 如果不是 POST 或 PUT 请求，设置为空对象
		}
		// 创建额外信息
		runtimeUuid, err := uuid.NewV4()
		runtimeUuidString := php2go.StrReplace("-", "", runtimeUuid.String(), -1)
		extraInfo := Extra{
			RuntimeUuid: runtimeUuidString,
			GoVersion:   "go1.20", // 这里可以替换为实际的 Go 版本
			User:        userInfo,
			Http:        httpInfo,
		}
		// 创建日志记录
		logEntry := LogEntry{
			Datetime:  time.Now(),
			Level:     200,
			LevelName: "INFO",
			Message:   "",
			Context:   c.URL.Path,
			Channel:   l.config.EsConfig.Channel,
			Extra:     extraInfo,
		}

		// 将日志记录转换为 JSON
		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(logEntry); err != nil {
			log.Fatalf("Error encoding logClient entry: %s", err)
		}

		// 将日志记录发送到 Elasticsearch
		addRes, err := elasticsearchClient.Index(
			l.config.EsConfig.IndexName, // 索引名称
			&buf,                        // 日志数据
			elasticsearchClient.Index.WithContext(context.Background()),
		)
		if err != nil {
			log.Fatalf("Error indexing document: %s", err)
		}
		defer addRes.Body.Close()

		// 检查响应
		if addRes.IsError() {
			log.Printf("Error indexing document: %s", addRes.String())
		} else {
			fmt.Println("Log entry indexed successfully")
		}
	}

}
