package postgres

import (
	"errors"
	"fmt"
)

const (
	DbConnectionFailedMsg = "db connection failed"
	SettingNodeFailedMsg  = "unable to set db node"
)

func ErrDBConnectionFailed(connectErr error) error {
	return formatError(DbConnectionFailedMsg, connectErr.Error())
}

func ErrUnableToSetNode(setErr error) error {
	return formatError(SettingNodeFailedMsg, setErr.Error())
}

func formatError(msg, err string) error {
	return errors.New(fmt.Sprintf("%s: %s", msg, err))
}
