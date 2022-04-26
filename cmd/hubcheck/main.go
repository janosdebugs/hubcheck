package main

import (
	"flag"
	"fmt"
	"os"

	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/hublog"
	"go.debugged.it/hubcheck/rules/actionspermissions"
	"go.debugged.it/hubcheck/rules/defaultrepopermission"
	"go.debugged.it/hubcheck/rules/orgadmins"
	"go.debugged.it/hubcheck/rules/twofactor"
	"go.debugged.it/hubcheck/rules/workflowapprovals"
)

func main() {
	logger := hublog.New()
	org := ""
	printRules := false

	flag.StringVar(&org, "org", "", "Organization ID (in case you have access to more than one organization)")
	flag.BoolVar(&printRules, "rules", false, "List all rules.")

	flag.Parse()

	rules := []hubcheck.Rule{
		twofactor.New(),
		defaultrepopermission.New(),
		actionspermissions.New(),
		workflowapprovals.New(),
		orgadmins.New(),
	}
	if printRules {
		for _, rule := range rules {
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
		rules...,
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
			if result.FixURL != "" && result.DocURL != "" {
				logger.WithLevel(result.Level).Logf(
					"[%s] %s: %s Please visit %s to change this setting. More information is available at %s.",
					rule,
					result.Title,
					result.Description,
					result.FixURL,
					result.DocURL,
				)
			} else {
				logger.WithLevel(result.Level).Logf(
					"[%s] %s: %s",
					rule,
					result.Title,
					result.Description,
				)
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}
