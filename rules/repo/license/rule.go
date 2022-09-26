package license

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
	return "https://docs.github.com/articles/adding-a-license-to-a-repository/"
}

func (r rule) Name() string {
	return "Repository license"
}

func (r rule) Description() string {
	return "Public repositories should have a license."
}

func (r rule) ID() string {
	return "public-repo-license"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	if repo.License != nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Notice,
				Repository:  repo.Name,
				Title:       "Repository has a license",
				Description: fmt.Sprintf("This repository is licensed under the %s.", repo.License.Name),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/community/license/new",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	if repo.Visibility != "public" {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       "Repository has no license",
				Description: fmt.Sprintf("This repository does not have a license file, but this repository is not public. Consider adding a license."),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/community/license/new",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	return []hubcheck.RuleResult{
		{
			Level:       hublog.Warning,
			Repository:  repo.Name,
			Title:       "Repository has no license",
			Description: fmt.Sprintf("This repository does not have a license file."),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/community/license/new",
				url.QueryEscape(org.Login),
				url.QueryEscape(repo.Name),
			),
			DocURL: r.DocURL(),
		},
	}, nil
}
