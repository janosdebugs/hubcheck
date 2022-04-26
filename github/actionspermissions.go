package github

type ActionsPermissions struct {
	client Client `json:"-"`

	EnabledRepositories string `json:"enabled_repositories"`
	AllowedActions      string `json:"allowed_actions"`
}
