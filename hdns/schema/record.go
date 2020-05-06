package schema

type BaseRecord struct {
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type Record struct {
	BaseRecord
	ID       string `json:"id"`
	Created  Time   `json:"created"`
	Modified Time   `json:"modified"`
}

type RecordCreateRequest struct {
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type RecordCreateResponse struct {
	Record Record `json:"record"`
}

type RecordGetResponse struct {
	Record Record `json:"record"`
}

type RecordUpdateRequest struct {
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type RecordUpdateResponse struct {
	Record Record `json:"record"`
}

type RecordBulkCreateRequest struct {
	Records []RecordCreateRequest `json:"records"`
}

type RecordBulkCreateResponse struct {
	InvalidRecords []BaseRecord `json:"invalid_records"`
	Records        []Record     `json:"records"`
	ValidRecords   []BaseRecord `json:"valid_records"`
}

type RecordBulkUpdateRequest struct {
	Records []RecordUpdateRequest `json:"records"`
}

type RecordBulkUpdateResponse struct {
	FailedRecords []BaseRecord `json:"failed_records"`
	Records       []Record     `json:"records"`
}

type RecordAllResponse struct {
	Records []Record `json:"records"`
}
