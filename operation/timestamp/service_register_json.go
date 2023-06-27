package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type ServiceRegisterFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender   mitumbase.Address `json:"sender"`
	Target   mitumbase.Address `json:"target"`
	Service  types.ContractID  `json:"service"`
	Currency types.CurrencyID  `json:"currency"`
}

func (fact ServiceRegisterFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ServiceRegisterFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Target:                fact.target,
		Service:               fact.service,
		Currency:              fact.currency,
	})
}

type ServiceRegisterFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Target   string `json:"target"`
	Service  string `json:"service"`
	Currency string `json:"currency"`
}

func (fact *ServiceRegisterFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of ServiceRegisterFact")

	var u ServiceRegisterFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Target, u.Service, u.Currency)
}

type serviceRegisterMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op ServiceRegister) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(serviceRegisterMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *ServiceRegister) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of ServiceRegister")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
