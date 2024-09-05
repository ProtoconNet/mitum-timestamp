package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-timestamp/operation/timestamp"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-timestamp/operation/processor"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-timestamp-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *currencyprocessor.OperationProcessor
	var setA *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]
	var setB *hint.CompatibleSet[currencycmds.NewOperationProcessorInternalWithProposalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		currencycmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &setA,
		currencycmds.OperationProcessorsMapBContextKey, &setB,
	); err != nil {
		return pctx, err
	}

	//err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	//if err != nil {
	//	return pctx, err
	//}
	err := opr.SetGetNewProcessorFunc(processor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	if err := opr.SetProcessor(
		timestamp.RegisterModelHint,
		timestamp.NewRegisterModelProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessorWithProposal(
		timestamp.IssueHint,
		timestamp.NewIssueProcessor(),
	); err != nil {
		return pctx, err
	}

	_ = setB.Add(timestamp.IssueHint, func(height base.Height, proposal base.ProposalSignFact, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		if err := opr.SetProposal(&proposal); err != nil {
			return nil, err
		}
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	_ = setA.Add(timestamp.RegisterModelHint, func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	pctx = context.WithValue(pctx, currencycmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, setA)        //revive:disable-line:modifies-parameter
	pctx = context.WithValue(pctx, currencycmds.OperationProcessorsMapBContextKey, setB) //revive:disable-line:modifies-parameter

	return pctx, nil
}
