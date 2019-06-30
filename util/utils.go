package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kowshikRoy/cft/model"
)

func CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("Couldn't create directory: %v", err)
	}
	return nil
}

func Crawl(contestID int, standings *model.Standings) error {
	url := "https://codeforces.com/api/contest.standings?from=1&count=1&contestId=" + strconv.Itoa(contestID)
	res, err := http.Get(url)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	out, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(out, standings)
	if err != nil {
		fmt.Printf("%v", err)
		return fmt.Errorf("Coundn't Parse JSON: %v", err)
	}
	if standings.Status != "OK" {
		return fmt.Errorf(standings.Comment)
	}
	return nil
}

func CrawlTest(c chan bool, testDir, index, url string) {

	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Errorf("Couldn't Fetch the testfile for %s: %v", url, err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	folder := path.Join(testDir, index)
	CreateDir(folder)
	doc.Find(".sample-test .input pre").Each(func(i int, s *goquery.Selection) {
		temp, _ := s.Html()
		text := strings.ReplaceAll(temp, "<br/>", "\n")
		fileName := path.Join(folder, "input-"+strconv.Itoa(i)+".txt")
		ioutil.WriteFile(fileName, []byte(text), 0755)
	})
	doc.Find(".sample-test .output pre").Each(func(i int, s *goquery.Selection) {
		temp, _ := s.Html()
		text := strings.ReplaceAll(temp, "<br/>", "\n")
		fileName := path.Join(folder, "output-"+strconv.Itoa(i)+".txt")
		ioutil.WriteFile(fileName, []byte(text), 0755)
	})
	c <- true
}
