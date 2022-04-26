package defaultrepopermission

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
	return "Default repository permissions"
}

func (r rule) Description() string {
	return "To ensure that organization members cannot carry out destructive actions, such as force-pushing and thereby deleting history, the default repository permissions should not be set to admin."
}

func (r rule) DocURL() string {
	return "https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/setting-base-permissions-for-an-organization"
}

func (r rule) ID() string {
	return "default-repository-permission"
}

func (r rule) Run(org *github.Organization) ([]hubcheck.RuleResult, error) {
	if org.DefaultRepositoryPermission == "" {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Title:       "OrgRule execution failed",
				Description: "Are you an admin?",
				FixURL: fmt.Sprintf(
					"https://github.com/organizations/%s/settings/member_privileges",
					url.QueryEscape(org.Login),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}
	if org.DefaultRepositoryPermission != "admin" {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Notice,
				Title:       fmt.Sprintf("Default repository permissions are %s", org.DefaultRepositoryPermission),
				Description: r.Description(),
				FixURL: fmt.Sprintf(
					"https://github.com/organizations/%s/settings/member_privileges",
					url.QueryEscape(org.Login),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}
	return []hubcheck.RuleResult{
		{
			Level:       hublog.Error,
			Title:       fmt.Sprintf("Default repository permissions are %s", org.DefaultRepositoryPermission),
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/organizations/%s/settings/member_privileges",
				url.QueryEscape(org.Login),
			),
			DocURL: r.DocURL(),
		},
	}, nil
}
