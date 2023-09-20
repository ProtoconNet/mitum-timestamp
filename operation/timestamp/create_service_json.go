package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type CreateServiceFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender   mitumbase.Address `json:"sender"`
	Target   mitumbase.Address `json:"target"`
	Currency types.CurrencyID  `json:"currency"`
}

func (fact CreateServiceFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(CreateServiceFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Target:                fact.target,
		Currency:              fact.currency,
	})
}

type CreateServiceFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Target   string `json:"target"`
	Currency string `json:"currency"`
}

func (fact *CreateServiceFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateServiceFact")

	var u CreateServiceFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Target, u.Currency)
}

type createServiceMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op CreateService) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(createServiceMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *CreateService) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of CreateService")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
