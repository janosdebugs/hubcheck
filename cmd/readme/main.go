package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "regexp"

    "go.debugged.it/hubcheck/rules"
)

func main() {
    output := "<!-- region Rules -->\n\n"
    for _, rule := range rules.New() {
        output += fmt.Sprintf("### %s\n\n%s\n\nRead more: %s\n\n", rule.Name(), rule.Description(), rule.DocURL())
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
