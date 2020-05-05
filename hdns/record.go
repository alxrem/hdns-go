package hdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/alxrem/hdns-go/hdns/schema"
)

type Record struct {
	ID       string
	Name     string
	TTL      int
	Type     string
	Value    string
	ZoneID   string
	Created  schema.Time
	Modified schema.Time
}

type RecordClient struct {
	client *Client
}

func (c *RecordClient) GetByID(ctx context.Context, id string) (*Record, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/records/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, resp, err
	}
	return RecordFromSchema(body.Record), resp, nil
}

func (c *RecordClient) All(ctx context.Context) ([]*Record, error) {
	req, err := c.client.NewRequest(ctx, "GET", "/records", nil)
	if err != nil {
		return nil, err
	}

	var body schema.RecordAllResponse
	_, err = c.client.Do(req, &body)
	if err != nil {
		return nil, err
	}

	records := make([]*Record, 0, len(body.Records))
	for _, z := range body.Records {
		records = append(records, RecordFromSchema(z))
	}

	return records, nil
}

// RecordCreateOpts specifies parameters for creating a Record.
type RecordCreateOpts struct {
	Name   string
	TTL    int
	Type   string
	Value  string
	ZoneID string
}

// Create creates a Record.
func (c *RecordClient) Create(ctx context.Context, opts RecordCreateOpts) (*Record, *Response, error) {
	reqBody := schema.RecordCreateRequest{
		Name:   opts.Name,
		TTL:    opts.TTL,
		Type:   opts.Type,
		Value:  opts.Value,
		ZoneID: opts.ZoneID,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/records", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.RecordCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}

	return RecordFromSchema(respBody.Record), resp, nil
}

func (c *RecordClient) Delete(ctx context.Context, id string) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/records/%s", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}

// RecordUpdateOpts specifies parameters for updating a Record.
type RecordUpdateOpts struct {
	Name   string
	TTL    int
	Type   string
	Value  string
	ZoneID string
}

// Update updates a Record.
func (c *RecordClient) Update(ctx context.Context, id string, opts RecordUpdateOpts) (*Record, *Response, error) {
	reqBody := schema.RecordUpdateRequest{
		Name: opts.Name,
		TTL:  opts.TTL,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	path := fmt.Sprintf("/records/%s", id)
	req, err := c.client.NewRequest(ctx, "PUT", path, bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody schema.RecordUpdateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}

	return RecordFromSchema(respBody.Record), resp, nil
}
