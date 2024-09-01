package main

/*
 *    Copyright 2024 YggdrasilSoft, LLC
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at

 *      http://www.apache.org/licenses/LICENSE-2.0

 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * ---
 * Author: Gary L. Greene, Jr. <greeneg@yggdrasilsoft.com>
 */

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	fqdn "github.com/Showmax/go-fqdn"
	"github.com/matishsiao/goInfo"
	//	"github.com/pborman/getopt/v2"
	//	"golang.org/x/term"
)

type osReleaseStruct struct {
	Id      string
	Version string
}

type UpdateStruct struct {
	Updates     []Update `json:"updates"`
	UpdateCount int      `json:"updateCount"`
	FQDN        string   `json:"fqdn"`
	OsFamily    string   `json:"osFamily"`
	OsId        string   `json:"osId"`
	OsVersion   string   `json:"osVersion"`
	HostArch    string   `json:"hostArchitecture"`
}

type Update struct {
	Kind       string `json:"kind" xml:"kind,attr"`
	Name       string `json:"name" xml:"name,attr"`
	Version    string `json:"version" xml:"edition,attr"`
	Arch       string `json:"arch" xml:"arch,attr"`
	OldVersion string `json:"oldVersion" xml:"edition-old,attr"`
	Summary    string `json:"summary" xml:"summary"`
}

type UpdateList struct {
	Updates []Update `xml:"update"`
}

type UpdateStatus struct {
	UpdateList UpdateList `xml:"update-list"`
}

type Stream struct {
	UpdateStatus UpdateStatus `xml:"update-status"`
}

func getOsId(content string) (string, error) {
	var osId string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			osId = strings.Split(line, "=")[1] // get back the actual ID value
			osId = strings.Trim(osId, "\"")
		}
	}

	return osId, nil
}

func getOsVersion(content string) (string, error) {
	var osVersion string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VERSION_ID=") {
			osVersion = strings.Split(line, "=")[1] // get back the actual ID value
			osVersion = strings.Trim(osVersion, "\"")
		}
	}

	return osVersion, nil
}

func ReadOsReleaseFile(releaseFile string, r *osReleaseStruct) error {
	fileContent, err := os.ReadFile(releaseFile)
	if err != nil {
		fatalCheckError(err)
	}

	// line by line parse output into struct entries
	r.Id, err = getOsId(string(fileContent))
	if err != nil {
		return err
	}
	r.Version, err = getOsVersion(string(fileContent))
	if err != nil {
		return err
	}

	return nil
}

func DetectOs() (osReleaseStruct, string, error) {
	var osVariant string
	ors := osReleaseStruct{}
	osReleaseFile := "/etc/os-release"
	if _, err := os.Stat(osReleaseFile); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return ors, "", err
		}
	} else {
		err := ReadOsReleaseFile(osReleaseFile, &ors)
		if err != nil {
			return ors, "", err
		}
	}

	if ors.Id == "opensuse-leap" {
		osVariant = "suse"
	} else if ors.Id == "ubuntu" {
		osVariant = "debian"
	}

	return ors, osVariant, nil
}

func processAptOutput(text string) (string, error) {
	var output string

	// drop the unneeded first line
	text = strings.Replace(text, "Listing...\n", "", -1)
	log.Fatal(text)

	return output, nil
}

func getAptListOutput() (string, error) {
	var output string

	out, err := exec.Command("apt", "list", "--upgradeable").Output()
	if err != nil {
		return "", err
	}

	output = string(out)
	return output, nil
}

func processZypperOutput(text string) (Stream, error) {
	s := Stream{}
	err := xml.Unmarshal([]byte(text), &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func getZypperLuOutput() (string, error) {
	var output string

	out, err := exec.Command("zypper", "--no-color", "--no-refresh", "-x", "lu").Output()
	if err != nil {
		return "", err
	}

	output = string(out)
	return output, nil
}

func processUpdates(output string, osFamily string) (UpdateStruct, error) {
	us := UpdateStruct{}

	count := 0
	if osFamily == "suse" {
		s, err := processZypperOutput(output)
		if err != nil {
			return us, err
		}
		updates := make([]Update, 0)
		for _, pkg := range s.UpdateStatus.UpdateList.Updates {
			count++
			u := Update{}
			u.Kind = pkg.Kind
			u.Name = pkg.Name
			u.Version = pkg.Version
			u.Arch = pkg.Arch
			u.OldVersion = pkg.OldVersion
			u.Summary = pkg.Summary

			updates = append(updates, u)
		}
		us.Updates = updates
	} else if osFamily == "debian" {
		_, err := processAptOutput(output)
		if err != nil {
			return us, err
		}
	}

	us.UpdateCount = count

	return us, nil
}

func getOsFamily() (string, error) {
	gi, err := goInfo.GetInfo()
	if err != nil {
		return "", err
	}

	return gi.GoOS, nil
}

func getHostArch() (string, error) {
	gi, err := goInfo.GetInfo()
	if err != nil {
		return "", err
	}
	return gi.Platform, nil
}

func main() {
	output := ""
	ors, osVariant, err := DetectOs()
	if err != nil {
		os.Exit(1)
	}

	if osVariant == "suse" {
		output, err = getZypperLuOutput()
		if err != nil {
			log.Fatal(err)
		}
	} else if osVariant == "debian" {
		output, err = getAptListOutput()
		if err != nil {
			log.Fatal(err)
		}
	}

	us, err := processUpdates(output, osVariant)
	if err != nil {
		log.Fatal(err)
	}

	us.OsId = ors.Id
	us.OsVersion = ors.Version

	fqdn, err := fqdn.FqdnHostname()
	if err != nil {
		log.Fatal(err)
	}
	us.FQDN = fqdn

	osFamily, err := getOsFamily()
	if err != nil {
		log.Fatal(err)
	}
	us.OsFamily = osFamily

	hostArch, err := getHostArch()
	if err != nil {
		log.Fatal(err)
	}
	us.HostArch = hostArch

	// convert the UpdateStruct to JSON text
	j, err := json.Marshal(us)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(os.Stdout, string(j))
}
