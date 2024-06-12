package cmds

type TimestampCommand struct {
	Issue         IssueCommand         `cmd:"" name:"issue" help:"issue new timestamp item"`
	RegisterModel RegisterModelCommand `cmd:"" name:"register-model" help:"register timestamp model"`
}
