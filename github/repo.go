package github

import "time"

type Repository struct {
	client   Client `json:"-"`
	orgLogin string `json:"-"`

	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	Fork            bool      `json:"fork"`
	Url             string    `json:"url"`
	Homepage        string    `json:"homepage"`
	ForksCount      int       `json:"forks_count"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Size            int       `json:"size"`
	DefaultBranch   string    `json:"default_branch"`
	OpenIssuesCount int       `json:"open_issues_count"`
	IsTemplate      bool      `json:"is_template"`
	Topics          []string  `json:"topics"`
	HasIssues       bool      `json:"has_issues"`
	HasProjects     bool      `json:"has_projects"`
	HasWiki         bool      `json:"has_wiki"`
	HasPages        bool      `json:"has_pages"`
	HasDownloads    bool      `json:"has_downloads"`
	Archived        bool      `json:"archived"`
	Disabled        bool      `json:"disabled"`
	Visibility      string    `json:"visibility"`
	PushedAt        time.Time `json:"pushed_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r Repository) GetActionsPermissions() (*ActionsPermissions, error) {
	return r.client.GetGitHubActionsRepoPermissions(r.orgLogin, r.Name)
}
