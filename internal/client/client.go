package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type InfoClient struct {
	client http.Client
}

const (
	apiAge         = "api.agify.io"
	apiGender      = "api.genderize.io"
	apiNationality = "api.nationalize.io"
)

var (
	ErrNoNationality   = errors.New("get 0 nationalities")
	ErrCannotGetAge    = errors.New("can't get age")
	ErrCannotGetGender = errors.New("can't get gender")
)

func New() InfoClient {
	return InfoClient{
		client: http.Client{},
	}
}

func (c *InfoClient) InfoAge(name string) (int64, error) {
	q := url.Values{}
	q.Add("name", name)

	data, err := c.doRequest(apiAge, q)
	if err != nil {
		return 0, err
	}

	var res InfoJsonAge
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	if res.Count == 0 {
		return 0, ErrCannotGetAge
	}

	return res.Age, nil
}

func (c *InfoClient) InfoGender(name string) (string, error) {
	q := url.Values{}
	q.Add("name", name)

	data, err := c.doRequest(apiGender, q)
	if err != nil {
		return "", err
	}

	var res InfoJsonGender
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	if res.Count == 0 {
		return "", ErrCannotGetGender
	}

	return res.Gender, nil
}

func (c *InfoClient) InfoNationality(name string) (string, error) {
	q := url.Values{}
	q.Add("name", name)

	data, err := c.doRequest(apiNationality, q)
	if err != nil {
		return "", err
	}

	var res InfoJsonNationality
	if err := json.Unmarshal(data, &res); err != nil {
		return "", err
	}

	if len(res.CountryProbs) == 0 {
		return "", ErrNoNationality
	}

	return res.CountryProbs[0].CountryId, nil
}

func (c *InfoClient) doRequest(host string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   host,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %s", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
