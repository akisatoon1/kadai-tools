package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DebugSet struct {
	CFileName  string
	InputFiles []string
}

func (debugObj *DebugSet) storeInputFiles(kadaiNum string) {
	kadaiLevel := strings.TrimPrefix(debugObj.CFileName, fmt.Sprintf("kadai%v", kadaiNum))
	kadaiLevel = strings.TrimSuffix(kadaiLevel, ".c")
	debugObj.InputFiles, _ = filepath.Glob(fmt.Sprintf("./inputFiles/input%v[0-9].txt", kadaiLevel))
}

func (debugObj *DebugSet) debug() error {
	executable := strings.TrimSuffix(debugObj.CFileName, ".c")

	// compile
	cmdCompile := exec.Command("gcc", "-Wall", "-o", executable, debugObj.CFileName, "-lm")
	cmdCompile.Stdout = os.Stdout
	cmdCompile.Stderr = os.Stderr
	fmt.Printf("compile: gcc -Wall -o %v %v -lm\n", executable, debugObj.CFileName)
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
	kadaiNum := getKadaiNum() // ex. kadai[01]a
	var debugObjs []DebugSet

	// if no args, debug all c files.
	if len(args) == 0 {
		cFiles, _ := filepath.Glob(fmt.Sprintf("kadai%v[a-z].c", kadaiNum))
		for _, file := range cFiles {
			debugObjs = append(debugObjs, DebugSet{CFileName: file})
		}
	} else {
		for _, arg := range args {
			cFiles, _ := filepath.Glob(fmt.Sprintf("kadai%v%v.c", kadaiNum, arg))
			// if no cFiles, ignore it.
			if len(cFiles) == 0 {
				fmt.Printf("kadai%v%v.cは存在しません\n", kadaiNum, arg)
				continue
			}
			debugObjs = append(debugObjs, DebugSet{CFileName: cFiles[0]})
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
