package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AppendFact) unpack(
	enc encoder.Encoder,
	sa, ta, pid string,
	rqts uint64,
	data, cid string,
) error {
	switch sender, err := mitumbase.DecodeAddress(sa, enc); {
	case err != nil:
		return err
	default:
		fact.sender = sender
	}

	switch contract, err := mitumbase.DecodeAddress(ta, enc); {
	case err != nil:
		return err
	default:
		fact.contract = contract
	}

	fact.projectID = pid
	fact.requestTimeStamp = rqts
	fact.data = data
	fact.currency = types.CurrencyID(cid)

	return nil
}
