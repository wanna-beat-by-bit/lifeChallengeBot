package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	ce "tgBot/pkg/customError"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(token string, host string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Request(method string, getParams url.Values) (bytes []byte, err error) {
	defer func() { err = ce.WrapIfError("Error while making request", err) }()

	urlForRequest := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: getParams.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, urlForRequest.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Updates(offset int, limit int) (ups []Update, err error) {
	defer func() { err = ce.WrapIfError("Error while fetching updates", err) }()

	getParams := url.Values{}
	getParams.Add("offset", strconv.Itoa(offset))
	getParams.Add("limit", strconv.Itoa(limit))

	data, err := c.Request(getUpdatesMethod, getParams)
	if err != nil {
		return nil, err
	}

	var bucket MessageWrapper
	if err := json.Unmarshal(data, &bucket); err != nil {
		return nil, err
	}

	if !bucket.Ok {
		return nil, err
	}

	return bucket.Result, nil
}

func (c *Client) SendMessage(text string, chatId int) error {
	getParams := url.Values{}
	getParams.Add("text", text)
	getParams.Add("chat_id", strconv.Itoa(chatId))

	_, err := c.Request(sendMessageMethod, getParams)
	if err != nil {
		return ce.Wrap("Error while sending message", err)
	}

	return nil
}
