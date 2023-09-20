package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *AppendFact) unmarshal(
	enc encoder.Encoder,
	sa, ta, pid string,
	rqts uint64,
	data, cid string,
) error {
	e := util.StringError("failed to unmarshal AppendFact")

	switch sender, err := mitumbase.DecodeAddress(sa, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.sender = sender
	}

	switch target, err := mitumbase.DecodeAddress(ta, enc); {
	case err != nil:
		return e.Wrap(err)
	default:
		fact.target = target
	}

	fact.projectID = pid
	fact.requestTimeStamp = rqts
	fact.data = data
	fact.currency = types.CurrencyID(cid)

	return nil
}
