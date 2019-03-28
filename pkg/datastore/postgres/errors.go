package postgres

import (
	"errors"
	"fmt"
)

const (
	BeginTransactionFailedMsg = "failed to begin transaction"
	DbConnectionFailedMsg     = "db connection failed"
	DeleteQueryFailedMsg      = "delete query failed"
	InsertQueryFailedMsg      = "insert query failed"
	SettingNodeFailedMsg      = "unable to set db node"
)

func ErrBeginTransactionFailed(beginErr error) error {
	return formatError(BeginTransactionFailedMsg, beginErr.Error())
}

func ErrDBConnectionFailed(connectErr error) error {
	return formatError(DbConnectionFailedMsg, connectErr.Error())
}

func ErrDBDeleteFailed(deleteErr error) error {
	return formatError(DeleteQueryFailedMsg, deleteErr.Error())
}

func ErrDBInsertFailed(insertErr error) error {
	return formatError(InsertQueryFailedMsg, insertErr.Error())
}

func ErrUnableToSetNode(setErr error) error {
	return formatError(SettingNodeFailedMsg, setErr.Error())
}

func formatError(msg, err string) error {
	return errors.New(fmt.Sprintf("%s: %s", msg, err))
}
