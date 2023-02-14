package types

import (
	"bytes"
	"encoding/json"
)

const (
	HTTPSingleExecutionMode   = "httpSingle"
	HTTPMultipleExecutionMode = "httpMultiple"
)

type NullExecutionMode struct {
	Value string
	Valid bool
}

func NewNullExecutionMode(s string, valid bool) NullExecutionMode {
	return NullExecutionMode{s, valid}
}

func NullExecutionModeFrom(s string) NullExecutionMode {
	return NullExecutionMode{s, true}
}

func (em NullExecutionMode) ValueOrZero() string {
	if !em.Valid {
		return ""
	}
	return em.Value
}

func (em NullExecutionMode) MarshalJSON() ([]byte, error) {
	if !em.Valid {
		return []byte(`null`), nil
	}
	return json.Marshal(em.Value)
}

func (em *NullExecutionMode) UnmarshalJSON(data []byte) error {
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
