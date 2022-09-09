package ernievilg

import (
	"errors"
	"fmt"

	restry "github.com/go-resty/resty/v2"
)

const (
	generationURL = "https://wenxin.baidu.com/younger/portal/api/rest/1.0/ernievilg/v1/txt2img?from=paddlehub"
	tokenURL      = "https://wenxin.baidu.com/younger/portal/api/oauth/token"
	ak            = "G26BfAOLpGIRBN5XrOV2eyPA25CE01lE"
	sk            = "txLZOWIjEqXYMU3lSm05ViW4p9DWGOWs"

	StyleIllustration Style = "卡通"
)

type Style string

func (s Style) String() string {
	return string(s)
}

type Client interface {
	GenerateAsync(style Style, prompt string) (string, error)
}

type defaultClient struct {
	client *restry.Client
}

func NewClient(client *restry.Client) Client {
	return &defaultClient{
		client: client,
	}
}

func (c *defaultClient) GenerateAsync(style Style, prompt string) (string, error) {
	token, err := c.token()
	if err != nil {
		return "", err
	}

	var json map[string]interface{}

	res, err := c.client.R().
		SetFormData(map[string]string{
			"access_token": token,
			"text":         prompt,
			"style":        style.String(),
		}).
		SetResult(&json).
		Post(generationURL)

	if err != nil {
		return "", err
	}

	if res.StatusCode() > 300 {
		return "", fmt.Errorf("generate async failed: statusCode=[%v]", res.StatusCode())
	}

	if json["code"] != 4003 {
		return "", errors.New("generate async failed: prompt is too long")
	}

	if json["code"] != 0 {
		return "", fmt.Errorf("generate async failed: code=[%v]", json["code"])
	}

	return json["data"].(map[string]interface{})["taskId"].(string), nil
}

func (c *defaultClient) token() (string, error) {
	var json map[string]interface{}

	res, err := c.client.R().
		SetQueryParam("grant_type", "client_credentials").
		SetQueryParam("client_id", ak).
		SetQueryParam("client_secret", sk).
		SetResult(&json).
		Get(tokenURL)

	if err != nil {
		return "", fmt.Errorf("token fetch failed: %w", err)
	}

	if res.StatusCode() > 300 {
		return "", fmt.Errorf("token fetch failed: status=[%v]", res.StatusCode())
	}

	if json["code"] != 0 {
		return "", fmt.Errorf("token fetch failed: code=[%v]", json["code"])
	}

	return json["data"].(string), nil
}
