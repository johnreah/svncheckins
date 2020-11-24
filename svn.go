package main

import (
	"encoding/xml"
	"net/http"
	"io/ioutil"
	"os"
)

type SvnInfoEntryCommit struct {
	Revision string `xml:"revision,attr"`
	Author   string `xml:"author"`
	Date     string `xml:"date"`
}

type SvnInfoEntryRepo struct {
	Root string `xml:"root"`
	UUID string `xml:"uuid"`
}

type SvnInfoEntry struct {
	Revision    string             `xml:"revision,attr"`
	Path        string             `xml:"path,attr"`
	Kind        string             `xml:"kind,attr"`
	URL         string             `xml:"url"`
	RelativeURL string             `xml:"relative-url"`
	Repository  SvnInfoEntryRepo   `xml:"repository"`
	Commit      SvnInfoEntryCommit `xml:"commit"`
}

type SvnInfo struct {
	XMLName xml.Name `xml:"info"`
	Entry struct {
		Revision    string             `xml:"revision,attr"`
		Path        string             `xml:"path,attr"`
		Kind        string             `xml:"kind,attr"`
		URL         string             `xml:"url"`
		RelativeURL string             `xml:"relative-url"`
		Repository  SvnInfoEntryRepo   `xml:"repository"`
		Commit      SvnInfoEntryCommit `xml:"commit"`
	} `xml:"entry"`
}

type Svn struct{}

func (svn *Svn) Info(sourceUrl string, flags []string) (SvnInfo, error) {
	result := SvnInfo{}
	commandArgs := []string{"svn", "info"}
	commandArgs = append(commandArgs, sourceUrl)
	if len(flags) > 0 {
		for _, flag := range flags {
			commandArgs = append(commandArgs, flag)
		}
	}
	commandArgs = append(commandArgs, "--xml")
	xml_data, err := commandResult(commandArgs)
	if err != nil {
		return result, err
	}
	err = xml.Unmarshal(xml_data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type SvnLists struct {
	List SvnList `xml:"list"`
}

type SvnList struct {
	Entries []SvnListEntry `xml:"entry"`
}

type SvnListEntry struct {
	Name string `xml:"name"`
	Commit SvnListEntryCommit `xml:"commit"`
}

type SvnListEntryCommit struct {
	Revision string `xml:"revision,attr"`
	Author string `xml:"author"`
	Date string `xml:"date"`
}

func (svn *Svn) List(sourceUrl string, flags []string) (SvnList, error) {
	commandArgs := []string{"svn", "list"}
	commandArgs = append(commandArgs, sourceUrl)
	if len(flags) > 0 {
		for _, flag := range flags {
			commandArgs = append(commandArgs, flag)
		}
	}
	commandArgs = append(commandArgs, "--xml")
	xml_data, err := commandResult(commandArgs)
	if err != nil {
		return SvnList{}, err
	}
	var svnLists SvnLists
	err = xml.Unmarshal(xml_data, &svnLists)
	if err != nil {
		return SvnList{}, err
	}
	return svnLists.List, nil
}

type PomProject struct {
	ArtifactId string `xml:"artifactId"`
	Version string `xml:"version"`
}

func (svn *Svn) GetPomInfo(pomUrl string) (PomProject, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", pomUrl, nil)
	request.SetBasicAuth(os.Getenv("BEANSTALK_USERNAME"), os.Getenv("BEANSTALK_PASSWORD"))
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		var project PomProject
		if err = xml.Unmarshal(contents, &project); err != nil {
			panic(err)
		}
		return project, nil
	}
}
