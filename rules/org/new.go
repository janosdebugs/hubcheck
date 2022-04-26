package org

import (
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/org/actionspermissions"
	"go.debugged.it/hubcheck/rules/org/defaultrepopermission"
	"go.debugged.it/hubcheck/rules/org/orgadmins"
	"go.debugged.it/hubcheck/rules/org/twofactor"
	"go.debugged.it/hubcheck/rules/org/workflowapprovals"
)

func New() []hubcheck.OrgRule {
	return []hubcheck.OrgRule{
		twofactor.New(),
		defaultrepopermission.New(),
		actionspermissions.New(),
		workflowapprovals.New(),
		orgadmins.New(),
	}
}
