package database

import "kv_db/internal/database/comd"

type Query struct {
	commandID comd.CmdID
	arguments []string
}

func NewQuery(commandID comd.CmdID, arguments []string) Query {
	return Query{
		commandID: commandID,
		arguments: arguments,
	}
}

func (c *Query) CommandID() comd.CmdID {
	return c.commandID
}

func (c *Query) Arguments() []string {
	return c.arguments
}
