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
    ListOrganizations() ([]Organization, error)
    GetOrg(id string) (Organization, error)
    GetGitHubActionsOrgPermissions(id string) (ActionsOrgPermissions, error)
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
        logger:      logger,
        accessToken: accessToken,
        cli:         cli,
    }, nil
}

type client struct {
    accessToken string
    cli         *http.Client
    logger      hublog.Logger
}

func (c *client) GetGitHubActionsOrgPermissions(id string) (ActionsOrgPermissions, error) {
    resp := &ActionsOrgPermissions{}
    if err := getRequest(c, "GET", "orgs/"+url.PathEscape(id)+"/actions/permissions", resp); err != nil {
        return ActionsOrgPermissions{}, fmt.Errorf(
            "failed to fetch GitHub Actions permissions for organization %s (%w)",
            id,
            err,
        )
    }
    resp.client = c
    return *resp, nil
}

func (c *client) GetOrg(id string) (Organization, error) {
    org := &Organization{}
    if err := getRequest(c, "GET", "orgs/"+url.PathEscape(id), org); err != nil {
        return Organization{}, fmt.Errorf("failed to fetch organization %s (%w)", id, err)
    }
    org.client = c
    return *org, nil
}

func (c *client) ListOrganizations() ([]Organization, error) {
    orgs, err := listRequest[Organization](c, "GET", "user/orgs")
    if err != nil {
        return nil, fmt.Errorf("failed to list organizations (%w)", err)
    }
    for _, org := range orgs {
        org.client = c
    }
    return orgs, nil
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
