package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	HostURL    string
	HttpClient *http.Client
	apiKey     string
}

func NewClient(host, apiKey *string) (*Client, error) {
	c := Client{
		HttpClient: &http.Client{},
	}

	if host != nil {
		c.HostURL = *host
	}

	if apiKey != nil {
		c.apiKey = *apiKey
	}

	return &c, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-API-Key", c.apiKey)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (b &[]byte) ValidateRequest (error) {
	var responses []postResponse
	err := json.Unmarshal(b, &responses)
	if err != nil {
		return err
	}

	for _, res := range responses {
		if res.Type != Success {
			return errors.New(res.Message)
		}
	}

	return nil
}

func (c *Client) GetAlias(id int64) (*AliasResponse, error) {
	url := c.HostURL + "/api/v1/get/alias/" + strconv.FormatInt(id, 10)

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var item AliasResponse
	err = json.Unmarshal(res, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) GetAllAliases() (*[]AliasResponse, error) {
	url := c.HostURL + "/api/v1/get/alias/all"

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var aliases []AliasResponse
	err = json.Unmarshal(res, &aliases)

	if err != nil {
		return nil, err
	}

	return &aliases, nil
}

func (c *Client) DeleteAlias(id int64) error {
	url := c.HostURL + "/api/v1/delete/alias"

	data := []byte(`[
		"` + strconv.FormatInt(id, 10) + `"
	]`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}

	var responses []postResponse
	err = json.Unmarshal(res, &responses)

	if responses[0].Type != Success {
		return errors.New(responses[0].Message)
	}

	return nil
}

func (c *Client) GetDomain(domain string) (*DomainResponse, error) {
	url := c.HostURL + "/api/v1/get/domain/" + domain

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var item DomainResponse
	err = json.Unmarshal(res, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) GetAllDomains() (*[]DomainResponse, error) {
	url := c.HostURL + "/api/v1/get/domain/all"

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var domains []DomainResponse
	err = json.Unmarshal(res, &domains)

	if err != nil {
		return nil, err
	}

	return &domains, nil
}

func (c *Client) DeleteDomain(domain string) error {
	url := c.HostURL + "/api/v1/delete/domain"

	data := []byte(`[
		"` + domain + `"
	]`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}

	var responses []postResponse
	err = json.Unmarshal(res, &responses)

	if responses[0].Type != Success {
		return errors.New(responses[0].Message)
	}

	return nil
}

func (c *Client) GetMailbox(username string) (*MailboxResponse, error) {
	url := c.HostURL + "/api/v1/get/mailbox/" + username

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var item MailboxResponse
	err = json.Unmarshal(res, &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Client) GetAllMailboxes() (*[]MailboxResponse, error) {
	url := c.HostURL + "/api/v1/get/mailbox/all"

	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var mailboxes []MailboxResponse
	err = json.Unmarshal(res, &mailboxes)

	if err != nil {
		return nil, err
	}

	return &mailboxes, nil
}

func (c *Client) DeleteMailbox(mailbox string) error {
	url := c.HostURL + "/api/v1/delete/mailbox"

	data := []byte(`[
		"` + mailbox + `"
	]`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}

	var responses []postResponse
	err = json.Unmarshal(res, &responses)

	if responses[0].Type != Success {
		return errors.New(responses[0].Message)
	}

	return nil
}
