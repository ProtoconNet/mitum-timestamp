package cmds

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	digestisaac "github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/operation/timestamp"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: common.BaseStateHint, Instance: common.BaseState{}},
	{Hint: common.NodeHint, Instance: common.BaseNode{}},

	{Hint: currencytypes.AddressHint, Instance: currencytypes.Address{}},
	{Hint: currencytypes.AmountHint, Instance: currencytypes.Amount{}},
	{Hint: currencytypes.AccountHint, Instance: currencytypes.Account{}},
	{Hint: currencytypes.AccountKeysHint, Instance: currencytypes.BaseAccountKeys{}},
	{Hint: currencytypes.AccountKeyHint, Instance: currencytypes.BaseAccountKey{}},
	{Hint: currencytypes.ContractAccountKeysHint, Instance: currencytypes.ContractAccountKeys{}},
	{Hint: currencytypes.NilFeeerHint, Instance: currencytypes.NilFeeer{}},
	{Hint: currencytypes.FixedFeeerHint, Instance: currencytypes.FixedFeeer{}},
	{Hint: currencytypes.RatioFeeerHint, Instance: currencytypes.RatioFeeer{}},
	{Hint: currencytypes.CurrencyPolicyHint, Instance: currencytypes.CurrencyPolicy{}},
	{Hint: currencytypes.CurrencyDesignHint, Instance: currencytypes.CurrencyDesign{}},

	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.TimeStampItemHint, Instance: types.TimeStampItem{}},

	{Hint: currency.CreateAccountsItemMultiAmountsHint, Instance: currency.CreateAccountsItemMultiAmounts{}},
	{Hint: currency.CreateAccountsItemSingleAmountHint, Instance: currency.CreateAccountsItemSingleAmount{}},
	{Hint: currency.CreateAccountsHint, Instance: currency.CreateAccounts{}},
	{Hint: currency.KeyUpdaterHint, Instance: currency.KeyUpdater{}},
	{Hint: currency.TransfersItemMultiAmountsHint, Instance: currency.TransfersItemMultiAmounts{}},
	{Hint: currency.TransfersItemSingleAmountHint, Instance: currency.TransfersItemSingleAmount{}},
	{Hint: currency.TransfersHint, Instance: currency.Transfers{}},
	{Hint: currency.CurrencyRegisterHint, Instance: currency.CurrencyRegister{}},
	{Hint: currency.CurrencyPolicyUpdaterHint, Instance: currency.CurrencyPolicyUpdater{}},
	{Hint: currency.SuffrageInflationHint, Instance: currency.SuffrageInflation{}},
	{Hint: currency.GenesisCurrenciesHint, Instance: currency.GenesisCurrencies{}},
	{Hint: currency.GenesisCurrenciesFactHint, Instance: currency.GenesisCurrenciesFact{}},

	{Hint: extension.CreateContractAccountsItemMultiAmountsHint, Instance: extension.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extension.CreateContractAccountsItemSingleAmountHint, Instance: extension.CreateContractAccountsItemSingleAmount{}},
	{Hint: extension.CreateContractAccountsHint, Instance: extension.CreateContractAccounts{}},
	{Hint: extension.WithdrawsItemMultiAmountsHint, Instance: extension.WithdrawsItemMultiAmounts{}},
	{Hint: extension.WithdrawsItemSingleAmountHint, Instance: extension.WithdrawsItemSingleAmount{}},
	{Hint: extension.WithdrawsHint, Instance: extension.Withdraws{}},

	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},

	{Hint: timestamp.AppendHint, Instance: timestamp.Append{}},
	{Hint: timestamp.ServiceRegisterHint, Instance: timestamp.ServiceRegister{}},

	{Hint: state.TimeStampItemStateValueHint, Instance: state.TimeStampItemStateValue{}},
	{Hint: state.ServiceDesignStateValueHint, Instance: state.ServiceDesignStateValue{}},
	{Hint: state.TimeStampLastIndexStateValueHint, Instance: state.TimeStampLastIndexStateValue{}},

	{Hint: statecurrency.AccountStateValueHint, Instance: statecurrency.AccountStateValue{}},
	{Hint: statecurrency.BalanceStateValueHint, Instance: statecurrency.BalanceStateValue{}},
	{Hint: statecurrency.CurrencyDesignStateValueHint, Instance: statecurrency.CurrencyDesignStateValue{}},
	{Hint: stateextension.ContractAccountStateValueHint, Instance: stateextension.ContractAccountStateValue{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: currencydigest.AccountValueHint, Instance: currencydigest.AccountValue{}},
	{Hint: currencydigest.OperationValueHint, Instance: currencydigest.OperationValue{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: currency.CreateAccountsFactHint, Instance: currency.CreateAccountsFact{}},
	{Hint: currency.KeyUpdaterFactHint, Instance: currency.KeyUpdaterFact{}},
	{Hint: currency.TransfersFactHint, Instance: currency.TransfersFact{}},
	{Hint: currency.SuffrageInflationFactHint, Instance: currency.SuffrageInflationFact{}},
	{Hint: currency.CurrencyRegisterFactHint, Instance: currency.CurrencyRegisterFact{}},
	{Hint: currency.CurrencyPolicyUpdaterFactHint, Instance: currency.CurrencyPolicyUpdaterFact{}},

	{Hint: extension.CreateContractAccountsFactHint, Instance: extension.CreateContractAccountsFact{}},
	{Hint: extension.WithdrawsFactHint, Instance: extension.WithdrawsFact{}},

	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},

	{Hint: timestamp.AppendFactHint, Instance: timestamp.AppendFact{}},
	{Hint: timestamp.ServiceRegisterFactHint, Instance: timestamp.ServiceRegisterFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
