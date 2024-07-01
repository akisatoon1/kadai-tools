package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/akisatoon1/manaba"
)

var username = os.Getenv("MANABA_ID")
var password = os.Getenv("MANABA_PASS")

func submit(args []string) error {
	// check args
	if len(args) != 1 {
		return fmt.Errorf("引数の数は1個です。")
	}
	if arg := args[0]; len(arg) != 1 || arg[0] < "a"[0] || "z"[0] < arg[0] {
		return fmt.Errorf("引数が間違っています")
	}

	// login
	jar, _ := cookiejar.New(nil)
	err := manaba.Login(jar, username, password)
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
	kadaiNum := getKadaiNum() // ex. kadai[01]a
	kadaiName := "kadai" + kadaiNum + args[0]
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

	// submit flow
	err = manaba.UploadFile(jar, url, kadaiName+".c")
	if err != nil {
		return err
	}
	err = manaba.SubmitReports(jar, url)
	if err != nil {
		return err
	}

	return nil
}
