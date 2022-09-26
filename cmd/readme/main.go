package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	orgRules "go.debugged.it/hubcheck/rules/org"
	repoRules "go.debugged.it/hubcheck/rules/repo"
)

func main() {
	output := "<!-- region Rules -->\n\n"
	for _, rule := range orgRules.New() {
		output += fmt.Sprintf("### %s\n\n%s\n\nRead more: %s\n\n", rule.Name(), rule.Description(), rule.DocURL())
	}
	for _, rule := range repoRules.New(nil, "") {
		output += fmt.Sprintf("### %s\n\n%s\n\n", rule.Name(), rule.Description())
		if rule.DocURL() != "" {
			output += fmt.Sprintf("Read more: %s\n\n", rule.DocURL())
		}
	}
	output += "<!-- endregion -->\n\n"

	readme, err := ioutil.ReadFile("README.md")
	if err != nil {
		log.Fatalln(err)
	}
	re := regexp.MustCompile("(?s)<!-- region Rules -->.*<!-- endregion -->")
	readme = re.ReplaceAll(readme, []byte(output))
	if err := ioutil.WriteFile("README.md", readme, 0644); err != nil {
		log.Fatalln(err)
	}
}
