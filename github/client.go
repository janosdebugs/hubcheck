package github

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.debugged.it/hubcheck/hublog"
)

type Client interface {
	ListOrganizations() ([]*Organization, error)
	GetOrg(login string) (*Organization, error)
	GetGitHubActionsOrgPermissions(login string) (*ActionsPermissions, error)
	ListOrgAdmins(login string) ([]*OrgMember, error)
	ListOrgRepositories(login string) ([]*Repository, error)
	GetGitHubActionsRepoPermissions(login string, repoName string) (*ActionsPermissions, error)
	RepoVulnerabilityAlertsEnabled(login string, repoName string) (bool, error)
	ListContents(login string, repoName string) ([]RepoDirEntry, error)
}

func NewClient(logger hublog.Logger, accessToken string) (Client, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain the system certificate pool (%w)", err)
	}

	cli := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs:    certPool,
				MinVersion: tls.VersionTLS13,
			},
		},
	}

	return &client{
		logger:           logger,
		accessToken:      accessToken,
		cli:              cli,
		repoContentCache: map[string][]RepoDirEntry{},
	}, nil
}

type client struct {
	accessToken string
	cli         *http.Client
	logger      hublog.Logger

	repoContentCache map[string][]RepoDirEntry
}

func (c *client) RepoVulnerabilityAlertsEnabled(login string, repoName string) (bool, error) {
	statusCode, _, body, err := c.request(
		"GET",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/vulnerability-alerts", login, repoName),
	)
	if err != nil {
		return false, fmt.Errorf("failed to query repository %s vulnerability alert settings (%w)", repoName, err)
	}
	switch statusCode {
	case 204:
		return true, nil
	case 403:
		return false, fmt.Errorf(
			"you do not have permissions to query vulnerability alerts for repository %s",
			repoName,
		)
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf(
			"unexpected HTTP status code for repository %s vulnerability alerts: %d (%s)",
			repoName,
			statusCode,
			body,
		)
	}
}

func (c *client) GetGitHubActionsRepoPermissions(login string, repoName string) (*ActionsPermissions, error) {
	resp := &ActionsPermissions{}
	if err := getRequest(
		c,
		"GET",
		"repos/"+url.PathEscape(login)+"/"+url.PathEscape(repoName)+"/actions/permissions",
		resp,
	); err != nil {
		return nil, fmt.Errorf(
			"Failed to fetch GitHub Actions permissions for repo %s/%s. (Did you forget to add the admin permissions to your GitHub token?) (%w)",
			login,
			repoName,
			err,
		)
	}
	resp.client = c
	return resp, nil
}

func (c *client) ListOrgRepositories(login string) ([]*Repository, error) {
	repos, err := listRequest[*Repository](c, "GET", fmt.Sprintf("orgs/%s/repos", url.PathEscape(login)))
	if err != nil {
		return nil, fmt.Errorf("Failed to list organization repositories. (%w)", err)
	}
	for _, repo := range repos {
		repo.client = c
		repo.orgLogin = login
	}
	return repos, nil
}

func (c *client) ListOrgAdmins(id string) ([]*OrgMember, error) {
	members, err := listRequest[*OrgMember](c, "GET", fmt.Sprintf("orgs/%s/members?role=admin", url.PathEscape(id)))
	if err != nil {
		return nil, fmt.Errorf("Failed to list organization %s members. (%w)", id, err)
	}
	return members, nil
}

func (c *client) GetGitHubActionsOrgPermissions(id string) (*ActionsPermissions, error) {
	resp := &ActionsPermissions{}
	if err := getRequest(c, "GET", "orgs/"+url.PathEscape(id)+"/actions/permissions", resp); err != nil {
		return nil, fmt.Errorf(
			"Failed to fetch GitHub Actions permissions for organization %s. (%w)",
			id,
			err,
		)
	}
	resp.client = c
	return resp, nil
}

func (c *client) GetOrg(id string) (*Organization, error) {
	org := &Organization{}
	if err := getRequest(c, "GET", "orgs/"+url.PathEscape(id), org); err != nil {
		return nil, fmt.Errorf("Failed to fetch organization %s. (%w)", id, err)
	}
	org.client = c
	return org, nil
}

func (c *client) ListOrganizations() ([]*Organization, error) {
	orgs, err := listRequest[*Organization](c, "GET", "user/orgs")
	if err != nil {
		return nil, fmt.Errorf("Failed to list organizations. (%w)", err)
	}
	for _, org := range orgs {
		org.client = c
	}
	return orgs, nil
}

