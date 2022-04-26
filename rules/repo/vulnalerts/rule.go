package vulnalerts

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
	return "https://docs.github.com/en/code-security/dependabot/dependabot-alerts/about-dependabot-alerts"
}

func (r rule) Name() string {
	return "Vulnerability alerts"
}

func (r rule) Description() string {
	return "Vulnerability alerts warn if a library used as a dependency has a known vulnerability and should be updated."
}

func (r rule) ID() string {
	return "repo-vulnerability-alerts"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	vulnerabilityAlertsEnabled, err := repo.VulnerabilityAlertsEnabled()
	if err != nil {
		return nil, err
	}

	okResult := []hubcheck.RuleResult{
		{
			Level:       hublog.Notice,
			Repository:  repo.Name,
			Title:       "Vulnerability alerts are enabled",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/settings/security_analysis",
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
			Title:       "Vulnerability alerts are disabled",
			Description: r.Description(),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/settings/security_analysis",
				url.QueryEscape(org.Login),
				url.QueryEscape(repo.Name),
			),
			DocURL: r.DocURL(),
		},
	}

	if vulnerabilityAlertsEnabled {
		return okResult, nil
	} else {
		return errResult, nil
	}
}
