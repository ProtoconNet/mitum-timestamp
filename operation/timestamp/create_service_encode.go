package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateServiceFact) unmarshal(
	enc encoder.Encoder,
	sa,
	ta,
	svc,
	cid string,
) error {
	e := util.StringError("failed to unmarshal CreateServiceFact")

	fact.currency = types.CurrencyID(cid)

	sender, err := mitumbase.DecodeAddress(sa, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.sender = sender
	target, err := mitumbase.DecodeAddress(ta, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.target = target
	fact.service = types.ContractID(svc)

	return nil
}
