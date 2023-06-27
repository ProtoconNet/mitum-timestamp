package timestamp

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statetimestamp "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"sync"

	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var serviceRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ServiceRegisterProcessor)
	},
}

func (ServiceRegister) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ServiceRegisterProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewServiceRegisterProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new ServiceRegisterProcessor")

		nopp := serviceRegisterProcessorPool.Get()
		opp, ok := nopp.(*ServiceRegisterProcessor)
		if !ok {
			return nil, errors.Errorf("expected servicesRegisterProcessor, not %T", nopp)
		}

		b, err := mitumbase.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *ServiceRegisterProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess serviceRegister")

	fact, ok := op.Fact().(ServiceRegisterFact)
	if !ok {
		return ctx, nil, e.Errorf("expected ServiceRegisterFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q: %w", fact.Sender(), err), nil
	}

	_, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender address is contract account, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing: %w", err), nil
	}

	st, err := state.ExistsState(stateextension.StateKeyContractAccount(fact.Target()), "key of contract account", getStateFunc)
	if err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("target contract account not found, %q: %w", fact.Target(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q: %w", fact.Target(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender is not owner of contract account, %q, %q", fact.Sender(), ca.Owner()), nil
	}

	if !ca.IsActive() {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("deactivated contract account, %q", fact.Target()), nil
	}

	if err := state.CheckNotExistsState(statetimestamp.StateKeyServiceDesign(fact.Target(), fact.Service()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("service design already exists, %q: %w", fact.Service(), err), nil
	}

	return ctx, nil, nil
}

func (opp *ServiceRegisterProcessor) Process(
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process ServiceRegister")

	fact, ok := op.Fact().(ServiceRegisterFact)
	if !ok {
		return nil, nil, e.Errorf("expected ServiceRegisterFact, not %T", op.Fact())
	}

	sts := make([]mitumbase.StateMergeValue, 2)
	pids := []string{}

	design := types.NewDesign(fact.Service(), pids...)
	if err := design.IsValid(nil); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("invalid service design, %q: %w", fact.Service(), err), nil
	}

	sts[0] = statetimestamp.NewStateMergeValue(
		statetimestamp.StateKeyServiceDesign(fact.target, design.Service()),
		statetimestamp.NewServiceDesignStateValue(design),
	)

	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q: %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to check fee of currency, %q: %w", fact.Currency(), err), nil
	}

	st, err := state.ExistsState(statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of sender balance", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender balance not found, %q: %w", fact.Sender(), err), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := statecurrency.StateBalanceValue(st); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get balance value, %q: %w", statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError("not enough balance of sender, %q", fact.Sender()), nil
	}

	v, ok := sb.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[1] = state.NewStateMergeValue(
		sb.Key(),
		statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *ServiceRegisterProcessor) Close() error {
	serviceRegisterProcessorPool.Put(opp)

	return nil
}
