package wrapper

import (
	"AnimeLifeBackEnd/global"
	"net/http"
	"time"
)

func Get(url string) (resp *http.Response, err error) {
	global.Logger.Infof("GET %s", url)
	// log the request time
	start := time.Now()
	resp, err = http.Get(url)
	if err != nil {
		global.Logger.Errorf("GET %s failed, err: %v", url, err)
		// TODO add retry logic
	}
	global.Logger.Infof("GET %s finished, cost: %v", url, time.Since(start))
	return resp, err
}

func GetWithHeader(url string, header map[string]string) (resp *http.Response, err error) {
	global.Logger.Infof("GET %s", url)
	// log the request time
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		global.Logger.Errorf("GET %s failed, err: %v", url, err)
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		global.Logger.Errorf("GET %s failed, err: %v", url, err)
		// TODO add retry logic
	}
	global.Logger.Infof("GET %s finished, cost: %v", url, time.Since(start))
	return resp, err
}
