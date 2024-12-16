package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *IssueFact) unpack(
	enc encoder.Encoder,
	sa, ta, pid string,
	rqts uint64,
	data, cid string,
) error {
	switch sender, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return err
	default:
		fact.sender = sender
	}

	switch contract, err := base.DecodeAddress(ta, enc); {
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
