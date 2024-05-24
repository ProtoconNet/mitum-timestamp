package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CreateServiceFact) unpack(
	enc encoder.Encoder,
	sa, ta, cid string,
) error {
	fact.currency = types.CurrencyID(cid)

	sender, err := mitumbase.DecodeAddress(sa, enc)
	if err != nil {
		return err
	}
	fact.sender = sender
	contract, err := mitumbase.DecodeAddress(ta, enc)
	if err != nil {
		return err
	}
	fact.contract = contract

	return nil
}
