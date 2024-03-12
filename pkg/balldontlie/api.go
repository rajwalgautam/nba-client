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
var pageSize = 100

var reqCounter = 0
var rcmu = &sync.Mutex{}

type Client interface {
	Ping() error
	GamesOnDate(d string) ([]Game, error)
	GamesDateRange(start, end string) ([]Game, error)
	PublishStatsByGameId(id int, queue chan []Stats) ([]Stats, error)
	StatsByDateRange(start, end string) ([]Stats, error)
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

	cursor := 0
	allGames := make([]Game, 0)
	for {
		time.Sleep(time.Second * 2)
		url := fmt.Sprintf("%s&cursor=%d", u, cursor)
		b, err := api.sendReq(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("request: %s", err)
		}
		sw := gameWrapper{}
		err = json.Unmarshal(b, &sw)
		if err != nil {
			return nil, fmt.Errorf("json: %s", err)
		}
		if len(sw.Data) > 0 {
			fmt.Printf("found games from %s to %s\n", sw.Data[0].Date, sw.Data[len(sw.Data)-1].Date)
			allGames = append(allGames, sw.Data...)
		}
		if len(sw.Data) < pageSize {
			return allGames, nil
		}
		md, err := getMetadata(b)
		if err != nil {
			return nil, fmt.Errorf("metadata err: %s", err)
		}
		cursor = md.NextCursor
	}
}

func getMetadata(b []byte) (Metadata, error) {
	md := metadataWrapper{}
	err := json.Unmarshal(b, &md)
	return md.Metadata, err
}

func (api *bdlClient) StatsByDateRange(start, end string) ([]Stats, error) {
	u, _ := url.JoinPath(api.base.String(), "stats")
	u = fmt.Sprintf("%s?start_date=%s&end_date=%s&per_page=100", u, start, end)

	cursor := 0
	allStats := make([]Stats, 0)
	for {
		time.Sleep(time.Second * 2)
		url := fmt.Sprintf("%s&cursor=%d", u, cursor)
		b, err := api.sendReq(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("request: %s", err)
		}
		sw := statsWrapper{}
		err = json.Unmarshal(b, &sw)
		if err != nil {
			return nil, fmt.Errorf("json: %s", err)
		}
		if len(sw.Data) > 0 {
			allStats = append(allStats, sw.Data...)
		}
		if len(sw.Data) < pageSize {
			return allStats, nil
		}
		cursor += pageSize
	}
}

func (api *bdlClient) PublishStatsByGameId(id int, queue chan []Stats) ([]Stats, error) {
	u, _ := url.JoinPath(api.base.String(), "stats")
	u = fmt.Sprintf("%s?per_page=100&game_ids[]=%d", u, id)

	cursor := 0
	allStats := make([]Stats, 0)
	for {
		time.Sleep(time.Second * 2)
		url := fmt.Sprintf("%s&cursor=%d", u, cursor)
		b, err := api.sendReq(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("request: %s", err)
		}
		sw := statsWrapper{}
		err = json.Unmarshal(b, &sw)
		if err != nil {
			return nil, fmt.Errorf("json: %s", err)
		}
		if len(sw.Data) > 0 {
			fmt.Printf("published stats for %s\n", sw.Data[0].Game.Date)
			queue <- sw.Data
		}
		if len(sw.Data) < pageSize {
			return allStats, nil
		}
		cursor += pageSize
	}
}
