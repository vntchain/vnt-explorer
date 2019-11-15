package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"github.com/astaxie/beego/httplib"
)

type Param struct {
	Key   string
	Value string
}

func CallApi(requestUrl string, params []Param) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: httplib.TimeoutDialer(5 * time.Second, 5 * time.Second),
			DisableKeepAlives:true,
			MaxIdleConns:100,
			MaxConnsPerHost:100,
		},
	}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	for _, param := range params {
		q.Add(param.Key, param.Value)
	}

	req.Header.Add("User-Agent", "xxx")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("Failed to get resp body")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error Response : %s", resp.Status)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read resp.Body err: %s", err)
	}
	return respBody, nil
}
