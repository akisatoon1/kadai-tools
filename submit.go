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
	//
	// 不完全
	// submit -t -r のような形式
	// 指定した形式以外だった時の処理をしていない
	//

	// check args
	var objs []string
	var options []string
	for _, arg := range args {
		if arg[0:1] == "-" && len(arg[1:]) == 1 {
			options = append(options, arg[1:])
		} else if arg[0:1] != "-" && len(arg) == 1 {
			objs = append(objs, arg)
		} else {
			return fmt.Errorf("引数の形式が間違っています")
		}
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
	kadaiNum := getKadaiNum()                 // ex. kadai[01]a
	kadaiName := "kadai" + kadaiNum + objs[0] // 後で
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

	if len(options) == 1 && options[0] == "r" {
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

	// submit flow
	fileName := kadaiName + ".c"
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
	body, err := os.ReadFile(ExeDir + "/login.json")
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
