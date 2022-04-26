package twofactor

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

func (r rule) Name() string {
	return "Two-factor enforcement"
}

func (r rule) Description() string {
	return "To ensure that authorized members of an organization are not easily compromised by a password theft you should enforce two-factor authentication in your organization."
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-two-factor-authentication-for-your-organization/requiring-two-factor-authentication-in-your-organization"
}

func (r rule) ID() string {
	return "two-factor"
}

func (r rule) Run(org *github.Organization) ([]hubcheck.RuleResult, error) {
	if org.TwoFactorRequirementEnabled == nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Title:       "OrgRule execution failed",
				Description: "Are you an admin?",
				FixURL: fmt.Sprintf(
					"https://github.com/organizations/%s/settings/security",
					url.QueryEscape(org.Login),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}
	if *org.TwoFactorRequirementEnabled {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Notice,
				Title:       "Two-factor authentication enforcement is enabled",
				Description: r.Description(),
				FixURL: fmt.Sprintf(
					"https://github.com/organizations/%s/settings/security",
					url.QueryEscape(org.Login),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}
	return []hubcheck.RuleResult{
		{
			Level:       hublog.Error,
			Title:       "Two-factor authentication enforcement is not enabled",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/organizations/%s/settings/security",
				url.QueryEscape(org.Login),
			),
			DocURL: r.DocURL(),
		},
	}, nil
}
