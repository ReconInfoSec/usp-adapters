package usp_itglue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/refractionPOINT/go-uspclient"
	"github.com/refractionPOINT/usp-adapters/utils"
)

const (
	logURL = "/logs"
)

var URL = map[string]string{
	"enterprise": "https://api.itglue.com",
}

type ITGlueConfig struct {
	ClientOptions uspclient.ClientOptions `json:"client_options" yaml:"client_options"`
	ApiKey        string                  `json:"apikey" yaml:"apikey"`
	Endpoint      string                  `json:"endpoint" yaml:"endpoint"`
}

func (c *ITGlueConfig) Validate() error {
	if err := c.ClientOptions.Validate(); err != nil {
		return fmt.Errorf("client_options: %v", err)
	}
	if c.ApiKey == "" {
		return errors.New("missing API key")
	}
	if c.Endpoint == "" {
		return errors.New("missing endpoint")
	}
	_, ok := URL[c.Endpoint]
	if !strings.HasPrefix(c.Endpoint, "https://") && !ok {
		return fmt.Errorf("invalid endpoint, not https or in %v", URL)
	}
	return nil
}

func NewITGlueAdapter(conf ITGlueConfig) (*ITGlueAdapter, chan struct{}, error) {
	var err error
	a := &ITGlueAdapter{
		config: conf,
		ctx:    context.Background(),
		doStop: utils.NewEvent(),
	}

	if strings.HasPrefix(conf.Endpoint, "https://") {
		a.endpoint = conf.Endpoint
	} else if v, ok := URL[conf.Endpoint]; ok {
		a.endpoint = v
	} else {
		return nil, nil, fmt.Errorf("not a valid api endpoint: %s", conf.Endpoint)
	}

	a.uspClient, err = uspclient.NewClient(conf.ClientOptions)
	if err != nil {
		return nil, nil, err
	}

	a.httpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
		},
	}

	a.chStopped = make(chan struct{})

	a.wgSenders.Add(3)
	go a.GetLogs()

	go func() {
		a.wgSenders.Wait()
		close(a.chStopped)
	}()

	return a, a.chStopped, nil
}

type ITGlueAdapter struct {
	config     ITGlueConfig
	lastLogID  string
	lastCalled bool

	uspClient  *uspclient.Client
	httpClient *http.Client

	endpoint string

	chStopped chan struct{}
	wgSenders sync.WaitGroup
	doStop    *utils.Event

	ctx context.Context
}

type ITGlueResponse struct {
	Data []ITGlueRecord `json:"data"`
}

type ITGlueRecord struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
}

func (adapter *ITGlueAdapter) HandleRequest(w http.ResponseWriter, r *http.Request) {
	response, err := adapter.GetLogs()
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// If this is not the first call and we have a last log ID,
	// filter the response to only include logs with an ID greater
	// than the last log ID.
	if adapter.lastCalled && adapter.lastLogID != "" {
		filteredResponse := ITGlueResponse{}
		for _, record := range response.Data {
			if record.ID > adapter.lastLogID {
				filteredResponse.Data = append(filteredResponse.Data, record)
			}
		}
		response = filteredResponse
	}

	// Update the last log ID and set lastCalled to true
	if len(response.Data) > 0 {
		lastLog := response.Data[len(response.Data)-1]
		adapter.lastLogID = lastLog.ID
		adapter.lastCalled = true
	}

	json.NewEncoder(w).Encode(response)
}

func (adapter ITGlueAdapter) GetLogs() (ITGlueResponse, error) {
	var response ITGlueResponse

	url := adapter.config.Endpoint + logURL

	// If this is not the first call and we have a last log ID,
	// add a filter parameter to only retrieve logs with an ID greater
	// than the last log ID.
	if adapter.lastCalled && adapter.lastLogID != "" {
		url += "?filter[id]=gt:" + adapter.lastLogID
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", adapter.config.ApiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
