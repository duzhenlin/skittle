// Package log
// Created by Duzhenlin
// @Author   Duzhenlin
// @Email: duzhenlin@vip.qq.com
// @Date: 2025/7/8
// @Time: 15:49

package log

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/duzhenlin/skittle/v2/src/config"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net/http"
)

type ESLogger struct {
	client *elasticsearch.Client
	index  string
}

func NewESLogger(config *config.Config) (*ESLogger, error) {

	cfg := elasticsearch.Config{
		Addresses: []string{config.EsConfig.Address},
		Username:  config.EsConfig.Username,
		Password:  config.EsConfig.Password,
		// 需要忽略证书
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ESLogger{client: client, index: config.EsConfig.IndexName}, nil
}

func (l *ESLogger) WriteLog(doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:      l.index,
		DocumentID: "", // 自动生成ID
		Body:       bytes.NewReader(data),
		Refresh:    "true", // 可选，开发调试时可用
	}
	res, err := req.Do(context.Background(), l.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es index error: %s", res.String())
	}
	return nil
}
