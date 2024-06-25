package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	RegisterModelFactHint = hint.MustNewHint("mitum-timestamp-register-model-operation-fact-v0.0.1")
	RegisterModelHint     = hint.MustNewHint("mitum-timestamp-register-model-operation-v0.0.1")
)

type RegisterModelFact struct {
	mitumbase.BaseFact
	sender   mitumbase.Address
	contract mitumbase.Address
	currency types.CurrencyID
}

func NewRegisterModelFact(token []byte, sender, contract mitumbase.Address, currency types.CurrencyID) RegisterModelFact {
	bf := mitumbase.NewBaseFact(RegisterModelFactHint, token)
	fact := RegisterModelFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RegisterModelFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if err := util.CheckIsValiders(nil, false,
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

func (fact RegisterModelFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RegisterModelFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RegisterModelFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact RegisterModelFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact RegisterModelFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact RegisterModelFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact RegisterModelFact) Addresses() ([]mitumbase.Address, error) {
	return []mitumbase.Address{fact.sender, fact.contract}, nil
}

func (fact RegisterModelFact) Currency() types.CurrencyID {
	return fact.currency
}

type RegisterModel struct {
	common.BaseOperation
}

func NewRegisterModel(fact RegisterModelFact) (RegisterModel, error) {
	return RegisterModel{BaseOperation: common.NewBaseOperation(RegisterModelHint, fact)}, nil
}
