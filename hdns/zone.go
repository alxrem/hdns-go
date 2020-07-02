package hdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alxrem/hdns-go/hdns/schema"
)

type ZoneTxtVerification struct {
	Name  string
	Token string
}

type Zone struct {
	ID              string
	Name            string
	TTL             int
	Created         schema.Time
	IsSecondaryDNS  bool
	LegacyDNSHost   string
	LegacyNS        []string
	Modified        schema.Time
	NS              []string
	Owner           string
	Paused          bool
	Permission      string
	Project         string
	RecordsCount    int
	Registrar       string
	Status          string
	TXTVerification ZoneTxtVerification
	Verified        schema.Time
}

type ZoneClient struct {
	client *Client
}

func (c *ZoneClient) GetByID(ctx context.Context, id string) (*Zone, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/zones/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ZoneGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}
	return ZoneFromSchema(body.Zone), resp, nil
}

func (c *ZoneClient) All(ctx context.Context) ([]*Zone, error) {
	req, err := c.client.NewRequest(ctx, "GET", "/zones", nil)
	if err != nil {
		return nil, err
	}

	var body schema.ZoneAllResponse
	_, err = c.client.Do(req, &body)
	if err != nil {
		return nil, err
	}

	zones := make([]*Zone, 0, len(body.Zones))
	for _, z := range body.Zones {
		zones = append(zones, ZoneFromSchema(z))
	}

	return zones, nil
}

// ZoneCreateOpts specifies parameters for creating a Zone.
type ZoneCreateOpts struct {
	Name string
	TTL  int
}

// Create creates a Zone.
func (c *ZoneClient) Create(ctx context.Context, opts ZoneCreateOpts) (*Zone, *Response, error) {
	reqBody := schema.ZoneCreateRequest{
		Name: opts.Name,
		TTL:  opts.TTL,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/zones", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.ZoneCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(respBody.Zone), resp, nil
}

func (c *ZoneClient) Delete(ctx context.Context, id string) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/zones/%s", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// ZoneUpdateOpts specifies parameters for updating a Zone.
type ZoneUpdateOpts struct {
	Name string
	TTL  int
}

// Update updates a Zone.
func (c *ZoneClient) Update(ctx context.Context, id string, opts ZoneUpdateOpts) (*Zone, *Response, error) {
	reqBody := schema.ZoneUpdateRequest{
		Name: opts.Name,
		TTL:  opts.TTL,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/zones/%s", id)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.ZoneUpdateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}

	return ZoneFromSchema(respBody.Zone), resp, nil
}
