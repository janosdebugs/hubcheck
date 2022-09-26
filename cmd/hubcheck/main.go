package main

import (
	"flag"
	"fmt"
	"os"

	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/hublog"
	orgRules "go.debugged.it/hubcheck/rules/org"
	repoRules "go.debugged.it/hubcheck/rules/repo"
)

func main() {
	logger := hublog.New()
	org := ""
	printRules := false

	flag.StringVar(&org, "org", "", "Organization ID (in case you have access to more than one organization)")
	flag.BoolVar(&printRules, "rules", false, "List all rules.")

	flag.Parse()

	orgRuleList := orgRules.New()
	repoRuleList := repoRules.New()
	if printRules {
		for _, rule := range orgRuleList {
			fmt.Printf("## %s\n\n%s\n\nRead more: %s\n\n", rule.Name(), rule.Description(), rule.DocURL())
		}
		for _, rule := range repoRuleList {
			fmt.Printf("## %s\n\n%s\n\nRead more: %s\n\n", rule.Name(), rule.Description(), rule.DocURL())
		}
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		logger.WithLevel(hublog.Error).Logf("Please set the GITHUB_TOKEN environment variable.")
		os.Exit(1)
	}

	hc, err := hubcheck.New(logger, token, org)
	if err != nil {
		logger.WithLevel(hublog.Error).Loge(err)
		os.Exit(1)
	}

	results, err := hc.Run(
		orgRuleList,
		repoRuleList,
	)
	if err != nil {
		logger.WithLevel(hublog.Error).Loge(err)
		os.Exit(1)
	}

	failed := false
	for rule, resultList := range results {
		for _, result := range resultList {
			if result.Level == hublog.Warning || result.Level == hublog.Error {
				failed = true
			}
			name := result.Title
			if result.Repository != "" {
				name = fmt.Sprintf("%s on %s", name, result.Repository)
			}
			if result.FixURL != "" && result.DocURL != "" {
				logger.WithLevel(result.Level).Logf(
					"[%s] %s: %s Please visit %s to fix this issue. More information may be available at %s.",
					rule,
					name,
					result.Description,
					result.FixURL,
					result.DocURL,
				)
			} else {
				logger.WithLevel(result.Level).Logf(
					"[%s] %s: %s",
					rule,
					name,
					result.Description,
				)
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}
