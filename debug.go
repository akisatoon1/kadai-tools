package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DebugSet struct {
	Lang       string
	CFileName  string
	InputFiles []string
}

func (debugObj *DebugSet) storeInputFiles(kadaiNum string) {
	ext, _ := getFileExt(debugObj.Lang)
	kadaiLevel := strings.TrimPrefix(debugObj.CFileName, fmt.Sprintf("kadai%v", kadaiNum))
	kadaiLevel = strings.TrimSuffix(kadaiLevel, "."+ext)
	debugObj.InputFiles, _ = filepath.Glob(fmt.Sprintf("./inputFiles/input%v[0-9].txt", kadaiLevel))
}

func (debugObj *DebugSet) debug() error {
	// apply lang option
	ext, _ := getFileExt(debugObj.Lang)
	compiler, _ := getCompilerName(debugObj.Lang)

	executable := strings.TrimSuffix(debugObj.CFileName, "."+ext)

	// compile
	cmdCompile := exec.Command(compiler, "-Wall", "-o", executable, debugObj.CFileName, "-lm")
	cmdCompile.Stdout = os.Stdout
	cmdCompile.Stderr = os.Stderr
	fmt.Printf("compile: %v -Wall -o %v %v -lm\n", compiler, executable, debugObj.CFileName)
	err := cmdCompile.Run()
	if err != nil {
		return err
	}

	// execute
	for _, file := range debugObj.InputFiles {
		fp, err := os.Open("./" + file)
		if err != nil {
			return err
		}
		defer fp.Close()

		cmdExe := exec.Command("./" + executable)
		cmdExe.Stdin = fp
		cmdExe.Stdout = os.Stdout
		cmdExe.Stderr = os.Stderr
		fmt.Printf("Output of %v:\n", file)
		err = cmdExe.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func debug(args []string) error {
	// get lang option
	lang, err := getLang()
	if err != nil {
		return err
	}
	ext, _ := getFileExt(lang)

	kadaiNum := getKadaiNum() // ex. kadai[01]a
	var debugObjs []DebugSet

	// if no args, debug all c files.
	if len(args) == 0 {
		cFiles, _ := filepath.Glob(fmt.Sprintf("kadai%v[a-z].%v", kadaiNum, ext))
		for _, file := range cFiles {
			debugObjs = append(debugObjs, DebugSet{Lang: lang, CFileName: file})
		}
	} else {
		for _, arg := range args {
			cFiles, _ := filepath.Glob(fmt.Sprintf("kadai%v%v.%v", kadaiNum, arg, ext))
			// if no cFiles, ignore it.
			if len(cFiles) == 0 {
				fmt.Printf("kadai%v%v.%vは存在しません\n", kadaiNum, arg, ext)
				continue
			}
			debugObjs = append(debugObjs, DebugSet{Lang: lang, CFileName: cFiles[0]})
		}
	}

	var tmp []DebugSet
	for _, debugObj := range debugObjs {
		debugObj.storeInputFiles(kadaiNum)
		tmp = append(tmp, debugObj)
	}
	debugObjs = tmp
	for _, debugObj := range debugObjs {
		err := debugObj.debug()
		if err != nil {
			return err
		}
	}

	return nil
}
