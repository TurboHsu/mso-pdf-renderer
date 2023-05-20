package process

import (
	"log"
	"mso-pdf-renderer/manager"
	"os"
	"os/exec"
)

var RunningPath string

func init() {
	RunningPath, _ = os.Getwd()
	// Judge whether script file exists
	if _, err := os.Stat(RunningPath + "/scripts/ppt.vbs"); os.IsNotExist(err) {
		log.Panicln("[E] script file not found, exiting")
	}

	// Clear cache
	os.RemoveAll(RunningPath + "/cache")
	os.Mkdir(RunningPath+"/cache", os.ModePerm)
}

func Convert(uuid string) {
	// Find routine
	routine := manager.FindRoutine(uuid)
	if routine == nil {
		log.Println("[E] invalid uuid: ", uuid)
		return
	}

	// Check whether file exists
	if _, err := os.Stat(RunningPath + "/cache/" + uuid + routine.FileExtension); os.IsNotExist(err) {
		log.Println("[E] file not found: ", RunningPath+"/cache/"+uuid+routine.FileExtension)
		return
	}

	// Convert
	switch routine.FileExtension {
	case ".ppt", ".pptx":
		convertPPT(RunningPath+"/cache/"+uuid+routine.FileExtension, RunningPath+"/cache/"+uuid+".pdf")
	}
}

func convertPPT(pptPath string, pdfPath string) {
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
