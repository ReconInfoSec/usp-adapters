package usp_itglue

import (
	"bytes"
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
	"github.com/refractionPOINT/go-uspclient/protocol"
	"github.com/refractionPOINT/usp-adapters/utils"
)

const (
	logsURL = "/logs"
)

type opRequest struct {
	Limit     int    `json:"limit,omitempty"`
	StartTime string `json:"created-at,omitempty"`
}

var URL = map[string]string{
	"enterprise": "https://api.itglue.com",
}

type ITGlueAdapter struct {
	conf       ITGlueConfig
	uspClient  *uspclient.Client
	httpClient *http.Client

	endpoint string

	chStopped chan struct{}
	wgSenders sync.WaitGroup
	doStop    *utils.Event

	ctx context.Context
}

type ITGlueConfig struct {
	ClientOptions uspclient.ClientOptions `json:"client_options" yaml:"client_options"`
	Token         string                  `json:"token" yaml:"token"`
	Endpoint      string                  `json:"endpoint" yaml:"endpoint"`
}

func (c *ITGlueConfig) Validate() error {
	if err := c.ClientOptions.Validate(); err != nil {
		return fmt.Errorf("client_options: %v", err)
	}
	if c.Token == "" {
		return errors.New("missing token")
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
		conf:   conf,
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
	go a.fetchEvents(logsURL)

	go func() {
		a.wgSenders.Wait()
		close(a.chStopped)
	}()

	return a, a.chStopped, nil
}

func (a *ITGlueAdapter) Close() error {
	a.conf.ClientOptions.DebugLog("closing")
	a.doStop.Set()
	a.wgSenders.Wait()
	err1 := a.uspClient.Drain(1 * time.Minute)
	_, err2 := a.uspClient.Close()
	a.httpClient.CloseIdleConnections()

	if err1 != nil {
		return err1
	}

	return err2
}

func (a *ITGlueAdapter) fetchEvents(url string) {
	defer a.wgSenders.Done()
	defer a.conf.ClientOptions.DebugLog(fmt.Sprintf("fetching of %s events exiting", url))

	lastCursor := ""
	for !a.doStop.WaitFor(30 * time.Second) {
		// The makeOneRequest function handles error
		// handling and fatal error handling.
		items, newCursor := a.makeOneRequest(url, lastCursor)
		lastCursor = newCursor
		if items == nil {
			continue
		}

		for _, item := range items {
			msg := &protocol.DataMessage{
				JsonPayload: item,
				TimestampMs: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
			}
			if err := a.uspClient.Ship(msg, 10*time.Second); err != nil {
				if err == uspclient.ErrorBufferFull {
					a.conf.ClientOptions.OnWarning("stream falling behind")
					err = a.uspClient.Ship(msg, 1*time.Hour)
				}
				if err == nil {
					continue
				}
				a.conf.ClientOptions.OnError(fmt.Errorf("Ship(): %v", err))
				a.doStop.Set()
				return
			}
		}
	}
}

func (a *ITGlueAdapter) makeOneRequest(url string, lastCursor string) ([]utils.Dict, string) {
	// Prepare the request body.
	reqData := opRequest{}

	a.conf.ClientOptions.DebugLog(fmt.Sprintf("requesting from %s starting now", url))
	reqData.StartTime = time.Now().Format(time.RFC3339)
	reqData.Limit = 1000

	b, err := json.Marshal(reqData)
	if err != nil {
		a.doStop.Set()
		return nil, ""
	}

	// Prepare the request.
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s?filter[created_at]=%s", a.endpoint, url, reqData.StartTime), bytes.NewBuffer(b))
	if err != nil {
		a.doStop.Set()
		return nil, ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.conf.Token))

	// Issue the request.
	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.conf.ClientOptions.OnError(fmt.Errorf("http.Client.Do(): %v", err))
		return nil, lastCursor
	}
	defer resp.Body.Close()

	// Evaluate if success.
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		a.conf.ClientOptions.OnWarning(fmt.Sprintf("itglue api non-200: %s\nREQUEST: %s\nRESPONSE: %s", resp.Status, string(b), string(body)))
		return nil, lastCursor
	}

	// Parse the response.
	respData := utils.Dict{}
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(&respData); err != nil {
		a.conf.ClientOptions.OnError(fmt.Errorf("itglue api invalid json: %v", err))
		return nil, lastCursor
	}

	// Report if a cursor was returned
	// as well as the items.
	lastCursor = respData.FindOneString("cursor")
	items, _ := respData.GetListOfDict("items")
	return items, lastCursor
}
