package schema

type Record struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TTL      int    `json:"ttl"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	ZoneID   string `json:"zone_id"`
	Created  Time   `json:"created"`
	Modified Time   `json:"modified"`
}

type RecordCreateResponse struct {
	Record Record `json:"record"`
}

type RecordGetResponse struct {
	Record Record `json:"record"`
}

type RecordUpdateResponse struct {
	Record Record `json:"record"`
}

type RecordCreateRequest struct {
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type RecordUpdateRequest struct {
	Name   string `json:"name"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	ZoneID string `json:"zone_id"`
}

type RecordAllResponse struct {
	Records []Record `json:"records"`
}
