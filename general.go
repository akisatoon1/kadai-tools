package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var SettingsFile string = "settings.json"

func getKadaiNum() string {
	curDir, _ := os.Getwd()
	return filepath.Base(curDir)
}

func getLang() (string, error) {
	body, err := os.ReadFile(ExeDir + "/" + SettingsFile)
	if err != nil {
		return "", err
	}
	data := struct {
		Lang string `json:"lang"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if data.Lang != "c" && data.Lang != "c++" {
		return "", fmt.Errorf("対応しているlangは'c'または'c++'のみです。現在のlang: '%v'", data.Lang)
	}
	return data.Lang, nil
}
