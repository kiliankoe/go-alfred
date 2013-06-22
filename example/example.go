package main

import (
	"os"
	"github.com/pascalw/go-alfred"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func main() {
	queryTerms := os.Args[1:]

	// optimize query terms for fuzzy matching
	alfred.InitTerms(queryTerms)

	// create a new alfred workflow response
	response := alfred.NewResponse()
	repos := getRepos()

	for _, repo := range repos {
		// check if the repo name fuzzy matches the query terms
		if ! alfred.MatchesTerms(queryTerms, repo.Name) { continue }

		// it matched so add a new response item
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: true,
			Uid: repo.URL,
			Title: repo.Name,
			Arg: repo.URL,
		})
	}

	// finally print the resulting Alfred Workflow XML
	response.Print()
}

type Repo struct {
	Name string
	URL string
}

func getRepos() []Repo {
	resp, err := http.Get("https://api.github.com/users/pascalw/repos")
	if err != nil { panic(err) }

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var repos []Repo
	json.Unmarshal(body, &repos)

	return repos
}
