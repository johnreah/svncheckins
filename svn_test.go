package main

import (
	"testing"
	"encoding/xml"
	"io/ioutil"
	"fmt"
)

type TProject struct {
	ArtifactId string `xml:"artifactId"`
	Version string `xml:"version"`
}

func TestPom(t *testing.T) {
	xmlFile, err := ioutil.ReadFile("pom.xml")
	if err != nil {
		panic(err)
	}
	var project TProject
	xml.Unmarshal(xmlFile, &project)
	fmt.Printf("%s\n", project.ArtifactId)
	fmt.Printf("%s\n", project.Version)
}

func TestList(t *testing.T) {
	var svn Svn
	const repoUrl = "https://<cloud-svn-host>/projectA/branches"
	svnList, err := svn.List(repoUrl, []string{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Repo %s has %d subdirs\n", repoUrl, len(svnList.Entries))
	for _, br := range svnList.Entries {
		fmt.Printf("%s %s %s\n", br.Name, br.Commit.Author, br.Commit.Date)
	}
}

