package types

import (
	"bytes"
	"encoding/json"
)

const (
	HTTPSingleLoadZone   = "http_single"
	HTTPMultipleLoadZone = "http_multiple"
)

type LoadZone struct {
	Location string `json:"location"`
	Fraction int    `json:"fraction"`
}

type NullLoadDistribution struct {
	Value []LoadZone
	Valid bool
}

func NewNullLoadDistribution(loadZone []LoadZone, valid bool) NullLoadDistribution {
	return NullLoadDistribution{loadZone, valid}
}

func NullLoadDistributionFrom(loadZone []LoadZone) NullLoadDistribution {
	return NullLoadDistribution{loadZone, true}
}

func (em NullLoadDistribution) MarshalJSON() ([]byte, error) {
	if !em.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(em.Value)
}

func (em *NullLoadDistribution) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`null`)) {
		em.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &em.Value); err != nil {
		return err
	}
	em.Valid = true
	return nil
}
