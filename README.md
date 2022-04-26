# HubCheck: an organization security checker for GitHub

This tool checks your organization settings on GitHub. **This tool is not affiliated with GitHub.**

## Running this tool

In other to run this tool you will have to export the `GITHUB_TOKEN` environment variable to grant permissions to read the organization settings.

You can then run:

```
go run cmd/hubcheck/main.go
```

## Rules

### Two-factor enforcement

To ensure that authorized members of an organization are not easily compromised by a password theft you should enforce two-factor authentication in your organization.

Read more: https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-two-factor-authentication-for-your-organization/requiring-two-factor-authentication-in-your-organization

### Default repository permissions

To ensure that organization members cannot carry out destructive actions, such as force-pushing and thereby deleting history, the default repository permissions should not be set to admin.

Read more: https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/setting-base-permissions-for-an-organization

### Limit GitHub Actions

Allowing all GitHub Actions to run introduces the risk of accidentally exposing sensitive credentials to untrusted, or even malicious developers.

Read more: https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization
