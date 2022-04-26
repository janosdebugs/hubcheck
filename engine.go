package hubcheck

import (
	"fmt"

	"go.debugged.it/hubcheck/github"
	"go.debugged.it/hubcheck/hublog"
)

type Rule interface {
	ID() string
	Name() string
	Description() string
	DocURL() string
}

type OrgRule interface {
	Rule
	Run(org *github.Organization) ([]RuleResult, error)
}

type RepoRule interface {
	Rule
	Run(org *github.Organization, repo *github.Repository) ([]RuleResult, error)
}

type RuleResult struct {
	Level       hublog.Level
	Repository  string
	Title       string
	Description string
	FixURL      string
	DocURL      string
}

type HubCheck interface {
	Run(orgRules []OrgRule, repoRules []RepoRule) (map[string][]RuleResult, error)
}

func New(logger hublog.Logger, token string, orgID string) (HubCheck, error) {
	if token == "" {
		return nil, fmt.Errorf("no access token provided")
	}
	ghClient, err := github.NewClient(logger, token)
	if err != nil {
		return nil, err
	}
	var org *github.Organization
	if orgID != "" {
		org, err = ghClient.GetOrg(orgID)
		if err != nil {
			return nil, err
		}
	} else {
		orgs, err := ghClient.ListOrganizations()
		if err != nil {
			return nil, err
		}
		if len(orgs) == 0 {
			return nil, fmt.Errorf("no organization found using your access token")
		} else if len(orgs) == 1 {
			org = orgs[0]
		} else {
			return nil, fmt.Errorf("more than one organization found in your account, please provide the organization ID using the -org parameter")
		}

	}
	return &hubCheck{
		logger: logger,
		client: ghClient,
		org:    org,
	}, nil
}

type hubCheck struct {
	client github.Client
	org    *github.Organization
	logger hublog.Logger
}

func (h hubCheck) Run(orgRules []OrgRule, repoRules []RepoRule) (map[string][]RuleResult, error) {
	results := map[string][]RuleResult{}
	for _, rule := range orgRules {
		h.logger.WithLevel(hublog.Debug).Logf("Processing rule %s...", rule.ID())
		result, err := rule.Run(h.org)
		if err != nil {
			results[rule.ID()] = []RuleResult{
				{
					Level:       hublog.Warning,
					Title:       "Rule execution failed",
					Description: err.Error(),
				},
			}
		} else {
			results[rule.ID()] = result
		}
	}
	repos, err := h.org.ListRepositories()
	if err != nil {
		return results, fmt.Errorf("failed to list organization repositories (%w)", err)
	}
	for _, rule := range repoRules {
		for _, repo := range repos {
			if _, ok := results[rule.ID()]; !ok {
				results[rule.ID()] = nil
			}
			h.logger.WithLevel(hublog.Debug).Logf("Processing rule %s on repository %s...", rule.ID(), repo.Name)
			result, err := rule.Run(h.org, repo)
			if err != nil {
				results[rule.ID()] = append(
					results[rule.ID()],
					RuleResult{
						Level:       hublog.Warning,
						Title:       fmt.Sprintf("Rule execution failed on repository %s", repo.Name),
						Description: err.Error(),
					},
				)
			} else {
				results[rule.ID()] = append(results[rule.ID()], result...)
			}
		}
	}
	return results, nil
}
