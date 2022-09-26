package readme

import (
	"fmt"
	"net/url"
	"strings"

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
	return "https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes"
}

func (r rule) Name() string {
	return "Repository README"
}

func (r rule) Description() string {
	return "Repositories should have a README file."
}

func (r rule) ID() string {
	return "readme"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	repoContents, err := repo.ListContents()
	if err != nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       "Cannot check README",
				Description: fmt.Sprintf("Failed to list repository contents. (%v)", err),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/new/%s?readme=1",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	var found *github.RepoDirEntry
	for _, f := range repoContents {
		if f.Path == f.Name && strings.HasPrefix(f.Name, "README") {
			found = &f
			break
		}
	}

	if found == nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Error,
				Repository:  repo.Name,
				Title:       "Repository has no README",
				Description: fmt.Sprintf("The repository has no README file."),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/new/%s?readme=1",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	if found.Size < 1024 {
		return []hubcheck.RuleResult{
			{
				Level:      hublog.Warning,
				Repository: repo.Name,
				Title:      "Repository has very short README",
				Description: fmt.Sprintf(
					"The repository has a README file named %s, but it is too short to be useful.",
					found.Path,
				),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/edit/%s/%s",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
					found.Path,
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	return []hubcheck.RuleResult{
		{
			Level:       hublog.Notice,
			Repository:  repo.Name,
			Title:       "Repository has a README",
			Description: fmt.Sprintf("The repository has a README file named %s.", found),
			DocURL:      r.DocURL(),
		},
	}, nil
}
