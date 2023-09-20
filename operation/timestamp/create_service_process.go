package timestamp

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statetimestamp "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"

	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var createServiceProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateServiceProcessor)
	},
}

func (CreateService) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateServiceProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewCreateServiceProcessor() currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateServiceProcessor")

		nopp := createServiceProcessorPool.Get()
		opp, ok := nopp.(*CreateServiceProcessor)
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

func (opp *CreateServiceProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return ctx, nil, e.Errorf("expected CreateServiceFact, not %T", op.Fact())
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := state.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("sender not found, %q; %w", fact.Sender(), err), nil
	}

	_, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	if err := state.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			"sender address is contract account, %q",
			fact.Sender(),
		), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	st, err := state.ExistsState(
		stateextension.StateKeyContractAccount(fact.Target()),
		"key of contract account",
		getStateFunc,
	)
	if err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			"target contract account not found, %q; %w",
			fact.Target(),
			err,
		), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get state value of contract account, %q; %w",
			fact.Target(),
			err,
		), nil
	}

	if !(ca.Owner().Equal(fact.Sender()) || ca.IsOperator(fact.Sender())) {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"sender is neither the owner nor the operator of the target contract account, %q",
			fact.Sender(),
		), nil
	}

	if ca.IsActive() {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"a design is already registered, %q",
			fact.Target().String(),
		), nil
	}

	if err := state.CheckNotExistsState(statetimestamp.StateKeyServiceDesign(fact.Target()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			"service design already exists, %q; %w",
			fact.Target(),
			err,
		), nil
	}

	return ctx, nil, nil
}

func (opp *CreateServiceProcessor) Process(
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process CreateService")

	fact, ok := op.Fact().(CreateServiceFact)
	if !ok {
		return nil, nil, e.Errorf("expected CreateServiceFact, not %T", op.Fact())
	}

	sts := make([]mitumbase.StateMergeValue, 3)
	pids := []string(nil)

	design := types.NewDesign(pids...)
	if err := design.IsValid(nil); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Target(), err), nil
	}

	sts[0] = statetimestamp.NewStateMergeValue(
		statetimestamp.StateKeyServiceDesign(fact.Target()),
		statetimestamp.NewServiceDesignStateValue(design),
	)

	st, err := state.ExistsState(stateextension.StateKeyContractAccount(fact.Target()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("target contract account not found, %q; %w", fact.Target(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("failed to get state value of contract account, %q; %w", fact.Target(), err), nil
	}
	ca.SetIsActive(true)

	sts[1] = state.NewStateMergeValue(
		stateextension.StateKeyContractAccount(fact.Target()),
		stateextension.NewContractAccountStateValue(ca),
	)

	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %q; %w", fact.Currency(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	st, err = state.ExistsState(
		statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
		"key of sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"sender balance not found, %q; %w",
			fact.Sender(),
			err,
		), nil
	}
	sb := state.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := statecurrency.StateBalanceValue(st); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case b.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := sb.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", sb.Value()), nil
	}
	sts[2] = state.NewStateMergeValue(
		sb.Key(),
		statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CreateServiceProcessor) Close() error {
	createServiceProcessorPool.Put(opp)

	return nil
}
