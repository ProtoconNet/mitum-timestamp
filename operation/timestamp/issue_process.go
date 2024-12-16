package timestamp

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var issueProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(IssueProcessor)
	},
}

func (Issue) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type IssueProcessor struct {
	*base.BaseOperationProcessor
	proposal *base.ProposalSignFact
}

func NewIssueProcessor() ctypes.GetNewProcessorWithProposal {
	return func(
		height base.Height,
		proposal *base.ProposalSignFact,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new IssueProcessor")

		nopp := issueProcessorPool.Get()
		opp, ok := nopp.(*IssueProcessor)
		if !ok {
			return nil, e.Errorf("expected IssueProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b
		opp.proposal = proposal

		return opp, nil
	}
}

func (opp *IssueProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(IssueFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", IssueFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := cstate.CheckExistsState(state.DesignStateKey(fact.Contract()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("timestamp service for contract account %v",
				fact.Contract(),
			)), nil
	}

	k := state.LastIdxStateKey(fact.Contract(), fact.ProjectId())
	switch _, _, err := getStateFunc(k); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError("getting timestamp item last index failed, %q; %w", fact.Contract(), err), nil
	}

	return ctx, nil, nil
}

func (opp *IssueProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(IssueFact)

	st, err := cstate.ExistsState(state.DesignStateKey(fact.Contract()), "service design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design not found, %q; %w", fact.Contract(), err), nil
	}

	design, err := state.GetDesignFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design value not found, %q; %w", fact.Contract(), err), nil
	}

	design.AddProject(fact.ProjectId())
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
	}

	var idx uint64
	k := state.LastIdxStateKey(fact.Contract(), fact.ProjectId())
	switch st, found, err := getStateFunc(k); {
	case err != nil:
		return nil, base.NewBaseOperationProcessReasonError(
			"getting timestamp item lastindex failed, %q; %w",
			fact.Contract(),
			err,
		), nil
	case found:
		i, err := state.GetLastIdxFromState(st)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"getting timestamp item lastindex value failed, %q; %w",
				fact.Contract(),
				err,
			), nil
		}
		idx = i + 1
	case !found:
		idx = 0
		st = base.NewBaseState(base.NilHeight, k, nil, nil, nil)
	}

	proposal := *opp.proposal
	nowTime := uint64(proposal.ProposalFact().ProposedAt().Unix())

	tsItem := types.NewItem(
		fact.ProjectId(),
		fact.RequestTimeStamp(),
		nowTime,
		idx,
		fact.Data(),
	)
	if err := tsItem.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid timestamp; %w", err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc
	sts = append(sts, cstate.NewStateMergeValue(
		state.ItemStateKey(fact.Contract(), fact.ProjectId(), idx),
		state.NewItemStateValue(tsItem),
	))
	sts = append(sts, cstate.NewStateMergeValue(
		state.LastIdxStateKey(fact.Contract(), fact.ProjectId()),
		state.NewLastIdxStateValue(fact.ProjectId(), idx),
	))
	sts = append(sts, cstate.NewStateMergeValue(
		state.DesignStateKey(fact.Contract()),
		state.NewDesignStateValue(design),
	))

	return sts, nil, nil
}

func (opp *IssueProcessor) Close() error {
	opp.proposal = nil
	issueProcessorPool.Put(opp)

	return nil
}