type FileType string

//goland:noinspection GoUnusedConst
const (
	FileTypeDir       FileType = "dir"
	FileTypeFile      FileType = "file"
	FileTypeSubmodule FileType = "submodule"
	FileTypeSymlink   FileType = "symlink"
)

type RepoDirEntry struct {
	Type FileType `json:"type"`
	Size int      `json:"size"`
	Name string   `json:"name"`
	Path string   `json:"path"`
	Sha  string   `json:"sha"`
}

func (c *client) listContents(orgID string, repoID string, path string) ([]RepoDirEntry, error) {
	var response []RepoDirEntry
	urlPath := fmt.Sprintf("repos/%s/%s/contents/%s", url.PathEscape(orgID), url.PathEscape(repoID), path)
	if err := getRequest(c, "GET", urlPath, &response); err != nil {
		return nil, err
	}

	result := response
	for _, item := range response {
		if item.Type != FileTypeDir {
			continue
		}
		subResult, err := c.listContents(orgID, repoID, item.Path)
		if err != nil {
			return nil, err
		}
		result = append(result, subResult...)
	}
	return result, nil
}

func (c *client) ListContents(orgID string, repoID string) ([]RepoDirEntry, error) {
	if contents, ok := c.repoContentCache[orgID+"/"+repoID]; ok {
		return contents, nil
	}
	contents, err := c.listContents(orgID, repoID, "")
	if err != nil {
		return nil, err
	}
	c.repoContentCache[orgID+"/"+repoID] = contents
	return contents, nil
}

type errorResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (c *client) request(method string, url string) (statusCode int, headers http.Header, body []byte, err error) {
	c.logger.WithLevel(hublog.Debug).Logf("HTTP --> %s %s", method, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to construct HTTP request (%w)", err)
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+c.accessToken)
	req.Header.Add("User-Agent", "HubCheck")
	response, err := c.cli.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("HTTP request failed (%w)", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to read response body (%v)", err)
	}

	c.logger.WithLevel(hublog.Debug).Logf("HTTP <-- %d", response.StatusCode)

	return response.StatusCode, response.Header, body, nil
}

func getRequest[T any](c *client, method string, path string, responseObject *T) error {
	status, _, body, err := c.request(method, "https://api.github.com/"+path)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	switch status {
	case 200:
	default:
		var errDetails *errorResponse
		if err := decoder.Decode(errDetails); err != nil {
			return fmt.Errorf(
				"unexpected HTTP response code: %d (%s)",
				status,
				body,
			)
		}
		return fmt.Errorf(
			"unexpected HTTP response code: %d (%s; %s )",
			status,
			errDetails.Message,
			errDetails.DocumentationURL,
		)
	}

	if err := decoder.Decode(responseObject); err != nil {
		return fmt.Errorf("failed to decode GitHub response (%v; %s)", err, body)
	}

	return nil
}

// listRequest lists items of a certain type while observing pagination.
// This is a non-receiver method due to https://github.com/golang/go/issues/49085
func listRequest[T any](c *client, method string, path string) ([]T, error) {
	nextLink := "https://api.github.com/" + path
	var result []T
	for {
		status, headers, body, err := c.request(method, nextLink)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(bytes.NewReader(body))
		switch status {
		case 200:
		default:
			var errDetails *errorResponse
			if err := decoder.Decode(errDetails); err != nil {
				return nil, fmt.Errorf(
					"unexpected HTTP response code: %d (%s)",
					status,
					body,
				)
			}
			return nil, fmt.Errorf(
				"unexpected HTTP response code: %d (%s; %s )",
				status,
				errDetails.Message,
				errDetails.DocumentationURL,
			)
		}

		nextLink = ""
		linkHeader := headers.Get("Link")
		if linkHeader != "" {
			linkHeaderParts := strings.Split(linkHeader, ",")
			for _, linkHeaderPart := range linkHeaderParts {
				parts := strings.SplitN(linkHeaderPart, ";", 2)
				if len(parts) == 2 {
					relPart := strings.TrimSpace(parts[1])
					if relPart == "rel=\"next\"" {
						nextLink = strings.Trim(parts[0], "<>")
					}
				}
			}
		}

		var items []T
		if err := decoder.Decode(&items); err != nil {
			return nil, fmt.Errorf("failed to decode GitHub response (%v; %s)", err, body)
		}
		if len(items) == 0 {
			return result, nil
		}

		result = append(result, items...)

		if nextLink == "" {
			return result, nil
		}
	}
}
