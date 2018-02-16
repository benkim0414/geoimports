package imports

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/benkim0414/geoauth"
)

const (
	baseURL = "https://api.geocreation.com.au/api"
)

type ImportType struct {
	ID                        string `json:"_id"`
	Name                      string `json:"name"`
	Type                      string `json:"type"`
	AvailableVersionSpottedAt string `json:"availableVersionSpottedAt"`
}

func (t *ImportType) Available() bool {
	spottedAt, err := time.Parse(time.RFC3339, t.AvailableVersionSpottedAt)
	if err != nil {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return today.Equal(spottedAt.Truncate(24 * time.Hour))
}

func (t ImportType) String() string {
	format := "%s %q"
	if t.Available() {
		format += " is available to import"
	} else {
		format += " is unavailable to import"
	}
	return fmt.Sprintf(format, t.Name, t.ID)
}

type Import struct {
	ID      string `json:"_id"`
	Type    string `json:"importTypeType"`
	Version string `json:"versionName"`
	Status  string `json:"status"`
}

// Client is a client for interacting with GEO imports.
type Client struct {
	hc *http.Client
}

// NewClient creates a new GEO imports client.
func NewClient(ctx context.Context, conf *geoauth.Config) (*Client, error) {
	tok, err := conf.KMSCredentialsToken(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		hc: conf.Client(ctx, tok),
	}, nil
}

func (c *Client) GetImportTypes(pull bool) ([]*ImportType, error) {
	url := fmt.Sprintf("%s/import_types/?pull=%t", baseURL, pull)
	res, err := c.hc.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response struct {
		ImportType struct {
			Result []*ImportType `json:"result"`
		} `json:"importType"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response.ImportType.Result, nil
}

func (c *Client) PutImport(id string) (*Import, error) {
	url := fmt.Sprintf("%s/imports/%s/transition", baseURL, id)
	payload := `{"status": "queued"}`
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response struct {
		Import struct {
			Result *Import `json:"result"`
		} `json:"import"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response.Import.Result, nil
}
