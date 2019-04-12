package composer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/knqyf263/trivy/pkg/git"
	"gopkg.in/yaml.v2"
)

type AdvisoryDB map[string][]Advisory

const (
	repoPath = "/tmp/composer"
	dbURL    = "https://github.com/FriendsOfPHP/security-advisories"
)

type Advisory struct {
	Cve       string
	Title     string
	Link      string
	Reference string
	Branches  map[string]Branch
}

type Branch struct {
	Versions []string
}

func (c *Scanner) UpdateDB() (err error) {
	if err := git.CloneOrPull(dbURL, repoPath); err != nil {
		return err
	}
	c.db, err = walk()
	return err
}

func walk() (AdvisoryDB, error) {
	advisoryDB := AdvisoryDB{}
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasPrefix(info.Name(), "CVE-") {
			return nil
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		advisory := Advisory{}
		err = yaml.Unmarshal(buf, &advisory)
		if err != nil {
			return err
		}
		advisories, ok := advisoryDB[advisory.Reference]
		if !ok {
			advisories = []Advisory{}
		}
		advisoryDB[advisory.Reference] = append(advisories, advisory)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return advisoryDB, nil
}
