package database

type CmdID int

const (
	UnknownCommandID CmdID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var (
	UnknownCommand = "UNKNOWN"
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
)

var commandNameToID = map[string]CmdID{
	UnknownCommand: UnknownCommandID,
	SetCommand:     SetCommandID,
	GetCommand:     GetCommandID,
	DelCommand:     DelCommandID,
}

func GetCommandIDByName(command string) CmdID {
	status, ok := commandNameToID[command]
	if !ok {
		return UnknownCommandID
	}

	return status
}
