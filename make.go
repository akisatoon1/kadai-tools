package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

var ReNum, _ = regexp.Compile("^[0-9][0-9]$") // check lesson number
var ReFile, _ = regexp.Compile("^[a-z]$")     // check file name
var Dir string
var ExeDir = exePathDir()

type fileData struct {
	Num   string
	Level string
}

func exePathDir() string {
	path, _ := os.Executable()
	dir := filepath.Dir(path)
	return dir
}

func make(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("2つ以上の引数が必要です。")
	}

	// check arguments format
	for i, arg := range args {
		switch i {
		case 0:
			if !ReNum.MatchString(arg) {
				return fmt.Errorf("2桁の数字(0埋め)にしてください。第1引数: %v", arg)
			}
			Dir = arg
		default:
			if !ReFile.MatchString(arg) {
				return fmt.Errorf("英小文字1文字にしてください。第%v引数: %v", i+1, arg)
			}
		}
	}

	// make directory
	err := os.Mkdir("./"+Dir, 0777)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("既に./%v/が存在しています。\n", Dir)
			fmt.Printf("./%v/以下を全て書き換えますか? [y/n]: ", Dir)
			var yn string
			fmt.Scan(&yn)
			if yn != "y" {
				fmt.Println("書き換えを中止しました。")
				fmt.Println("終了")
				return nil
			}
			os.RemoveAll("./" + Dir)
			fmt.Printf("./%v/ 削除\n", Dir)
		} else {
			return err
		}
	}
	os.Mkdir("./"+Dir, 0777)
	fmt.Printf("./%v/ 作成\n", Dir)

	// make inputFiles directory
	inputDir := fmt.Sprintf("./%v/inputFiles", Dir)
	os.Mkdir(inputDir, 0777)

	// read template
	tmpl, tmplErr := template.ParseFiles(ExeDir + "/tmpl.c")
	if tmplErr == nil {
		fmt.Printf("./tmpl.c 読み込み\n")
	} else {
		fmt.Println("テンプレート無し")
	}

	// make files and input files
	levels := args[1:]
	for _, level := range levels {
		filename := fmt.Sprintf("./%v/kadai%v%v.c", Dir, Dir, level)
		fpKadai, _ := os.Create(filename)
		fmt.Println(filename + " 作成")
		defer fpKadai.Close()

		if tmplErr == nil {
			tmpl.Execute(fpKadai, fileData{Num: Dir, Level: level})
			fmt.Println(filename + " テンプレート書き込み")
		}

		inputFile := fmt.Sprintf("%v/input%v1.txt", inputDir, level)
		fpInput, _ := os.Create(inputFile)
		fmt.Println(inputFile + " 作成")
		defer fpInput.Close()
	}

	fmt.Println("完了")
	return nil
}
