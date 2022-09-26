package repo

import (
	"github.com/gobwas/glob"
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/rules/repo/actionspermissions"
	"go.debugged.it/hubcheck/rules/repo/containing"
	"go.debugged.it/hubcheck/rules/repo/gitignore"
	"go.debugged.it/hubcheck/rules/repo/ide"
	"go.debugged.it/hubcheck/rules/repo/license"
	"go.debugged.it/hubcheck/rules/repo/readme"
	"go.debugged.it/hubcheck/rules/repo/vulnalerts"
)

func New(ignoreFilesList []glob.Glob, containingTerm string) []hubcheck.RepoRule {
	return []hubcheck.RepoRule{
		actionspermissions.New(),
		vulnalerts.New(),
		license.New(),
		readme.New(),
		gitignore.New(),
		ide.New(),
		containing.New(ignoreFilesList, containingTerm),
	}
}
