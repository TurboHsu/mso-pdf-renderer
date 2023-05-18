package process

import (
	"log"
	"os"
	"os/exec"
)

var RunningPath string

func init() {
	RunningPath, _ = os.Getwd()
	// Judge whether script file exists
	if _, err := os.Stat(RunningPath + "/scripts/ppt.vbs"); os.IsNotExist(err) {
		panic("script file not found, exiting")
	}

	// Check whether cache folder exists
	if _, err := os.Stat(RunningPath + "/cache"); os.IsNotExist(err) {
		os.Mkdir(RunningPath+"/cache", os.ModePerm)
	}
}

func ConvertPPT(pptPath string, pdfPath string) {
	// Check whether pptPath exists
	if _, err := os.Stat(pptPath); os.IsNotExist(err) {
		log.Println("[E] ppt file not found: ", pptPath)
		return
	}

	// Execute script
	cmd := exec.Command("cscript", RunningPath+"/scripts/ppt.vbs", pptPath, pdfPath)
	err := cmd.Run()
	if err != nil {
		log.Println("[E] failed to execute script: ", err)
	}

	// Check whether pdfPath exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Println("[E] pdf file not found: ", pdfPath)
		return
	}

	// Delete ppt file
	err = os.Remove(pptPath)
	if err != nil {
		log.Println("[E] failed to delete ppt file: ", err)
	}
	
}
