package github

import "time"

type Organization struct {
    client Client `json:"-"`

    Login                   string    `json:"login"`
    Id                      int       `json:"id"`
    NodeId                  string    `json:"node_id"`
    Url                     string    `json:"url"`
    ReposUrl                string    `json:"repos_url"`
    EventsUrl               string    `json:"events_url"`
    HooksUrl                string    `json:"hooks_url"`
    IssuesUrl               string    `json:"issues_url"`
    MembersUrl              string    `json:"members_url"`
    PublicMembersUrl        string    `json:"public_members_url"`
    AvatarUrl               string    `json:"avatar_url"`
    Description             string    `json:"description"`
    Name                    string    `json:"name"`
    Company                 string    `json:"company"`
    Blog                    string    `json:"blog"`
    Location                string    `json:"location"`
    Email                   string    `json:"email"`
    TwitterUsername         string    `json:"twitter_username"`
    IsVerified              bool      `json:"is_verified"`
    HasOrganizationProjects bool      `json:"has_organization_projects"`
    HasRepositoryProjects   bool      `json:"has_repository_projects"`
    PublicRepos             int       `json:"public_repos"`
    PublicGists             int       `json:"public_gists"`
    Followers               int       `json:"followers"`
    Following               int       `json:"following"`
    HtmlUrl                 string    `json:"html_url"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
    Type                    string    `json:"type"`
    TotalPrivateRepos       int       `json:"total_private_repos"`
    OwnedPrivateRepos       int       `json:"owned_private_repos"`
    PrivateGists            int       `json:"private_gists"`
    DiskUsage               int       `json:"disk_usage"`
    Collaborators           int       `json:"collaborators"`
    BillingEmail            string    `json:"billing_email"`
    Plan                    struct {
        Name         string `json:"name"`
        Space        int    `json:"space"`
        PrivateRepos int    `json:"private_repos"`
        FilledSeats  int    `json:"filled_seats"`
        Seats        int    `json:"seats"`
    } `json:"plan"`
    DefaultRepositoryPermission          string `json:"default_repository_permission"`
    MembersCanCreateRepositories         bool   `json:"members_can_create_repositories"`
    TwoFactorRequirementEnabled          bool   `json:"two_factor_requirement_enabled"`
    MembersAllowedRepositoryCreationType string `json:"members_allowed_repository_creation_type"`
    MembersCanCreatePublicRepositories   bool   `json:"members_can_create_public_repositories"`
    MembersCanCreatePrivateRepositories  bool   `json:"members_can_create_private_repositories"`
    MembersCanCreateInternalRepositories bool   `json:"members_can_create_internal_repositories"`
    MembersCanCreatePages                bool   `json:"members_can_create_pages"`
    MembersCanForkPrivateRepositories    bool   `json:"members_can_fork_private_repositories"`
}

func (o Organization) GetActionsPermissions() (ActionsOrgPermissions, error) {
    return o.client.GetGitHubActionsOrgPermissions(o.Login)
}
