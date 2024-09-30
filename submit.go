package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/akisatoon1/manaba"
)

func submit(args []string) error {
	// handle args
	var isResubmit bool
	var kadaiLevel string
	if l := len(args); l == 1 {
		if !ReFile.MatchString(args[0]) {
			return fmt.Errorf("引数が間違っています")
		}
		kadaiLevel = args[0]
		isResubmit = false
	} else if l == 2 {
		if args[0] != "-r" {
			return fmt.Errorf("引数が間違っています")
		}
		if !ReFile.MatchString(args[1]) {
			return fmt.Errorf("第2引数は英小文字1文字のみです")
		}
		kadaiLevel = args[1]
		isResubmit = true
	} else {
		return fmt.Errorf("引数の数は1つまたは2つです")
	}

	username, password, err := getUsernameAndPasswd()
	if err != nil {
		return err
	}

	// login
	jar, _ := cookiejar.New(nil)
	err = manaba.Login(jar, username, password)
	if err != nil {
		return err
	}

	// get url to submit
	client := &http.Client{Jar: jar}
	res, err := client.Get("https://room.chuo-u.ac.jp/ct/course_4932224_report") // URL of 'Cプロ演習'
	if err != nil {
		return err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil && err != io.EOF {
		return err
	}
	kadaiNum := getKadaiNum()                    // ex. kadai[01]a
	kadaiName := "kadai" + kadaiNum + kadaiLevel // 後で
	reg := regexp.MustCompile(kadaiName)
	var a *goquery.Selection
	doc.Find("h3.report-title").Find("a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if reg.MatchString(s.Text()) {
			a = s
			return false
		}
		return true
	})
	if a == nil {
		return fmt.Errorf(kadaiName)
	}
	url, _ := a.Attr("href")
	url = "https://room.chuo-u.ac.jp/ct/" + url

	if isResubmit {
		err = manaba.CancelSubmittion(jar, url)
		if err != nil {
			return err
		}
		err = manaba.DeleteAllFiles(jar, url)
		if err != nil {
			return err
		}
		fmt.Printf("'%v'のレポートを提出取り消ししました\n", kadaiName)
	}

	// determine file extension
	lang, err := getLang()
	if err != nil {
		return err
	}
	var ext string
	if lang == "c" {
		ext = "c"
	} else if lang == "c++" {
		ext = "cpp"
	}

	// submit flow
	fileName := kadaiName + "." + ext
	err = manaba.UploadFile(jar, url, fileName)
	if err != nil {
		return err
	}
	err = manaba.SubmitReports(jar, url)
	if err != nil {
		return err
	}

	fmt.Printf("'%v'を提出しました\n", fileName)
	return nil
}

func getUsernameAndPasswd() (string, string, error) {
	body, err := os.ReadFile(ExeDir + "/" + SettingsFile)
	if err != nil {
		return "", "", err
	}
	userData := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return "", "", err
	}
	return userData.Username, userData.Password, nil
}
