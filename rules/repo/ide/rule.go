package ide

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
	return "https://docs.github.com/en/get-started/getting-started-with-git/ignoring-files"
}

func (r rule) Name() string {
	return "IDE artifacts"
}

func (r rule) Description() string {
	return "Repositories should not have IDE artifacts committed (such as .vscode, .idea, *.iml, etc.)"
}

func (r rule) ID() string {
	return "ide"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	repoContents, err := repo.ListContents()
	if err != nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       "Cannot check IDE artifacts",
				Description: fmt.Sprintf("Failed to list repository contents. (%v)", err),
				DocURL:      r.DocURL(),
			},
		}, nil
	}

	results := []hubcheck.RuleResult{}

	for _, f := range repoContents {
		if f.Name == ".vscode" || f.Name == ".idea" || strings.HasSuffix(f.Name, ".iml") {
			results = append(results, hubcheck.RuleResult{
				Level:      hublog.Warning,
				Repository: repo.Name,
				Title:      "IDE artifacts found",
				Description: fmt.Sprintf(
					"IDE artifact found at %s. Please remove this IDE artifact for contributor friendliness.",
					f.Path,
				),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/delete/%s/%s",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
					f.Path,
				),
				DocURL: r.DocURL(),
			})
		}
	}

	if len(results) == 0 {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Notice,
				Repository:  repo.Name,
				Title:       "No IDE artifacts found",
				Description: fmt.Sprintf("The repository has no repository artifacts committed."),
				DocURL:      r.DocURL(),
			},
		}, nil
	}

	return results, nil
}
