package balldontlie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var baseURL = "https://api.balldontlie.io/v1"

var reqCounter = 0
var rcmu = &sync.Mutex{}

type Client interface {
	Ping() error
	GamesOnDate(d string) ([]Game, error)
	GamesDateRange(start, end string) ([]Game, error)
	StatsByGameId(id int) ([]Stats, error)
}

type bdlClient struct {
	httpc  *http.Client
	apikey string
	base   *url.URL
}

func New(k string) Client {
	b, _ := url.Parse(baseURL)
	return &bdlClient{
		httpc:  &http.Client{Timeout: time.Second * 10},
		apikey: k,
		base:   b,
	}
}

func (api *bdlClient) Ping() error {
	u, _ := url.JoinPath(api.base.String(), "teams/1")
	req, err := api.newReq(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	resp, err := api.httpc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}
	return nil
}

func (api *bdlClient) newReq(method string, url string, body []byte) (*http.Request, error) {
	rcmu.Lock()
	reqCounter++
	rcmu.Unlock()
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", api.apikey)
	return req, err
}

func (api *bdlClient) GamesOnDate(d string) ([]Game, error) {
	u, _ := url.JoinPath(api.base.String(), "games")
	u = fmt.Sprintf("%s?dates[]=%s&per_page=100", u, d)
	b, err := api.sendReq(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	gw := gameWrapper{}
	err = json.Unmarshal(b, &gw)
	if err != nil {
		return nil, err
	}
	if len(gw.Data) > 0 {
		return gw.Data, nil
	}

	return []Game{}, nil
}

func (api *bdlClient) sendReq(method string, url string, body []byte) ([]byte, error) {
	r, err := api.newReq(method, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := api.httpc.Do(r)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (api *bdlClient) GamesDateRange(start, end string) ([]Game, error) {
	u, _ := url.JoinPath(api.base.String(), "games")
	u = fmt.Sprintf("%s?start_date=%s&end_date=%s&per_page=100", u, start, end)
	b, err := api.sendReq(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	gw := gameWrapper{}
	err = json.Unmarshal(b, &gw)
	if err != nil {
		return nil, err
	}
	if len(gw.Data) > 0 {
		return gw.Data, nil
	}

	return nil, nil
}

func (api *bdlClient) StatsByGameId(id int) ([]Stats, error) {
	u, _ := url.JoinPath(api.base.String(), "stats")
	u = fmt.Sprintf("%s?game_ids[]=%d", u, id)
	b, err := api.sendReq(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %s", err)
	}
	sw := statsWrapper{}
	err = json.Unmarshal(b, &sw)
	if err != nil {
		fmt.Println(sw)
		return nil, fmt.Errorf("json: %s", err)
	}
	if len(sw.Data) > 0 {
		return sw.Data, nil
	}
	return nil, nil
}
