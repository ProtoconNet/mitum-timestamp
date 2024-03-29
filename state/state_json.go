package state

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-timestamp/types"

	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ServiceDesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"design"`
}

func (s ServiceDesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		ServiceDesignStateValueJSONMarshaler(s),
	)
}

type ServiceDesignStateValueJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	Design json.RawMessage `json:"design"`
}

func (s *ServiceDesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of ServiceDesignStateValue")

	var u ServiceDesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var sd types.Design
	if err := sd.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	s.Design = sd

	return nil
}

type TimeStampLastIndexStateValueJSONMarshaler struct {
	hint.BaseHinter
	ProjectID string `json:"projectid"`
	Index     uint64 `json:"index"`
}

func (s TimeStampLastIndexStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		TimeStampLastIndexStateValueJSONMarshaler(s),
	)
}

type TimeStampLastIndexStateValueJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	ProjectID string    `json:"projectid"`
	Index     uint64    `json:"index"`
}

func (s *TimeStampLastIndexStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TimeStampLastIndexStateValue")

	var u TimeStampLastIndexStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)
	s.ProjectID = u.ProjectID
	s.Index = u.Index

	return nil
}

type TimeStampItemStateValueJSONMarshaler struct {
	hint.BaseHinter
	TimeStampItem types.TimeStampItem `json:"timestampitem"`
}

func (s TimeStampItemStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		TimeStampItemStateValueJSONMarshaler(s),
	)
}

type TimeStampItemStateValueJSONUnmarshaler struct {
	Hint          hint.Hint       `json:"_hint"`
	TimeStampItem json.RawMessage `json:"timestampitem"`
}

func (s *TimeStampItemStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of TimeStampItemStateValue")

	var u TimeStampItemStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var t types.TimeStampItem
	if err := t.DecodeJSON(u.TimeStampItem, enc); err != nil {
		return e.Wrap(err)
	}
	s.TimeStampItem = t

	return nil
}
