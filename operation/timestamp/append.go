package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	AppendFactHint = hint.MustNewHint("mitum-timestamp-append-operation-fact-v0.0.1")
	AppendHint     = hint.MustNewHint("mitum-timestamp-append-operation-v0.0.1")
)

type AppendFact struct {
	mitumbase.BaseFact
	sender           mitumbase.Address
	contract         mitumbase.Address
	projectID        string
	requestTimeStamp uint64
	data             string
	currency         currencytypes.CurrencyID
}

func NewAppendFact(
	token []byte, sender, contract mitumbase.Address, projectID string,
	requestTimeStamp uint64, data string, currency currencytypes.CurrencyID) AppendFact {
	bf := mitumbase.NewBaseFact(AppendFactHint, token)
	fact := AppendFact{
		BaseFact:         bf,
		sender:           sender,
		contract:         contract,
		projectID:        projectID,
		requestTimeStamp: requestTimeStamp,
		data:             data,
		currency:         currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact AppendFact) IsValid(b []byte) error {
	if len(fact.projectID) < 1 || len(fact.projectID) > types.MaxProjectIDLen {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf(
					"invalid projectID length %v < 1 or > %v", len(fact.projectID), types.MaxProjectIDLen)))
	}

	if !currencytypes.ReSpcecialChar.Match([]byte(fact.projectID)) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("projectID ID %s, must match regex `^[^\\s:/?#\\[\\]@]*$`", fact.projectID)))
	}

	if len(fact.data) < 1 || len(fact.data) > types.MaxDataLen {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf("invalid data length %v < 1 or > %v", len(fact.data), types.MaxDataLen)))
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact AppendFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact AppendFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact AppendFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.projectID),
		util.Uint64ToBytes(fact.requestTimeStamp),
		[]byte(fact.data),
		fact.currency.Bytes(),
	)
}

func (fact AppendFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact AppendFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact AppendFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact AppendFact) ProjectId() string {
	return fact.projectID
}

func (fact AppendFact) RequestTimeStamp() uint64 {
	return fact.requestTimeStamp
}

func (fact AppendFact) Data() string {
	return fact.data
}

func (fact AppendFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

type Append struct {
	common.BaseOperation
}

func NewAppend(fact AppendFact) (Append, error) {
	return Append{BaseOperation: common.NewBaseOperation(AppendHint, fact)}, nil
}
