package repo

import (
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/repo/actionspermissions"
	"go.debugged.it/hubcheck/rules/repo/gitignore"
	"go.debugged.it/hubcheck/rules/repo/ide"
	"go.debugged.it/hubcheck/rules/repo/license"
	"go.debugged.it/hubcheck/rules/repo/readme"
	"go.debugged.it/hubcheck/rules/repo/vulnalerts"
)

func New() []hubcheck.RepoRule {
	return []hubcheck.RepoRule{
		actionspermissions.New(),
		vulnalerts.New(),
		license.New(),
		readme.New(),
		gitignore.New(),
		ide.New(),
	}
}
