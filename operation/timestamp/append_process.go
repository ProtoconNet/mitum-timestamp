package timestamp

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	statetimestamp "github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/pkg/errors"
	"sync"

	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var appendProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(AppendProcessor)
	},
}

func (Append) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type AppendProcessor struct {
	*mitumbase.BaseOperationProcessor
	getLastBlockFunc processor.GetLastBlockFunc
}

func NewAppendProcessor(getLastBlockFunc processor.GetLastBlockFunc) currencytypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new AppendProcessor")

		nopp := appendProcessorPool.Get()
		opp, ok := nopp.(*AppendProcessor)
		if !ok {
			return nil, e.Errorf("expected AppendProcessor, not %T", nopp)
		}

		b, err := mitumbase.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b
		opp.getLastBlockFunc = getLastBlockFunc

		return opp, nil
	}
}

func (opp *AppendProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	e := util.StringError("failed to preprocess Mint")

	fact, ok := op.Fact().(AppendFact)
	if !ok {
		return ctx, nil, e.Errorf("expected AppendFact, not %T", op.Fact())
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
		return ctx, mitumbase.NewBaseOperationProcessReasonError("contract account cannot Append timestamp, %q", fact.Sender()), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError("invalid signing; %w", err), nil
	}

	st, err := state.ExistsState(stateextension.StateKeyContractAccount(fact.Target()), "key of contract account", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("target contract account state not found, %q; %w", fact.Target(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("contract account value not found from state, %q; %w", fact.Target(), err), nil
	}

	if !(ca.Owner().Equal(fact.Sender()) || ca.IsOperator(fact.Sender())) {
		return nil, mitumbase.NewBaseOperationProcessReasonError("sender is neither the owner nor the operator of the target contract account, %q", fact.Sender()), nil
	}

	_, err = state.ExistsState(statetimestamp.StateKeyServiceDesign(fact.Target()), "key of service design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("service design not found, %q; %w", fact.Target(), err), nil
	}

	k := statetimestamp.StateKeyTimeStampLastIndex(fact.Target(), fact.ProjectId())
	switch _, _, err := getStateFunc(k); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError("getting timestamp item lastindex failed, %q; %w", fact.Target(), err), nil
	}

	_, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("get LastBlock failed; %w", err), nil
	} else if !found {
		return nil, mitumbase.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	return ctx, nil, nil
}

func (opp *AppendProcessor) Process( // nolint:dupl
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	e := util.StringError("failed to process Append")

	fact, ok := op.Fact().(AppendFact)
	if !ok {
		return nil, nil, e.Errorf("expected AppendFact, not %T", op.Fact())
	}

	st, err := state.ExistsState(statetimestamp.StateKeyServiceDesign(fact.Target()), "key of service design", getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("service design not found, %q; %w", fact.Target(), err), nil
	}

	design, err := statetimestamp.StateServiceDesignValue(st)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("service design value not found, %q; %w", fact.Target(), err), nil
	}

	design.AddProject(fact.ProjectId())

	var idx uint64
	k := statetimestamp.StateKeyTimeStampLastIndex(fact.Target(), fact.ProjectId())
	switch st, found, err := getStateFunc(k); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"getting timestamp item lastindex failed, %q; %w",
			fact.Target(),
			err,
		), nil
	case found:
		i, err := statetimestamp.StateTimeStampLastIndexValue(st)
		if err != nil {
			return nil, mitumbase.NewBaseOperationProcessReasonError(
				"getting timestamp item lastindex value failed, %q; %w",
				fact.Target(),
				err,
			), nil
		}
		idx = i + 1
	case !found:
		idx = 0
		st = mitumbase.NewBaseState(mitumbase.NilHeight, k, nil, nil, nil)
	}

	blockmap, found, err := opp.getLastBlockFunc()
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("get LastBlock failed; %w", err), nil
	} else if !found {
		return nil, mitumbase.NewBaseOperationProcessReasonError("LastBlock not found"), nil
	}

	tsItem := types.NewTimeStampItem(
		fact.ProjectId(),
		fact.RequestTimeStamp(),
		uint64(blockmap.Manifest().ProposedAt().Unix()),
		idx,
		fact.Data(),
	)
	if err := tsItem.IsValid(nil); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("invalid timestamp; %w", err), nil
	}

	var sts []mitumbase.StateMergeValue // nolint:prealloc
	sts = append(sts, state.NewStateMergeValue(
		statetimestamp.StateKeyTimeStampItem(fact.Target(), fact.ProjectId(), idx),
		statetimestamp.NewTimeStampItemStateValue(tsItem),
	))
	sts = append(sts, state.NewStateMergeValue(
		statetimestamp.StateKeyTimeStampLastIndex(fact.Target(), fact.ProjectId()),
		statetimestamp.NewTimeStampLastIndexStateValue(fact.ProjectId(), idx),
	))

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

	senderBalSt, err := state.ExistsState(
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

	switch senderBal, err := statecurrency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if currencyPolicy.Feeer().Receiver() != nil {
		if err := state.CheckExistsState(statecurrency.StateKeyAccount(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
			return nil, nil, err
		} else if feeRcvrSt, found, err := getStateFunc(statecurrency.StateKeyBalance(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
			return nil, nil, err
		} else if !found {
			return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
		} else if feeRcvrSt.Key() != senderBalSt.Key() {
			r, ok := feeRcvrSt.Value().(statecurrency.BalanceStateValue)
			if !ok {
				return nil, nil, errors.Errorf("expected %T, not %T", statecurrency.BalanceStateValue{}, feeRcvrSt.Value())
			}
			sts = append(sts, common.NewBaseStateMergeValue(
				feeRcvrSt.Key(),
				statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
				func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
				},
			))

			sts = append(sts, common.NewBaseStateMergeValue(
				senderBalSt.Key(),
				statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
				func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
					return statecurrency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
				},
			))
		}
	}

	return sts, nil, nil
}

func (opp *AppendProcessor) Close() error {
	appendProcessorPool.Put(opp)

	return nil
}
