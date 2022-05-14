package elastic

/*
@author: AxisZql
@desc: elasticsearch util
@date: 2022-5-6 11:09 PM
*/

import (
	"blog-server/common"
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

var (
	es7  *elastic.Client
	once sync.Once
)

func newElasticClient() {
	var err error
	es7, err = elastic.NewClient(elastic.Config{
		Addresses: []string{common.Conf.Es.Addr},
		Transport: &http.Transport{ //配置http连接池
			MaxIdleConns:          10,          //最大keep-alive连接数量
			ResponseHeaderTimeout: time.Second, // 设置响应超时时间
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	res, err := es7.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)
	//check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body:%s", err)
	}
	// Print Client ans server version numbers
	log.Printf("Client: %s", elastic.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
}

// GetElasticClient 单例模式初始化ES客户端
func GetElasticClient() *elastic.Client {
	once.Do(func() {
		newElasticClient()
	})
	return es7
}

// Query
// @index: 索引名称
// @query: 查询参数
// @dest：查询结果存放(必须是指针)
func Query(index string, query interface{}, dest interface{}) (err error) {
	p := reflect.ValueOf(dest)
	if p.Kind() != reflect.Ptr {
		return errors.New("dest must be pointer")
	}
	es := GetElasticClient()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if err := json.NewDecoder(res.Body).Decode(dest); err != nil {
		return err
	}
	return nil
}
