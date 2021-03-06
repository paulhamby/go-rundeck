package rundeck

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (rc *RundeckClient) Get(i interface{}, path string, options map[string]string) error {
	return rc.makeRequest(i, "GET", path, options)
}

func (client *RundeckClient) RawGet(path string, qp map[string]string) string {
	qs := url.Values{}
	for k, v := range qp {
		qs.Add(k, v)
	}
	base_req_path := client.Config.BaseURL
	u, err := url.Parse(base_req_path + "/api/12/" + path)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	u.RawQuery = qs.Encode()
	request, _ := http.NewRequest("GET", u.String(), nil)
	request.Header.Add("X-Rundeck-Auth-Token", client.Config.Token)
	request.Header.Add("Accept", "application/xml")
	request.Header.Add("user-agent", "rundeck-go.v12")
	r, err := client.Client.Do(request)
	if err != nil {
		return err.Error()
	} else {
		defer r.Body.Close()
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err.Error()
		} else {
			return string(contents[:])
		}
	}

}

func (client *RundeckClient) makeRequest(i interface{}, method string, path string, params map[string]string) error {
	qs := url.Values{}
	for q, p := range params {
		qs.Add(q, p)
	}
	base_req_path := client.Config.BaseURL + "/api/12/" + path
	u, err := url.Parse(base_req_path)
	if err != nil {
		return err
	}
	if len(params) != 0 {
		u.RawQuery = qs.Encode()
	}
	request, _ := http.NewRequest(method, u.String(), nil)
	request.Header.Add("X-Rundeck-Auth-Token", client.Config.Token)
	request.Header.Add("Accept", "application/xml")
	request.Header.Add("user-agent", "rundeck-go.v12")
	r, err := client.Client.Do(request)
	if err != nil {
		return err
	} else {
		defer r.Body.Close()
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if r.StatusCode == 404 {
			return errors.New(r.Status)
		}

		if r.StatusCode != 200 {
			var data RundeckError
			xml.Unmarshal(contents, &data)
			return errors.New("non-200 response: " + data.Message)
		} else {
			xml.Unmarshal(contents, &i)
			return nil
		}
	}
}
