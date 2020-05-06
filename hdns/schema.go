package hdns

import (
	"gitlab.com/alxrem/hdns-go/hdns/schema"
)

// This file provides converter functions to convert models in the
// schema package to models in the hdns package.

// ZoneFromSchema converts a schema.Zone to a Zone.
func ZoneFromSchema(s schema.Zone) *Zone {
	zone := &Zone{
		ID:             s.ID,
		Name:           s.Name,
		TTL:            s.TTL,
		Created:        s.Created,
		IsSecondaryDNS: s.IsSecondaryDNS,
		LegacyDNSHost:  s.LegacyDNSHost,
		LegacyNS:       []string{},
		Modified:       s.Modified,
		NS:             []string{},
		Owner:          s.Owner,
		Paused:         s.Paused,
		Permission:     s.Permission,
		Project:        s.Project,
		RecordsCount:   s.RecordsCount,
		Status:         s.Status,
		TXTVerification: ZoneTxtVerification{
			Name:  s.TXTVerification.Name,
			Token: s.TXTVerification.Token,
		},
		Verified: s.Verified,
	}
	for _, ns := range s.LegacyNS {
		zone.LegacyNS = append(zone.LegacyNS, ns)
	}
	for _, ns := range s.NS {
		zone.NS = append(zone.LegacyNS, ns)
	}
	return zone
}

// BaseRecordFromSchema converts a schema.BaseRecord to a BaseRecord.
func BaseRecordFromSchema(s schema.BaseRecord) *BaseRecord {
	return &BaseRecord{
		Name:   s.Name,
		TTL:    s.TTL,
		Type:   s.Type,
		Value:  s.Value,
		ZoneID: s.ZoneID,
	}
}

// RecordFromSchema converts a schema.Record to a Record.
func RecordFromSchema(s schema.Record) *Record {
	return &Record{
		BaseRecord: *BaseRecordFromSchema(s.BaseRecord),
		ID:         s.ID,
		Created:    s.Created,
		Modified:   s.Modified,
	}
}

func BaseRecordsFromSchema(s []schema.BaseRecord) []*BaseRecord {
	var baseRecords []*BaseRecord
	for _, r := range s {
		baseRecords = append(baseRecords, BaseRecordFromSchema(r))
	}
	return baseRecords
}

func RecordsFromSchema(s []schema.Record) []*Record {
	var records []*Record
	for _, r := range s {
		records = append(records, RecordFromSchema(r))
	}
	return records
}

// PaginationFromSchema converts a schema.MetaPagination to a Pagination.
func PaginationFromSchema(s schema.MetaPagination) Pagination {
	return Pagination{
		Page:         s.Page,
		PerPage:      s.PerPage,
		PreviousPage: s.PreviousPage,
		NextPage:     s.NextPage,
		LastPage:     s.LastPage,
		TotalEntries: s.TotalEntries,
	}
}

// ErrorFromSchema converts a schema.Error to an Error.
func ErrorFromSchema(s schema.Error) Error {
	e := Error{
		Code:    ErrorCode(s.Code),
		Message: s.Message,
	}

	switch d := s.Details.(type) {
	case schema.ErrorDetailsInvalidInput:
		details := ErrorDetailsInvalidInput{
			Fields: []ErrorDetailsInvalidInputField{},
		}
		for _, field := range d.Fields {
			details.Fields = append(details.Fields, ErrorDetailsInvalidInputField{
				Name:     field.Name,
				Messages: field.Messages,
			})
		}
		e.Details = details
	}
	return e
}
