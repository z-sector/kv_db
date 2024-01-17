package wal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/internal/database/comd"
	"kv_db/pkg/dlog"
)

func TestReadLogs(t *testing.T) {
	t.Parallel()

	reader := NewFSReader("test_data", dlog.NewNonSlog())

	logs, err := reader.ReadLogs()
	require.NoError(t, err)
	require.Equal(t, 9, len(logs))

	// from tests_data/wal_1000.log
	require.Equal(t, LogData{LSN: 1, CommandID: comd.SetCommandID, Arguments: []string{"key_1", "value_1"}}, logs[0])
	require.Equal(t, LogData{LSN: 2, CommandID: comd.SetCommandID, Arguments: []string{"key_2", "value_2"}}, logs[1])
	require.Equal(t, LogData{LSN: 3, CommandID: comd.SetCommandID, Arguments: []string{"key_3", "value_3"}}, logs[2])

	// from tests_data/wal_2000.log
	require.Equal(t, LogData{LSN: 4, CommandID: comd.SetCommandID, Arguments: []string{"key_4", "value_4"}}, logs[3])
	require.Equal(t, LogData{LSN: 5, CommandID: comd.SetCommandID, Arguments: []string{"key_5", "value_5"}}, logs[4])
	require.Equal(t, LogData{LSN: 6, CommandID: comd.SetCommandID, Arguments: []string{"key_6", "value_6"}}, logs[5])

	// from tests_data/wal_3000.log
	require.Equal(t, LogData{LSN: 7, CommandID: comd.SetCommandID, Arguments: []string{"key_7", "value_7"}}, logs[6])
	require.Equal(t, LogData{LSN: 8, CommandID: comd.SetCommandID, Arguments: []string{"key_8", "value_8"}}, logs[7])
	require.Equal(t, LogData{LSN: 9, CommandID: comd.SetCommandID, Arguments: []string{"key_9", "value_9"}}, logs[8])
}
