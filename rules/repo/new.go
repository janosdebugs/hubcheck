package repo

import (
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/repo/actionspermissions"
)

func New() []hubcheck.RepoRule {
	return []hubcheck.RepoRule{
		actionspermissions.New(),
	}
}
