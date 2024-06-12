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

type TestCreateServiceProcessor struct {
	*test.BaseTestOperationProcessorNoItem[RegisterModel]
}

func NewTestCreateServiceProcessor(tp *test.TestProcessor) TestCreateServiceProcessor {
	t := test.NewBaseTestOperationProcessorNoItem[RegisterModel](tp)
	return TestCreateServiceProcessor{&t}
}

func (t *TestCreateServiceProcessor) Create() *TestCreateServiceProcessor {
	t.Opr, _ = NewRegisterModelProcessor()(
		base.GenesisHeight,
		t.GetStateFunc,
		nil, nil,
	)
	return t
}

func (t *TestCreateServiceProcessor) SetCurrency(
	cid string, am int64, addr base.Address, target []currencytypes.CurrencyID, instate bool,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.SetCurrency(cid, am, addr, target, instate)

	return t
}

func (t *TestCreateServiceProcessor) SetAmount(
	am int64, cid currencytypes.CurrencyID, target []currencytypes.Amount,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.SetAmount(am, cid, target)

	return t
}

func (t *TestCreateServiceProcessor) SetContractAccount(
	owner base.Address, priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.SetContractAccount(owner, priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateServiceProcessor) SetAccount(
	priv string, amount int64, cid currencytypes.CurrencyID, target []test.Account, inState bool,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.SetAccount(priv, amount, cid, target, inState)

	return t
}

func (t *TestCreateServiceProcessor) LoadOperation(fileName string,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.LoadOperation(fileName)

	return t
}

func (t *TestCreateServiceProcessor) Print(fileName string,
) *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.Print(fileName)

	return t
}

func (t *TestCreateServiceProcessor) SetService(
	contract base.Address,
) *TestCreateServiceProcessor {
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

func (t *TestCreateServiceProcessor) MakeOperation(
	sender base.Address, privatekey base.Privatekey, contract base.Address, currency currencytypes.CurrencyID,
) *TestCreateServiceProcessor {
	op, _ := NewRegisterModel(
		NewRegisterModelFact(
			[]byte("token"),
			sender,
			contract,
			currency,
		))
	_ = op.Sign(privatekey, t.NetworkID)
	t.Op = op

	return t
}

func (t *TestCreateServiceProcessor) RunPreProcess() *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.RunPreProcess()

	return t
}

func (t *TestCreateServiceProcessor) RunProcess() *TestCreateServiceProcessor {
	t.BaseTestOperationProcessorNoItem.RunProcess()

	return t
}
