package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	IssueFactHint = hint.MustNewHint("mitum-timestamp-issue-operation-fact-v0.0.1")
	IssueHint     = hint.MustNewHint("mitum-timestamp-issue-operation-v0.0.1")
)

type IssueFact struct {
	base.BaseFact
	sender           base.Address
	contract         base.Address
	projectID        string
	requestTimeStamp uint64
	data             string
	currency         ctypes.CurrencyID
}

func NewIssueFact(
	token []byte, sender, contract base.Address, projectID string,
	requestTimeStamp uint64, data string, currency ctypes.CurrencyID) IssueFact {
	bf := base.NewBaseFact(IssueFactHint, token)
	fact := IssueFact{
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

func (fact IssueFact) IsValid(b []byte) error {
	if len(fact.projectID) < 1 || len(fact.projectID) > types.MaxProjectIDLen {
		return common.ErrFactInvalid.Wrap(
			common.ErrValOOR.Wrap(
				errors.Errorf(
					"invalid projectID length %v < 1 or > %v", len(fact.projectID), types.MaxProjectIDLen)))
	}

	if !ctypes.ReValidSpcecialCh.Match([]byte(fact.projectID)) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Wrap(
				errors.Errorf("projectID ID %s, must match regex `^[^\\s:/?#\\[\\]$@]*$`", fact.projectID)))
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

func (fact IssueFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact IssueFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact IssueFact) Bytes() []byte {
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

func (fact IssueFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact IssueFact) Sender() base.Address {
	return fact.sender
}

func (fact IssueFact) Contract() base.Address {
	return fact.contract
}

func (fact IssueFact) ProjectId() string {
	return fact.projectID
}

func (fact IssueFact) RequestTimeStamp() uint64 {
	return fact.requestTimeStamp
}

func (fact IssueFact) Data() string {
	return fact.data
}

func (fact IssueFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact IssueFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender}, nil
}

func (fact IssueFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact IssueFact) FeePayer() base.Address {
	return fact.sender
}

func (fact IssueFact) FactUser() base.Address {
	return fact.sender
}

func (fact IssueFact) Signer() base.Address {
	return fact.sender
}

func (fact IssueFact) ActiveContractOwnerHandlerOnly() [][2]base.Address {
	return [][2]base.Address{{fact.contract, fact.sender}}
}

type Issue struct {
	extras.ExtendedOperation
}

func NewIssue(fact IssueFact) (Issue, error) {
	return Issue{
		ExtendedOperation: extras.NewExtendedOperation(IssueHint, fact),
	}, nil
}
