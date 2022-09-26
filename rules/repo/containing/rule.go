package containing

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gobwas/glob"
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/github"
	"go.debugged.it/hubcheck/hublog"
)

func New(ignoreFilesList []glob.Glob, term string) hubcheck.RepoRule {
	return &rule{
		ignoreFilesList: ignoreFilesList,
		term:            strings.ToLower(term),
	}
}

type rule struct {
	term            string
	ignoreFilesList []glob.Glob
}

func (r rule) DocURL() string {
	return ""
}

func (r rule) Name() string {
	if r.term == "" {
		return "Files containing a user-configurable term"
	}
	return fmt.Sprintf("Files containing '%s'", r.term)
}

func (r rule) Description() string {
	return "This rule alerts for files containing a user-configurable term."
}

func (r rule) ID() string {
	return "containing"
}

func (r rule) Run(org *github.Organization, repo *github.Repository) ([]hubcheck.RuleResult, error) {
	if r.term == "" {
		return nil, nil
	}

	repoContents, err := repo.ListContents()
	if err != nil {
		return []hubcheck.RuleResult{
			{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       "Cannot list repository contents.",
				Description: fmt.Sprintf("Failed to list repository contents. (%v)", err),
				DocURL:      r.DocURL(),
			},
		}, nil
	}

	var results []hubcheck.RuleResult
	for _, f := range repoContents {
		if f.Size > 204800 {
			results = append(results, hubcheck.RuleResult{
				Level:       hublog.Debug,
				Repository:  repo.Name,
				Title:       "File too large for analysis",
				Description: fmt.Sprintf("File %s is too large for content analysis, skipping..."),
			})
			continue
		}
		ignored := false
		for _, g := range r.ignoreFilesList {
			if g.Match(f.Path) {
				ignored = true
				break
			}
		}
		if ignored {
			results = append(results, hubcheck.RuleResult{
				Level:       hublog.Debug,
				Repository:  repo.Name,
				Title:       "File matches ignore pattern",
				Description: fmt.Sprintf("File %s matches ignore pattern, skipping analysis..."),
			})
			continue
		}
		contents, err := f.GetContents()
		if err != nil {
			results = append(results, hubcheck.RuleResult{
				Level:       hublog.Warning,
				Repository:  repo.Name,
				Title:       fmt.Sprintf("Failed to fetch %s", f.Path),
				Description: err.Error(),
			})
		}

		if strings.Contains(strings.ToLower(string(contents)), r.term) {
			results = append(results, hubcheck.RuleResult{
				Level:       hublog.Error,
				Repository:  repo.Name,
				Title:       fmt.Sprintf("File %s contains '%s'", f.Path, r.term),
				Description: fmt.Sprintf("This file contains the search term '%s'.", r.term),
				FixURL: fmt.Sprintf(
					"https://github.com/%s/%s/edit/%s/%s",
					url.QueryEscape(org.Login),
					url.QueryEscape(repo.Name),
					url.QueryEscape(repo.DefaultBranch),
					f.Path,
				),
			})
		}
	}

	return results, nil
}
