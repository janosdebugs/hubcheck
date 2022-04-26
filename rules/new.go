package rules

import (
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/actionspermissions"
	"go.debugged.it/hubcheck/rules/defaultrepopermission"
	"go.debugged.it/hubcheck/rules/orgadmins"
	"go.debugged.it/hubcheck/rules/twofactor"
	"go.debugged.it/hubcheck/rules/workflowapprovals"
)

func New() []hubcheck.Rule {
	return []hubcheck.Rule{
		twofactor.New(),
		defaultrepopermission.New(),
		actionspermissions.New(),
		workflowapprovals.New(),
		orgadmins.New(),
	}
}
