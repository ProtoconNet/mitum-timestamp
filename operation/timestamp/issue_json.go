package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type AppendFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender           mitumbase.Address `json:"sender"`
	Contract         mitumbase.Address `json:"contract"`
	ProjectID        string            `json:"project_id"`
	RequestTimeStamp uint64            `json:"request_timestamp"`
	Data             string            `json:"data"`
	Currency         types.CurrencyID  `json:"currency"`
}

func (fact IssueFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AppendFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		ProjectID:             fact.projectID,
		RequestTimeStamp:      fact.requestTimeStamp,
		Data:                  fact.data,
		Currency:              fact.currency,
	})
}

type AppendFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender           string `json:"sender"`
	Contract         string `json:"contract"`
	ProjectID        string `json:"project_id"`
	RequestTimeStamp uint64 `json:"request_timestamp"`
	Data             string `json:"data"`
	Currency         string `json:"currency"`
}

func (fact *IssueFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u AppendFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, u.Sender, u.Contract, u.ProjectID, u.RequestTimeStamp, u.Data, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type mintMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Issue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(mintMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Issue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
