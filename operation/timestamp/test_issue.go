package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/test"
	"github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statetimestamp "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type TestIssueProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Issue]
}

func NewTestIssueProcessor(tp *test.TestProcessor) TestIssueProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Issue](tp)
	return TestIssueProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestIssueProcessor) Create() *TestIssueProcessor {
	t.Opr, _ = NewIssueProcessor()(
		base.GenesisHeight,
		nil,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestIssueProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []currencytypes.CurrencyID, instate bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestIssueProcessor) SetAmount(
	am int64, cid currencytypes.CurrencyID, target []currencytypes.Amount,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestIssueProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestIssueProcessor) SetAccount(
	priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestIssueProcessor) LoadOperation(fileName string,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestIssueProcessor) Print(fileName string,
) *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestIssueProcessor) RunPreProcess() *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestIssueProcessor) RunProcess() *TestIssueProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestIssueProcessor) SetService(
	contract base.Address,
) *TestIssueProcessor {
	pids := []string(nil)
	design := types.NewDesign(pids...)

	st := common.NewBaseState(base.Height(1), statetimestamp.DesignStateKey(contract), statetimestamp.NewDesignStateValue(design), nil, []util.Hash{})
	t.SetState(st, true)

	cst, found, _ := t.MockGetter.Get(extension.StateKeyContractAccount(contract))
	if !found {
		panic("contract account not set")
	}
	status, err := extension.StateContractAccountValue(cst)
	if err != nil {
		panic(err)
	}

	nstatus := status.SetIsActive(true)
	cState := common.NewBaseState(base.Height(1), extension.StateKeyContractAccount(contract), extension.NewContractAccountStateValue(nstatus), nil, []util.Hash{})
	t.SetState(cState, true)

	return t
}

func (t *TestIssueProcessor) MakeOperation(
	sender base.Address,
	privatekey base.Privatekey,
	contract base.Address,
	projectID string,
	requestTimeStamp uint64,
	data string,
	currency currencytypes.CurrencyID,
) *TestIssueProcessor {
	op, _ := NewIssue(
		NewIssueFact(
			[]byte("token"),
			sender,
			contract,
			projectID,
			requestTimeStamp,
			data,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}
