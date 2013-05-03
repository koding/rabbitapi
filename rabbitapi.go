package rabbitapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Rabbit struct {
	Username string
	Password string
	Url      string
}

type Status struct {
	Status string
}

type Overview struct {
	Contexts []struct {
		Description string `json:"description"`
		Node        string `json:"node"`
		Path        string `json:"path"`
		Port        int    `json:"port"`
	} `json:"contexts"`
	ErlangVersion string `json:"erlang_version"`
	ExchangeTypes []struct {
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		Name        string `json:"name"`
	} `json:"exchange_types"`
	Listeners []struct {
		IpAdress string `json:"ip_address"`
		Node     string `json:"node"`
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
	} `json:"listeners"`
	ManagementVersion string   `json:"management_version"`
	MessageStats      []string `json:"message_stats"`
	Node              string   `json:"node"`
	ObjectTotals      struct {
		Channels    int `json:"channels"`
		Connections int `json:"connections`
		Consumers   int `json:"consumers`
		Exchanges   int `json:"exchanges`
		Queues      int `json:"queues`
	} `json:"object_totals"`
	QueueTotals struct {
		Messages        int `json:"messages"`
		MessagesDetails struct {
			Interval  int     `json:"interval"`
			LastEvent int     `json:"last_event"`
			Rate      float64 `json:"rate"`
		} `json:"messages_details"`
		MessagesReady        int `json:"messages_ready"`
		MessagesReadyDetails struct {
			Interval  int     `json:"interval"`
			LastEvent int     `json:"last_event"`
			Rate      float64 `json:"rate"`
		} `json:"messages_ready_details"`
		MessagesUnacknowledged        int `json:"messages_unacknowledged"`
		MessagesUnacknowledgedDetails struct {
			Interval  int     `json:"interval"`
			LastEvent int     `json:"last_event"`
			Rate      float64 `json:"rate"`
		} `json:"messages_unacknowledged_details"`
	} `json:"queue_totals"`
	RabbitmqVersion  string `json:"rabbitmq_version"`
	StatisticsDbNode string `json:"statistics_db_node"`
	StatisticsLevel  string `json:"statistics_level"`
}

// Auth is used to create the initial struct that will used for all api calls.
// If you pass empty strings, default values will be:
//
//    Username : guest
//    Password : guest
//    Url:  http://localhost:15672
func Auth(username, password, url string) *Rabbit {
	if username == "" {
		username = "guest"
	}

	if password == "" {
		password = "guest"
	}

	if url == "" {
		url = "http://localhost:15672"
	}

	return &Rabbit{
		Username: username,
		Password: password,
		Url:      url,
	}
}

// Various random bits of information that describe the whole system, like
// number of exchanges, connection information, erlang version, etc..
func (r *Rabbit) GetOverview() (Overview, error) {
	body, err := r.doRequest("GET", "/api/overview", nil)
	if err != nil {
		return Overview{}, err
	}

	overview := Overview{}
	err = json.Unmarshal(body, &overview)
	if err != nil {
		return Overview{}, err
	}

	return overview, nil
}

// AlivenessTest declares a test queue (with the name alivness test) for the
// given vhost, then publishes and consumes a message. Intended for use by
// monitoring tools. If everything is working correctly, will return an error
// of type nil. Note: the test queue will not be deleted (to to prevent queue
// churn if this is repeatedly pinged).
func (r *Rabbit) AlivenessTest(vhost string) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.doRequest("GET", "/api/aliveness-test/"+vhost, nil)
	if err != nil {
		return err
	}

	status := Status{}
	err = json.Unmarshal(body, &status)
	if err != nil {
		return err
	}

	if status.Status != "ok" {
		return fmt.Errorf("there is a problem on vhost '%s'", vhost)
	}
	return nil

}

// Our custom HTTP Request wrapper
func (r *Rabbit) doRequest(method, endpoint string, body []byte) ([]byte, error) {
	readerBody := bytes.NewBuffer(body)
	req, err := r.newRequest(method, endpoint, readerBody)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	switch method {
	case "PUT", "DELETE":
		if resp.StatusCode != 204 {
			return nil, fmt.Errorf(resp.Status)
		}
		return nil, nil
	case "GET":
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf(resp.Status)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return responseBody, nil
	default:
		return nil, errors.New("Method is not supported")
	}
}

// modified version of http.NewRequest to not escape %2f paths. unfortunaley
// rabbitmq uses a RESTful api and "/" is a resource for a lot of api calls
func (r *Rabbit) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	requestUrl := r.Url + endpoint
	u, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	u.Opaque = endpoint // get around the path encoding bug
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	req := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       rc,
		Host:       u.Host,
	}

	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
		}
	}

	return req, nil
}
