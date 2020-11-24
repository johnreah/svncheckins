package main

import (
	"fmt"
	"time"
	"sort"
)

const svnHost = "https://<cloud-svn-host>/"

type Repository struct {
	name  string
	parentPom string
}

var svn Svn

var repos = []Repository{
	{"projectA", "pom.xml"},
	{"projectB", "pom.xml"},
	{"projectC", "<subdirectory>/pom.xml"},
}

type RepoResult struct {
	RepoName string
	TrunkInfo SvnResult
	TrunkPomInfo PomResult
	BranchesInfo SvnResultSlice
}

type SvnResult struct {
	Url string
	LastCommitDate string
	LastCommitAuthor string
}

type SvnResultSlice []SvnResult
func (s SvnResultSlice) Less(i, j int) bool { return s[i].LastCommitDate > s[j].LastCommitDate }
func (s SvnResultSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SvnResultSlice) Len() int           { return len(s) }

type PomResult struct {
	ArtifactId string
	Version string
}

func main() {
	fmt.Println("Starting...")
	startTime := time.Now()

	results := make(map[string]RepoResult)

	chanResults := make(chan RepoResult)
	for _, repo := range repos {
		go getRepoInfo(repo, chanResults)
	}

	for range repos {
		result := <-chanResults
		results[result.RepoName] = result
	}

	var sortedRepoNames []string
	for key := range results {
		sortedRepoNames = append(sortedRepoNames, key)
	}
	sort.Strings(sortedRepoNames)

	fmt.Printf("Showing results for %d repositories...\n", len(repos))
	for _, key := range sortedRepoNames {
		result := results[key]
		fmt.Printf("Repository: %s\n", result.RepoName)
		fmt.Printf("    URL: %s\n", result.TrunkInfo.Url)
		fmt.Printf("    Date: %s\n", result.TrunkInfo.LastCommitDate)
		fmt.Printf("    Author: %s\n", result.TrunkInfo.LastCommitAuthor)
		fmt.Printf("    ArtifactId: %s\n", result.TrunkPomInfo.ArtifactId)
		fmt.Printf("    Version: %s\n", result.TrunkPomInfo.Version)
		fmt.Printf("    Branches: %d\n", len(result.BranchesInfo))
		sort.Sort(result.BranchesInfo)
		for _, br := range result.BranchesInfo {
			fmt.Printf("        %s %s %s\n", br.LastCommitDate, br.Url, br.LastCommitAuthor)
		}
	}

	fmt.Printf("\nFinished in %dms\n", time.Now().Sub(startTime)/1000000)
}

func getRepoInfo(repo Repository, c chan RepoResult) {
	chanTrunkResult := make(chan SvnResult)
	go getTrunkInfo(repo, chanTrunkResult)

	chanPomResult := make(chan PomResult)
	go getPomInfo(repo, chanPomResult)

	chanBranchesResult := make(chan SvnResultSlice)
	go getBranchesInfo(repo, chanBranchesResult)

	trunkResult := <- chanTrunkResult
	pomResult := <- chanPomResult
	branchesResult := <- chanBranchesResult

	c <- RepoResult{repo.name, trunkResult, pomResult, branchesResult}
}

func getTrunkInfo(repo Repository, c chan SvnResult) {
	repoUrl := svnHost + repo.name
	svninfo, err := svn.Info(repoUrl, []string{})
	if err != nil {
		panic(err)
	}
	var svnResult SvnResult
	svnResult.Url = repoUrl
	svnResult.LastCommitDate = svninfo.Entry.Commit.Date
	svnResult.LastCommitAuthor = svninfo.Entry.Commit.Author
	c <- svnResult
}

func getPomInfo(repo Repository, c chan PomResult) {
	pomUrl := svnHost + repo.name + "/trunk/" + repo.parentPom
	pomProject, err := svn.GetPomInfo(pomUrl)
	if err != nil {
		fmt.Printf("Problem getting POM info: %v", err)
		c <- PomResult{}
		return
	}
	c <- PomResult{pomProject.ArtifactId, pomProject.Version}
}

func getBranchesInfo(repo Repository, c chan SvnResultSlice) {
	repoUrl := svnHost + repo.name + "/branches"
	svnList, err := svn.List(repoUrl, []string{})
	if err != nil {
		fmt.Printf("Problem listing branches: %v", err)
		c <- SvnResultSlice{}
		return
	}
	var svnResults SvnResultSlice
	for _, br := range svnList.Entries {
		svnResults = append(svnResults, SvnResult{
			Url: br.Name,
			LastCommitDate: br.Commit.Date,
			LastCommitAuthor:br.Commit.Author,
		})
	}
	c <- svnResults
}

