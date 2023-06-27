package cmds

type TimestampCommand struct {
	Append          AppendCommand          `cmd:"" name:"ts-append" help:"append new timestamp item"`
	ServiceRegister ServiceRegisterCommand `cmd:"" name:"service-register" help:"register timestamp service"`
}
