package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type CreateServiceFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender   mitumbase.Address `json:"sender"`
	Contract mitumbase.Address `json:"contract"`
	Currency types.CurrencyID  `json:"currency"`
}

func (fact CreateServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Currency:              fact.currency,
	})
}

type CreateServiceFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	Currency string `json:"currency"`
}

func (fact *CreateServiceFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u CreateServiceFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, u.Sender, u.Contract, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type createServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(createServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateService) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
