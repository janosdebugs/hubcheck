package actionspermissions

import (
	"fmt"
	"net/url"

	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/github"
	"go.debugged.it/hubcheck/hublog"
)

func New() hubcheck.Rule {
	return &rule{}
}

type rule struct {
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization"
}

func (r rule) Name() string {
	return "Limit GitHub Actions"
}

func (r rule) Description() string {
	return "Allowing all GitHub Actions to run introduces the risk of accidentally exposing sensitive credentials to untrusted, or even malicious developers."
}

func (r rule) ID() string {
	return "github-actions-permissions"
}

func (r rule) Run(org github.Organization) ([]hubcheck.RuleResult, error) {
	actionsPermissions, err := org.GetActionsPermissions()
	if err != nil {
		return nil, err
	}

	okResult := []hubcheck.RuleResult{
		{
			Level:       hublog.Notice,
			Title:       "GitHub Actions are limited",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/organizations/%s/settings/actions",
				url.QueryEscape(org.Login),
			),
			DocURL: r.DocURL(),
		},
	}
	errResult := []hubcheck.RuleResult{
		{
			Level:       hublog.Error,
			Title:       "GitHub Actions are not limited",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/organizations/%s/settings/actions",
				url.QueryEscape(org.Login),
			),
			DocURL: r.DocURL(),
		},
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
