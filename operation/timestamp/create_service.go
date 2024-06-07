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
	CreateServiceFactHint = hint.MustNewHint("mitum-timestamp-create-service-operation-fact-v0.0.1")
	CreateServiceHint     = hint.MustNewHint("mitum-timestamp-create-service-operation-v0.0.1")
)

type CreateServiceFact struct {
	mitumbase.BaseFact
	sender   mitumbase.Address
	contract mitumbase.Address
	currency types.CurrencyID
}

func NewCreateServiceFact(token []byte, sender, contract mitumbase.Address, currency types.CurrencyID) CreateServiceFact {
	bf := mitumbase.NewBaseFact(CreateServiceFactHint, token)
	fact := CreateServiceFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateServiceFact) IsValid(b []byte) error {
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

func (fact CreateServiceFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateServiceFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateServiceFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CreateServiceFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CreateServiceFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CreateServiceFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact CreateServiceFact) Addresses() ([]mitumbase.Address, error) {
	return []mitumbase.Address{fact.sender, fact.contract}, nil
}

func (fact CreateServiceFact) Currency() types.CurrencyID {
	return fact.currency
}

type CreateService struct {
	common.BaseOperation
}

func NewCreateService(fact CreateServiceFact) (CreateService, error) {
	return CreateService{BaseOperation: common.NewBaseOperation(CreateServiceHint, fact)}, nil
}
