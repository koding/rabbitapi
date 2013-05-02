package rabbitapi

import (
	"bytes"
	"encoding/json"
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

// AlivenessTest declares a test queue for the given vhost, then publishes and
// consumes a message. Intended for use by monitoring tools. If everything is
// working correctly, will return HTTP status 200 with body. Note: the test
// queue will not be deleted (to to prevent queue churn if this is repeatedly
// pinged).
func (r *Rabbit) AlivenessTest(vhost string) error {
	if vhost == "/" {
		vhost = "%2f"
	}

	body, err := r.getRequest("/api/aliveness-test/" + vhost)
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

func (r *Rabbit) getRequest(endpoint string) ([]byte, error) {
	req, err := r.newRequest("GET", endpoint, nil)
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Rabbit) putRequest(endpoint string, body []byte) error {
	reader := bytes.NewBuffer(body)
	req, err := r.newRequest("PUT", endpoint, reader)
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

	if resp.StatusCode != 204 {
		return fmt.Errorf(resp.Status)
	}

	return nil
}

func (r *Rabbit) deleteRequest(endpoint string) error {
	req, err := r.newRequest("DELETE", endpoint, nil)
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

	if resp.StatusCode != 204 {
		return fmt.Errorf(resp.Status)
	}

	return nil
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
