package database

type Query struct {
	commandID CmdID
	arguments []string
}

func NewQuery(commandID CmdID, arguments []string) Query {
	return Query{
		commandID: commandID,
		arguments: arguments,
	}
}

func (c *Query) CommandID() CmdID {
	return c.commandID
}

func (c *Query) Arguments() []string {
	return c.arguments
}
