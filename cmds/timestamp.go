package cmds

type TimestampCommand struct {
	Append        AppendCommand        `cmd:"" name:"ts-append" help:"append new timestamp item"`
	CreateService CreateServiceCommand `cmd:"" name:"creates-service" help:"register timestamp service"`
}
