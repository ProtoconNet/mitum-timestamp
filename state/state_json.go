package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"design"`
}

func (sv DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		DesignStateValueJSONMarshaler(sv),
	)
}

type DesignStateValueJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	Design json.RawMessage `json:"design"`
}

func (sv *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	sv.BaseHinter = hint.NewBaseHinter(u.Hint)

	var sd types.Design
	if err := sd.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Design = sd

	return nil
}

type LastIdxStateValueJSONMarshaler struct {
	hint.BaseHinter
	ProjectID string `json:"project_id"`
	Index     uint64 `json:"last_idx"`
}

func (sv LastIdxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		LastIdxStateValueJSONMarshaler(sv),
	)
}

type LastIdxStateValueJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	ProjectID string    `json:"project_id"`
	Index     uint64    `json:"last_idx"`
}

func (sv *LastIdxStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of LastIdxStateValue")

	var u LastIdxStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	sv.BaseHinter = hint.NewBaseHinter(u.Hint)
	sv.ProjectID = u.ProjectID
	sv.Index = u.Index

	return nil
}

type TimeStampItemStateValueJSONMarshaler struct {
	hint.BaseHinter
	Item types.Item `json:"timestamp_item"`
}

func (sv ItemStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		TimeStampItemStateValueJSONMarshaler(sv),
	)
}

type ItemStateValueJSONUnmarshaler struct {
	Hint          hint.Hint       `json:"_hint"`
	TimeStampItem json.RawMessage `json:"timestamp_item"`
}

func (sv *ItemStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("decode json of ItemStateValue")

	var u ItemStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	sv.BaseHinter = hint.NewBaseHinter(u.Hint)

	var t types.Item
	if err := t.DecodeJSON(u.TimeStampItem, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Item = t

	return nil
}
