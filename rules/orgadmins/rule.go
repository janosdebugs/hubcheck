package orgadmins

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

func (r rule) Name() string {
	return "Organizations should have between 2 and 5 administrators"
}

func (r rule) Description() string {
	return "If an organization has only one administrator it is easy to lose access to it. If an organization has too many administrators it means that permissions are handled too liberally."
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/organizations/managing-membership-in-your-organization"
}

func (r rule) ID() string {
	return "organization-admins"
}

func (r rule) Run(org github.Organization) ([]hubcheck.RuleResult, error) {
	members, err := org.ListAdmins()
	if err != nil {
		return nil, err
	}
	if len(members) == 1 {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Error,
				Title:       "Your organization has only one admin",
				Description: r.Description(),
				FixURL:      fmt.Sprintf("https://github.com/orgs/%s/people", url.PathEscape(org.Login)),
				DocURL:      r.DocURL(),
			},
		}, nil
	}
	if len(members) > 5 {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Error,
				Title:       fmt.Sprintf("Too many admins (%d) in your organization", len(members)),
				Description: r.Description(),
				FixURL:      fmt.Sprintf("https://github.com/orgs/%s/people", url.PathEscape(org.Login)),
				DocURL:      r.DocURL(),
			},
		}, nil
	}
	return []hubcheck.RuleResult{
		{
			Level:       hublog.Notice,
			Title:       fmt.Sprintf("%d admins in your organization", len(members)),
			Description: r.Description(),
			FixURL:      fmt.Sprintf("https://github.com/orgs/%s/people", url.PathEscape(org.Login)),
			DocURL:      r.DocURL(),
		},
	}, nil
}
