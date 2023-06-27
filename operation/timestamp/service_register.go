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
	ServiceRegisterFactHint = hint.MustNewHint("mitum-timestamp-service-register-operation-fact-v0.0.1")
	ServiceRegisterHint     = hint.MustNewHint("mitum-timestamp-service-register-operation-v0.0.1")
)

type ServiceRegisterFact struct {
	mitumbase.BaseFact
	sender   mitumbase.Address
	target   mitumbase.Address
	service  types.ContractID
	currency types.CurrencyID
}

func NewServiceRegisterFact(token []byte, sender, target mitumbase.Address, service types.ContractID, currency types.CurrencyID) ServiceRegisterFact {
	bf := mitumbase.NewBaseFact(ServiceRegisterFactHint, token)
	fact := ServiceRegisterFact{
		BaseFact: bf,
		sender:   sender,
		target:   target,
		service:  service,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ServiceRegisterFact) IsValid(b []byte) error {
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

func (fact ServiceRegisterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ServiceRegisterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ServiceRegisterFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.target.Bytes(),
		fact.service.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact ServiceRegisterFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact ServiceRegisterFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact ServiceRegisterFact) Target() mitumbase.Address {
	return fact.target
}

func (fact ServiceRegisterFact) Service() types.ContractID {
	return fact.service
}

func (fact ServiceRegisterFact) Addresses() ([]mitumbase.Address, error) {
	return []mitumbase.Address{fact.sender, fact.target}, nil
}

func (fact ServiceRegisterFact) Currency() types.CurrencyID {
	return fact.currency
}

type ServiceRegister struct {
	common.BaseOperation
}

func NewServiceRegister(fact ServiceRegisterFact) (ServiceRegister, error) {
	return ServiceRegister{BaseOperation: common.NewBaseOperation(ServiceRegisterHint, fact)}, nil
}

func (op *ServiceRegister) HashSign(priv mitumbase.Privatekey, networkID mitumbase.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
