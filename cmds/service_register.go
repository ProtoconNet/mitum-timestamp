package cmds

import (
	"context"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	timestampservice "github.com/ProtoconNet/mitum-timestamp/operation/timestamp"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type ServiceRegisterCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Target   currencycmds.AddressFlag    `arg:"" name:"target" help:"target account to register policy" required:"true"`
	Service  string                      `arg:"" name:"service" help:"STO ID" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	target   base.Address
	service  currencytypes.ContractID
}

func (cmd *ServiceRegisterCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ServiceRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Target.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid target format; %q", cmd.Target)
	} else {
		cmd.target = a
	}

	service := currencytypes.ContractID(cmd.Service)
	if err := service.IsValid(nil); err != nil {
		return err
	} else {
		cmd.service = service
	}

	return nil
}

func (cmd *ServiceRegisterCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create service-register operation")

	fact := timestampservice.NewServiceRegisterFact([]byte(cmd.Token), cmd.sender, cmd.target, cmd.service, cmd.Currency.CID)

	op, err := timestampservice.NewServiceRegister(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
