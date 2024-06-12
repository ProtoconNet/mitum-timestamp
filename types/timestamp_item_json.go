package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type TimeStampItemJSONMarshaler struct {
	hint.BaseHinter
	ProjectID         string `json:"project_id"`
	RequestTimeStamp  uint64 `json:"request_timestamp"`
	ResponseTimeStamp uint64 `json:"response_timestamp"`
	TimeStampIdx      uint64 `json:"timestamp_idx"`
	Data              string `json:"data"`
}

func (t Item) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TimeStampItemJSONMarshaler{
		BaseHinter:        t.BaseHinter,
		ProjectID:         t.projectID,
		RequestTimeStamp:  t.requestTimeStamp,
		ResponseTimeStamp: t.responseTimeStamp,
		TimeStampIdx:      t.timestampIdx,
		Data:              t.data,
	})
}

type TimeStampItemJSONUnmarshaler struct {
	Hint              hint.Hint `json:"_hint"`
	ProjectID         string    `json:"project_id"`
	RequestTimeStamp  uint64    `json:"request_timestamp"`
	ResponseTimeStamp uint64    `json:"response_timestamp"`
	TimeStampIdx      uint64    `json:"timestamp_idx"`
	Data              string    `json:"data"`
}

func (t *Item) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of NFT")

	var u TimeStampItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return t.unmarshal(u.Hint, u.ProjectID, u.RequestTimeStamp, u.ResponseTimeStamp, u.TimeStampIdx, u.Data)
}
