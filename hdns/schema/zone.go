package schema

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

func (zt *Time) UnmarshalJSON(data []byte) error {
	var t time.Time

	err := json.Unmarshal(data, &t)
	if err == nil {
		*zt = Time{t}
		return nil
	}

	if string(data) == `""` {
		*zt = Time{}
		return nil
	}

	t, err = time.Parse(`"2006-01-02 15:04:05.999 -0700 MST"`, string(data))
	*zt = Time{t}
	return err
}

type ZoneTxtVerification struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Zone struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	TTL             int                 `json:"ttl"`
	Created         Time                `json:"created"`
	IsSecondaryDNS  bool                `json:"is_secondary_dns"`
	LegacyDNSHost   string              `json:"legacy_dns_host"`
	LegacyNS        []string            `json:"legacy_ns"`
	Modified        Time                `json:"modified"`
	NS              []string            `json:"ns"`
	Owner           string              `json:"owner"`
	Paused          bool                `json:"paused"`
	Permission      string              `json:"permission"`
	Project         string              `json:"project"`
	RecordsCount    int                 `json:"records_count"`
	Registrar       string              `json:"registrar"`
	Status          string              `json:"status"`
	TXTVerification ZoneTxtVerification `json:"txt_verification"`
	Verified        Time                `json:"verified"`
}

type ZoneCreateResponse struct {
	Zone Zone `json:"zone"`
}

type ZoneGetResponse struct {
	Zone Zone `json:"zone"`
}

type ZoneUpdateResponse struct {
	Zone Zone `json:"zone"`
}

type ZoneCreateRequest struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}

type ZoneUpdateRequest struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}

type ZoneAllResponse struct {
	Zones []Zone `json:"zones"`
}
