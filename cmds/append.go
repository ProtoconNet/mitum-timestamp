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

type AppendCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender           currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Target           currencycmds.AddressFlag    `arg:"" name:"target" help:"target address" required:"true"`
	Service          string                      `arg:"" name:"service" help:"service id" required:"true"`
	ProjectID        string                      `arg:"" name:"project id" help:"project id" required:"true"`
	RequestTimeStamp uint64                      `arg:"" name:"request timestamp" help:"request timestamp" required:"true"`
	Data             string                      `arg:"" name:"data" help:"data" required:"true"`
	Currency         currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender           base.Address
	target           base.Address
	service          currencytypes.ContractID
}

func NewAppendCommand() AppendCommand {
	cmd := NewBaseCommand()
	return AppendCommand{BaseCommand: *cmd}
}

func (cmd *AppendCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *AppendCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Target.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid target format, %q", cmd.Target)
	} else {
		cmd.target = a
	}

	service := currencytypes.ContractID(cmd.Service)
	if err := service.IsValid(nil); err != nil {
		return err
	} else {
		cmd.service = service
	}

	if len(cmd.ProjectID) < 1 {
		return errors.Errorf("invalid ProjectID, %s", cmd.ProjectID)
	}

	if len(cmd.Data) < 1 {
		return errors.Errorf("invalid data, %s", cmd.ProjectID)
	}

	if cmd.RequestTimeStamp < 1 {
		return errors.Errorf("invalid Request timestamp, %s", cmd.RequestTimeStamp)
	}

	return nil
}

func (cmd *AppendCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create append operation")

	fact := timestampservice.NewAppendFact([]byte(cmd.Token), cmd.sender, cmd.target, cmd.service, cmd.ProjectID, cmd.RequestTimeStamp, cmd.Data, cmd.Currency.CID)

	op, err := timestampservice.NewAppend(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
