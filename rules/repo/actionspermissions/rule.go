package actionspermissions

import (
	"fmt"
	"net/url"

	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/github"
	"go.debugged.it/hubcheck/hublog"
)

func New() hubcheck.RepoRule {
	return &rule{}
}

type rule struct {
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization"
}

func (r rule) Name() string {
	return "Limit GitHub Actions on repositories"
}

func (r rule) Description() string {
	return "Allowing all GitHub Actions to run introduces the risk of accidentally exposing sensitive credentials to untrusted, or even malicious developers."
}

func (r rule) ID() string {
	return "github-actions-repo-permissions"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	actionsPermissions, err := repo.GetActionsPermissions()
	if err != nil {
		return nil, err
	}

	okResult := []hubcheck.RuleResult{
		{
			Level:       hublog.Notice,
			Repository:  repo.Name,
			Title:       "GitHub Actions are limited",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/settings/actions",
				url.QueryEscape(org.Login),
				url.QueryEscape(repo.Name),
			),
			DocURL: r.DocURL(),
		},
	}
	errResult := []hubcheck.RuleResult{
		{
			Level:       hublog.Error,
			Repository:  repo.Name,
			Title:       "GitHub Actions are not limited",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/settings/actions",
				url.QueryEscape(org.Login),
				url.QueryEscape(repo.Name),
			),
			DocURL: r.DocURL(),
		},
	}

	if actionsPermissions.EnabledRepositories == "none" {
		return okResult, nil
	}
	switch actionsPermissions.AllowedActions {
	case "selected":
		return okResult, nil
	case "local_only":
		return okResult, nil
	default:
		return errResult, nil
	}
}
