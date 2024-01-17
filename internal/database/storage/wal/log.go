package wal

import (
	"kv_db/internal/database/comd"
	"kv_db/pkg/dfuture"
	"kv_db/pkg/dpromise"
)

type LogData struct {
	LSN       int64
	CommandID comd.CmdID
	Arguments []string
}

type Log struct {
	data         LogData
	writePromise dpromise.Promise[error]
}

func NewLog(lsn int64, commandID comd.CmdID, args []string) Log {
	return Log{
		data: LogData{
			LSN:       lsn,
			CommandID: commandID,
			Arguments: args,
		},
		writePromise: dpromise.NewPromise[error](),
	}
}

func (l *Log) Data() LogData {
	return l.data
}

func (l *Log) LSN() int64 {
	return l.data.LSN
}

func (l *Log) CommandID() comd.CmdID {
	return l.data.CommandID
}

func (l *Log) Arguments() []string {
	return l.data.Arguments
}

func (l *Log) SetResult(err error) {
	l.writePromise.Set(err)
}

func (l *Log) Result() dfuture.Future[error] {
	return l.writePromise.GetFuture()
}
