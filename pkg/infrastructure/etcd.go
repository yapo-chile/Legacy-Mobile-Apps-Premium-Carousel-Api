package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tidwall/gjson"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/repository"
)

// EtcdLogger defines log functions needed
type EtcdLogger interface {
	Debug(format string, params ...interface{})
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

// etcd works remote configuration service loader
type etcd struct {
	logger  EtcdLogger
	content *EtcdContent
}

// EtcdContent holds response from Etcd Service
type EtcdContent struct {
	Action string   `json:"action"`
	Node   EtcdNode `json:"node"`
}

// EtcdNode holds content from Etcd service
type EtcdNode struct {
	Key           string     `json:"key"`
	Value         string     `json:"value"`
	ModifiedIndex int        `json:"modifiedIndex"`
	CreatedIndex  int        `json:"createdIndex"`
	Nodes         []EtcdNode `json:"nodes"`
	IsDir         bool       `json:"dir"`
}

func readAll(resp *http.Response, logger EtcdLogger) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				logger.Error("ReadError %s", readErr.Error())
			} else {
				logger.Error("Response: %s", string(body))
			}
		}
		return nil, fmt.Errorf("response code invalid: %s", resp.Status)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logger.Error("Error on read body: %s", readErr.Error())
		return nil, readErr
	}
	return body, nil

}

// NewEtcd remote configuration loader
// This method get the configuration file from remote host.
// On success returns a Etcd struct, if not success returns an Error and Etcd
// with an EtcdContent empty
func NewEtcd(host, path, prefix string, logger EtcdLogger) (repository.Config, error) {
	etcd := &etcd{logger: logger, content: &EtcdContent{}}
	var netClient = &http.Client{
		Timeout: time.Second * 100,
	}
	url := fmt.Sprintf("%s%s%s", host, prefix, path)
	resp, err := netClient.Get(url)
	if err != nil {
		logger.Error("Error to get url %s", url)
		return etcd, err
	}
	defer resp.Body.Close() // nolint: errcheck
	body, err := readAll(resp, logger)
	if err != nil {
		return etcd, err
	}
	if jsonErr := json.Unmarshal(body, &etcd.content); jsonErr != nil {
		logger.Error("Error %s", jsonErr.Error())
		return etcd, jsonErr
	}
	logger.Info("Conf %s loaded", url)
	return etcd, nil
}

// Get gets the result of a GET method using given key
func (v etcd) Get(key string) string {
	if v.content == nil || &v.content.Node == nil {
		v.logger.Error("Empty conf")
		return ""
	}
	if v.content.Node.IsDir {
		v.logger.Error("Conf %s is a dir, i can not get some value", v.content.Node.Key)
		return ""
	}
	return gjson.Get(v.content.Node.Value, key).String()
}

// Get gets the result of a GET method using given key
func (v etcd) GetMap(key string) map[string]string {
	if v.content == nil || &v.content.Node == nil {
		v.logger.Error("Empty conf")
		return map[string]string{}
	}
	if v.content.Node.IsDir {
		v.logger.Error("Conf %s is a dir, i can not get some value", v.content.Node.Key)
		return map[string]string{}
	}
	result := make(map[string]string)
	mapRes := gjson.Get(v.content.Node.Value, key).Map()
	for k, v := range mapRes {
		result[k] = v.String()
	}
	return result
}
