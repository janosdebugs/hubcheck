package workflowapprovals

import (
	"fmt"
	"net/url"

	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/github"
	"go.debugged.it/hubcheck/hublog"
)

func New() hubcheck.OrgRule {
	return &rule{}
}

type rule struct {
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/actions/managing-workflow-runs/approving-workflow-runs-from-public-forks"
}

func (r rule) Name() string {
	return "Require workflow approvals (manual)"
}

func (r rule) Description() string {
	return "Workflow approvals cannot be checked automatically, please check them manually. When a pull request is submitted from a fork, GitHub actions should not be run automatically or you risk exposing sensitive credentials to untrusted code. You should change your settings to require approvals from a project maintainer in order to run workflows."
}

func (r rule) ID() string {
	return "github-actions-workflow-approvals"
}

func (r rule) Run(org *github.Organization) ([]hubcheck.RuleResult, error) {
	return []hubcheck.RuleResult{
		{
			Level:       hublog.Info,
			Title:       "Workflow approval requirements",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/organizations/%s/settings/actions",
				url.QueryEscape(org.Login),
			),
			DocURL: r.DocURL(),
		},
	}, nil
}
