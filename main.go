// Copyright 2022 Kirill Krasnov <kirill@kraeg.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Github-stars-backup application save your starred github repository to local disk
//
// Application parameters:
//
//   -users  <[user-or-organisation-comma-separated-list]>
//   -output [local-folder-name], default: ./repos
//   -maxrepo [maximum-repo-to-be-clone], default: 1000
//
// Usage examples:
//
//   go run . -users=kirill -output=./repo
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {

	// Parse parameters
	var userslist, output, maxrepo string
	flag.StringVar(&userslist, "users", "", "user or organisation comma separated list")
	flag.StringVar(&output, "output", "repos", "local folder name to save repositories")
	flag.StringVar(&maxrepo, "maxrepo", "1000", "maximum number of users repositories to be cloned")
	flag.Parse()

	// Parse users and limit
	users := strings.Split(userslist, ",")

	// Get list of repos with gh cli application
	for _, user := range users {
		getRepos(output, strings.TrimSpace(user), maxrepo)
	}
}

// Number of repositories to show in print
var reponum int

// getRepos get list of reopsitories and clone it
func getRepos(dir, user, maxrepo string) (repos []string) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/starred?page=1&per_page=%s", user, maxrepo))
	if err != nil {
		if err != nil {
			log.Printf("Can't get starred repos of %s: %s", user, err)
			return nil
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Can't read response body: %s", err)
		return nil
	}
	// fmt.Printf("The data is %s\n", string(body))

	type repoData struct {
		Clone       string `json:"clone_url,omitempty"`
		Description string `json:"description,omitempty"`
		FullName    string `json:"full_name,omitempty"`
		HasWiki     bool   `json:"has_wiki,omitempty"`
		SSHUrl      string `json:"ssh_url,omitempty"`
	}
	var jsonData []repoData

	if err := json.Unmarshal(body, &jsonData); err != nil {
		log.Printf("Can't parse response body to json: %s\n%s", err, string(body))
		return nil
	}

	for _, r := range jsonData {

		// Print repo name
		reponum++
		fmt.Printf("repo %3d: %s\n", reponum, r.FullName)
		repos = append(repos, r.FullName)

		// Clone repo
		_, err := exec.Command("git", "clone", "--mirror", r.SSHUrl, dir+"/"+r.FullName+".git").Output()
		if err != nil {
			log.Fatal(err)
		}

		if r.HasWiki {
			// Clone wiki repo
			err = exec.Command("git", "clone", "--mirror", "git@github.com:"+r.FullName+".wiki.git", dir+"/"+r.FullName+".wiki.git").Run()
			if err != nil {
				// log.Println(err)
			}
		}
	}
	return
}
