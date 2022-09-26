package gitignore

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
	return "https://docs.github.com/en/get-started/getting-started-with-git/ignoring-files"
}

func (r rule) Name() string {
	return "Repository .gitignore"
}

func (r rule) Description() string {
	return "Repositories should have a .gitignore file."
}

func (r rule) ID() string {
	return "gitignore"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	repoContents, err := repo.ListContents()
	if err != nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       "Cannot check .gitignore",
				Description: fmt.Sprintf("Failed to list repository contents. (%v)", err),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/new/%s?filename=.gitignore",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
				),
				DocURL: r.DocURL(),
			},
		}, nil
	}

	found := ""
	for _, f := range repoContents {
		if f.Path == "" && f.Name == ".gitignore" {
			found = f.Name
			break
		}
	}

	if found != "" {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Notice,
				Repository:  repo.Name,
				Title:       "Repository has a .gitignore file",
				Description: fmt.Sprintf("The repository has a .gitignore file named."),
				DocURL:      r.DocURL(),
			},
		}, nil
	}

	return []hubcheck.RuleResult{
		{
			Level:       hublog.Error,
			Repository:  repo.Name,
			Title:       "Repository has no .gitignore",
			Description: fmt.Sprintf("The repository has no .gitignore file."),
			FixURL: fmt.Sprintf(
				"https://github.com/%s/%s/new/%s?filename=.gitignore",
				url.QueryEscape(org.Login),
				url.QueryEscape(repo.Name),
				url.QueryEscape(repo.DefaultBranch),
			),
			DocURL: r.DocURL(),
		},
	}, nil
}
