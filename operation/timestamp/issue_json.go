package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type IssueFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender           base.Address     `json:"sender"`
	Contract         base.Address     `json:"contract"`
	ProjectID        string           `json:"project_id"`
	RequestTimeStamp uint64           `json:"request_timestamp"`
	Data             string           `json:"data"`
	Currency         types.CurrencyID `json:"currency"`
}

func (fact IssueFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(IssueFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		ProjectID:             fact.projectID,
		RequestTimeStamp:      fact.requestTimeStamp,
		Data:                  fact.data,
		Currency:              fact.currency,
	})
}

type IssueFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender           string `json:"sender"`
	Contract         string `json:"contract"`
	ProjectID        string `json:"project_id"`
	RequestTimeStamp uint64 `json:"request_timestamp"`
	Data             string `json:"data"`
	Currency         string `json:"currency"`
}

func (fact *IssueFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u IssueFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, u.Sender, u.Contract, u.ProjectID, u.RequestTimeStamp, u.Data, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type OperationMarshaler struct {
	common.BaseOperationJSONMarshaler
	extras.BaseOperationExtensionsJSONMarshaler
}

func (op Issue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Issue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
