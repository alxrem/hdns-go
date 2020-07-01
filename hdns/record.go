package hdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/alxrem/hdns-go/hdns/schema"
	"net/url"
)

type BaseRecord struct {
	Name   string
	TTL    int
	Type   string
	Value  string
	ZoneID string
}

type Record struct {
	BaseRecord
	ID       string
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

// ActionListOpts specifies options for listing actions.
type RecordListOpts struct {
	ListOpts
	ZoneID string
}

func (l RecordListOpts) values() url.Values {
	vals := l.ListOpts.values()
	if l.ZoneID != "" {
		vals.Add("zone_id", l.ZoneID)
	}
	return vals
}

// List returns a list of actions for a specific page.
func (c *RecordClient) List(ctx context.Context, opts RecordListOpts) ([]*Record, *Response, error) {
	path := "/records?" + opts.values().Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.RecordListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	records := make([]*Record, 0, len(body.Records))
	for _, record := range body.Records {
		records = append(records, RecordFromSchema(record))
	}
	return records, resp, nil
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

// RecordBulkCreateOpts specifies parameters for creating several Records at once.
type RecordBulkCreateOpts struct {
	Records []RecordCreateOpts
}

type RecordBulkCreateResult struct {
	InvalidRecords []*BaseRecord
	Records        []*Record
	ValidRecords   []*BaseRecord
}

// BulkCreate creates several Records at once.
func (c *RecordClient) BulkCreate(ctx context.Context, opts RecordBulkCreateOpts) (RecordBulkCreateResult, *Response, error) {
	reqBody := schema.RecordBulkCreateRequest{
		Records: []schema.RecordCreateRequest{},
	}

	for _, record := range opts.Records {
		reqRecordBody := schema.RecordCreateRequest{
			Name:   record.Name,
			TTL:    record.TTL,
			Type:   record.Type,
			Value:  record.Value,
			ZoneID: record.ZoneID,
		}
		reqBody.Records = append(reqBody.Records, reqRecordBody)
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return RecordBulkCreateResult{}, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/records/bulk", bytes.NewReader(reqBodyData))
	if err != nil {
		return RecordBulkCreateResult{}, nil, err
	}

	var respBody schema.RecordBulkCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return RecordBulkCreateResult{}, resp, err
	}

	result := RecordBulkCreateResult{
		InvalidRecords: BaseRecordsFromSchema(respBody.InvalidRecords),
		Records:        RecordsFromSchema(respBody.Records),
		ValidRecords:   BaseRecordsFromSchema(respBody.ValidRecords),
	}

	return result, resp, nil
}

// RecordBulkUpdateOpts specifies parameters for creating several Records at once.
type RecordBulkUpdateOpts struct {
	Records []RecordUpdateOpts
}

type RecordBulkUpdateResult struct {
	FailedRecords []*BaseRecord
	Records       []*Record
}

// BulkUpdate updates several Records at once.
func (c *RecordClient) BulkUpdate(ctx context.Context, opts RecordBulkUpdateOpts) (RecordBulkUpdateResult, *Response, error) {
	reqBody := schema.RecordBulkUpdateRequest{
		Records: []schema.RecordUpdateRequest{},
	}

	for _, record := range opts.Records {
		reqRecordBody := schema.RecordUpdateRequest{
			Name:   record.Name,
			TTL:    record.TTL,
			Type:   record.Type,
			Value:  record.Value,
			ZoneID: record.ZoneID,
		}
		reqBody.Records = append(reqBody.Records, reqRecordBody)
	}

	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return RecordBulkUpdateResult{}, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/records/bulk", bytes.NewReader(reqBodyData))
	if err != nil {
		return RecordBulkUpdateResult{}, nil, err
	}

	var respBody schema.RecordBulkUpdateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return RecordBulkUpdateResult{}, resp, err
	}

	result := RecordBulkUpdateResult{
		FailedRecords: BaseRecordsFromSchema(respBody.FailedRecords),
		Records:       RecordsFromSchema(respBody.Records),
	}

	return result, resp, nil
}
