package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	fc "github.com/fatih/color"
	"github.com/kowshikRoy/cft/model"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func ContestDir(contestID int) string {
	woringDir := viper.GetString("workdir")
	contestDir := path.Join(woringDir, strconv.Itoa(contestID))
	return contestDir
}

func TestDir(contestID int) string {
	contestDir := ContestDir(contestID)
	testDir := path.Join(contestDir, "tests")
	return testDir
}

func BinDir(contestID int) string {
	contestDir := ContestDir(contestID)
	binDir := path.Join(contestDir, "bin")
	return binDir
}
func CreateDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("Couldn't create directory: %v", err)
	}
	return nil
}

func CreateDev(contestID int) (string, string, string) {
	contestDir := ContestDir(contestID)
	testDir := TestDir(contestID)
	binDir := BinDir(contestID)
	CreateDir(contestDir)
	CreateDir(testDir)
	CreateDir(binDir)
	return contestDir, testDir, binDir
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

func GetSourceFileExtension(lang string) string {
	l, ok := Mapping[lang]
	if !ok {
		fmt.Println("Your language is not supported or you configured language options incorrectly")
		os.Exit(1)
	}
	return l.Extension
}

func santize(lang string) string {
	return strings.ToLower(lang)
}

func BuildFile(contestDir, problem, lang string) error {
	lang = santize(lang)
	temp := viper.GetStringMap("buildConfig")
	lo, ok := temp[lang]
	if !ok {
		fmt.Println("Your language is not supported or you configured language options incorrectly")
		os.Exit(1)
	}
	var l model.Language
	if mapstructure.Decode(lo, &l) != nil {
		fmt.Println("Couldn't parse your language buildconfig")
		os.Exit(1)
	}

	buildString := path.Join(contestDir, problem+l.Extension)
	outputFile := path.Join(contestDir, "bin", problem+l.OutputFileExtension[runtime.GOOS])
	args := strings.Split(l.BuildFlags, " ")
	args = append(args, buildString, "-o", outputFile)

	buildC := exec.Command(l.Compiler[runtime.GOOS], args...)
	out, err := buildC.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		os.Exit(1)
	}
	fmt.Println(fc.GreenString("Compiled Successfully!!"))
	return nil
}

func split(r rune) bool { return r == ' ' || r == '\n' || r == '\r' }
func cmpFunc(a, b string) bool {
	if strings.Compare(a, b) == 0 {
		return true
	}
	err := viper.GetFloat64("margin")
	A, e1 := strconv.ParseFloat(a, 64)
	B, e2 := strconv.ParseFloat(b, 64)
	if e1 == nil && e2 == nil && (B-A)/B < err {
		return true
	}
	return false
}
func cmp(t1, t2 []byte) bool {
	str1 := string(t1)
	str2 := string(t2)
	p1 := strings.FieldsFunc(str1, split)
	p2 := strings.FieldsFunc(str2, split)

	if len(p1) != len(p2) {
		return false
	}

	for i := 0; i < len(p1); i++ {
		if cmpFunc(p1[i], p2[i]) == false {
			return false
		}
	}
	return true
}
func run(c chan string, binFile, inFile, outFile string) {
	bytes, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	runC := exec.Command(binFile)
	stdin, err := runC.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.WriteString(stdin, string(bytes))
	if err != nil {
		log.Fatal(err)
	}
	beg := time.Now()
	out, err := runC.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := "-----------------------------------\n"
	result += fmt.Sprintf("Time taken: %s ", fc.HiMagentaString(time.Since(beg).String()))
	outf, err := ioutil.ReadFile(outFile)
	if err != nil {
		log.Fatalf("Couldn't read %s: %v", outFile, err)
	}
	ok := cmp(out, outf)
	if ok {
		result += fmt.Sprintln("Status:", fc.GreenString("OK"))
	} else {
		result += fmt.Sprintln("Status:", fc.HiRedString("WA"))
		result += fmt.Sprintln(fc.YellowString("Your Output:"))
		result += fmt.Sprintln(fc.HiBlackString(string(out)))
		result += fmt.Sprintln(fc.HiBlueString("Expected Answer:"))
		result += fmt.Sprintln(fc.HiMagentaString(string(outf)))
	}
	c <- result

}

func runSingleTest(ctx context.Context, binFile, inFile, outFile string) {
	c := make(chan string)
	beg := time.Now()
	go run(c, binFile, inFile, outFile)
	select {
	case <-ctx.Done():
		fmt.Println("Time Elapsed:", fc.CyanString(time.Since(beg).String()), "Status:", fc.HiRedString("TLE"))
	case out := <-c:
		fmt.Print(out)
	}
}

func RunTest(contestDir, problem string, param model.ExecParam) error {
	lang := santize(param.Lang)
	l, ok := Mapping[lang]
	if !ok {
		fmt.Println("Your language is not supported or you configured language options incorrectly")
		os.Exit(1)
	}
	testFileDir := path.Join(contestDir, "tests", problem)
	binFile := path.Join(contestDir, "bin", problem+l.OutputFileExtension[runtime.GOOS])
	files, _ := ioutil.ReadDir(testFileDir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "input") {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, time.Duration(param.TimeLimit)*time.Second)
			defer cancel()
			runSingleTest(ctx,
				binFile,
				path.Join(testFileDir, f.Name()),
				path.Join(testFileDir, strings.ReplaceAll(f.Name(), "input", "output")))

		}
	}
	return nil
}
