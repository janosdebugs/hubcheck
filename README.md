# HubCheck: an organization security checker for GitHub

This tool checks your organization settings on GitHub. **This tool is not affiliated with GitHub.**

## Running this tool

In other to run this tool you will have to export the `GITHUB_TOKEN` environment variable to grant permissions to read the organization settings.

You can then run:

```
go run cmd/hubcheck/main.go
```

## Rules

<!-- region Rules -->

### Two-factor enforcement

To ensure that authorized members of an organization are not easily compromised by a password theft you should enforce two-factor authentication in your organization.

Read more: https://docs.github.com/en/organizations/keeping-your-organization-secure/managing-two-factor-authentication-for-your-organization/requiring-two-factor-authentication-in-your-organization

### Default repository permissions

To ensure that organization members cannot carry out destructive actions, such as force-pushing and thereby deleting history, the default repository permissions should not be set to admin.

Read more: https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/setting-base-permissions-for-an-organization

### Limit GitHub Actions on the organization

Allowing all GitHub Actions to run introduces the risk of accidentally exposing sensitive credentials to untrusted, or even malicious developers.

Read more: https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization

### Require workflow approvals (manual)

Workflow approvals cannot be checked automatically, please check them manually. When a pull request is submitted from a fork, GitHub actions should not be run automatically or you risk exposing sensitive credentials to untrusted code. You should change your settings to require approvals from a project maintainer in order to run workflows.

Read more: https://docs.github.com/en/actions/managing-workflow-runs/approving-workflow-runs-from-public-forks

### Organizations should have between 2 and 5 administrators

If an organization has only one administrator it is easy to lose access to it. If an organization has too many administrators it means that permissions are handled too liberally.

Read more: https://docs.github.com/en/organizations/managing-membership-in-your-organization

### Limit GitHub Actions on repositories

Allowing all GitHub Actions to run introduces the risk of accidentally exposing sensitive credentials to untrusted, or even malicious developers.

Read more: https://docs.github.com/en/organizations/managing-organization-settings/disabling-or-limiting-github-actions-for-your-organization

### Vulnerability alerts

Vulnerability alerts warn if a library used as a dependency has a known vulnerability and should be updated.

Read more: https://docs.github.com/en/code-security/dependabot/dependabot-alerts/about-dependabot-alerts

### Repository license

Public repositories should have a license.

Read more: https://docs.github.com/articles/adding-a-license-to-a-repository/

### Repository README

Repositories should have a README file.

Read more: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes

### Repository .gitignore

Repositories should have a .gitignore file.

Read more: https://docs.github.com/en/get-started/getting-started-with-git/ignoring-files

### IDE artifacts

Repositories should not have IDE artifacts committed (such as .vscode, .idea, *.iml, etc.)

Read more: https://docs.github.com/en/get-started/getting-started-with-git/ignoring-files

### Files containing a user-configurable term

This rule alerts for files containing a user-configurable term.

<!-- endregion -->













