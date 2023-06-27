package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
)

type AppendFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender           mitumbase.Address `json:"sender"`
	Target           mitumbase.Address `json:"target"`
	Service          types.ContractID  `json:"service"`
	ProjectID        string            `json:"projectid"`
	RequestTimeStamp uint64            `json:"request_timestamp"`
	Data             string            `json:"data"`
	Currency         types.CurrencyID  `json:"currency"`
}

func (fact AppendFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AppendFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Target:                fact.target,
		Service:               fact.service,
		ProjectID:             fact.projectID,
		RequestTimeStamp:      fact.requestTimeStamp,
		Data:                  fact.data,
		Currency:              fact.currency,
	})
}

type AppendFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender           string `json:"sender"`
	Target           string `json:"target"`
	Service          string `json:"service"`
	ProjectID        string `json:"projectid"`
	RequestTimeStamp uint64 `json:"request_timestamp"`
	Data             string `json:"data"`
	Currency         string `json:"currency"`
}

func (fact *AppendFact) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of AppendFact")

	var u AppendFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	return fact.unmarshal(enc, u.Sender, u.Target, u.Service, u.ProjectID, u.RequestTimeStamp, u.Data, u.Currency)
}

type mintMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op Append) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(mintMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *Append) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Mint")

	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}