package repo

import (
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/repo/actionspermissions"
	"go.debugged.it/hubcheck/rules/repo/vulnalerts"
)

func New() []hubcheck.RepoRule {
	return []hubcheck.RepoRule{
		actionspermissions.New(),
		vulnalerts.New(),
	}
}
