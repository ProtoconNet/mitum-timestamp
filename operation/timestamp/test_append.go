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

type TestAppendProcessor struct {
	*test.BaseTestOperationProcessorNoItem[Append]
}

func NewTestAppendProcessor(tp *test.TestProcessor) TestAppendProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[Append](tp)
	return TestAppendProcessor{BaseTestOperationProcessorNoItem: &t}
}

func (t *TestAppendProcessor) Create() *TestAppendProcessor {
	t.Opr, _ = NewAppendProcessor(func() (base.BlockMap, bool, error) { return nil, true, nil })(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestAppendProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []currencytypes.CurrencyID, instate bool,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestAppendProcessor) SetAmount(
	am int64, cid currencytypes.CurrencyID, target []currencytypes.Amount,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestAppendProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestAppendProcessor) SetAccount(
	priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestAppendProcessor) LoadOperation(fileName string,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestAppendProcessor) Print(fileName string,
) *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestAppendProcessor) RunPreProcess() *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestAppendProcessor) RunProcess() *TestAppendProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}

func (t *TestAppendProcessor) SetService(
	contract base.Address,
) *TestAppendProcessor {
	pids := []string(nil)
	design := types.NewDesign(pids...)

	st := common.NewBaseState(base.Height(1), statetimestamp.StateKeyServiceDesign(contract), statetimestamp.NewServiceDesignStateValue(design), nil, []util.Hash{})
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

func (t *TestAppendProcessor) MakeOperation(
	sender base.Address,
	privatekey base.Privatekey,
	contract base.Address,
	projectID string,
	requestTimeStamp uint64,
	data string,
	currency currencytypes.CurrencyID,
) *TestAppendProcessor {
	op, _ := NewAppend(
		NewAppendFact(
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
