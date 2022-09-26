package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gobwas/glob"
	"go.debugged.it/hubcheck"
	"go.debugged.it/hubcheck/hublog"
	orgRules "go.debugged.it/hubcheck/rules/org"
	repoRules "go.debugged.it/hubcheck/rules/repo"
)

func main() {
	org := ""
	printRules := false
	logLevel := string(hublog.Info)
	ignoreFiles := "vendor/**;venv/**;virtualenv/**"
	reportFilesContaining := ""

	flag.StringVar(&org, "org", "", "Organization ID (in case you have access to more than one organization)")
	flag.BoolVar(&printRules, "rules", false, "List all rules.")
	flag.StringVar(&logLevel, "log-level", logLevel, "Minimum log level (debug, info, notice, warning, error).")
	flag.StringVar(&ignoreFiles, "ignore-files", ignoreFiles, "Vendor directories to ignore from analysis.")
	flag.StringVar(&reportFilesContaining, "report-files-containing", reportFilesContaining, "Report files containing this term.")
	flag.Parse()

	logger := hublog.New(hublog.Level(logLevel))

	var ignoreFilesList []glob.Glob
	if ignoreFiles != "" {
		for _, f := range strings.Split(ignoreFiles, ";") {
			g, err := glob.Compile(f)
			if err != nil {
				logger.WithLevel(hublog.Error).Loge(err)
			}
			ignoreFilesList = append(ignoreFilesList, g)
		}
	}

	orgRuleList := orgRules.New()
	repoRuleList := repoRules.New(ignoreFilesList, reportFilesContaining)
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

	print("# Report for the " + org + " GitHub organization\n\n")
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

			printResult(org, rule, result, hublog.Level(logLevel))
		}
	}
	if failed {
		os.Exit(1)
	}
}

func printResult(org string, rule string, result hubcheck.RuleResult, level hublog.Level) {
	prefix := ""
	suffix := "\033[0m\n\n"
	switch result.Level {
	case hublog.Debug:
		if level != hublog.Debug {
			return
		}
		prefix = "\033[37m## ⚙️ "
	case hublog.Info:
		if level != hublog.Debug && level != hublog.Info {
			return
		}
		prefix = "\033[36m## ℹ️ "
	case hublog.Notice:
		if level != hublog.Debug && level != hublog.Info && level != hublog.Notice {
			return
		}
		prefix = "\033[32m## ✅️ "
	case hublog.Warning:
		if level != hublog.Debug && level != hublog.Info && level != hublog.Notice && level != hublog.Warning {
			return
		}
		prefix = "\033[33m## ⚠️"
	case hublog.Error:
		prefix = "\033[31m## ❌ "
	}

	print(prefix + result.Title + " (`" + rule + "`)" + suffix)
	print(result.Description + "\n\n")
	if result.Repository != "" {
		print("- \033[1m**Repository:**\033[0m [" + org + "/" + result.Repository + "](https://github.com/" + org + "/" + result.Repository + ")\n")
	} else {
		print("- \033[1m**Organization:**\033[0m [" + org + "](https://github.com/" + org + ")\n")
	}
	if result.FixURL != "" {
		print("- \033[1m**Quick fix:**\033[0m " + result.FixURL + "\n")
	}
	if result.DocURL != "" {
		print("- \033[1m**Documentation:**\033[0m " + result.DocURL + "\n")
	}
	print("\n---\n\n")
}
