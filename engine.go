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
	Run(org github.Organization) ([]RuleResult, error)
}

type RuleResult struct {
	Level       hublog.Level
	Title       string
	Description string
	FixURL      string
	DocURL      string
}

type HubCheck interface {
	Run(rules ...Rule) (map[string][]RuleResult, error)
}

func New(logger hublog.Logger, token string, orgID string) (HubCheck, error) {
	if token == "" {
		return nil, fmt.Errorf("no access token provided")
	}
	ghClient, err := github.NewClient(logger, token)
	if err != nil {
		return nil, err
	}
	var org github.Organization
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
	org    github.Organization
	logger hublog.Logger
}

func (h hubCheck) Run(rules ...Rule) (map[string][]RuleResult, error) {
	results := map[string][]RuleResult{}
	for _, rule := range rules {
		h.logger.WithLevel(hublog.Debug).Logf("Processing rule %s...", rule.ID())
		result, err := rule.Run(h.org)
		if err != nil {
			results[rule.ID()] = []RuleResult{
				{
					Level:       hublog.Error,
					Title:       "Rule execution failed",
					Description: err.Error(),
				},
			}
		} else {
			results[rule.ID()] = result
		}
	}
	return results, nil
}
