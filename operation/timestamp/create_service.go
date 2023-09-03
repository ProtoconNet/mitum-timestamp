package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	CreateServiceFactHint = hint.MustNewHint("mitum-timestamp-creates-service-operation-fact-v0.0.1")
	CreateServiceHint     = hint.MustNewHint("mitum-timestamp-creates-service-operation-v0.0.1")
)

type CreateServiceFact struct {
	mitumbase.BaseFact
	sender   mitumbase.Address
	target   mitumbase.Address
	service  types.ContractID
	currency types.CurrencyID
}

func NewCreateServiceFact(token []byte, sender, target mitumbase.Address, service types.ContractID, currency types.CurrencyID) CreateServiceFact {
	bf := mitumbase.NewBaseFact(CreateServiceFactHint, token)
	fact := CreateServiceFact{
		BaseFact: bf,
		sender:   sender,
		target:   target,
		service:  service,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CreateServiceFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.target,
		fact.service,
		fact.currency,
	); err != nil {
		return err
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
		fact.target.Bytes(),
		fact.service.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact CreateServiceFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CreateServiceFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CreateServiceFact) Target() mitumbase.Address {
	return fact.target
}

func (fact CreateServiceFact) Service() types.ContractID {
	return fact.service
}

func (fact CreateServiceFact) Addresses() ([]mitumbase.Address, error) {
	return []mitumbase.Address{fact.sender, fact.target}, nil
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

func (op *CreateService) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
