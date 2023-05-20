package process

import (
	"github.com/TurboHsu/mso-pdf-renderer/manager"
	"log"
	"os"
	"os/exec"
)

var RunningPath string

func init() {
	RunningPath, _ = os.Getwd()
	// Judge whether script file exists
	if _, err := os.Stat(RunningPath + "/scripts/mso-convert.vbs"); os.IsNotExist(err) {
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
	extensionType, validation := manager.CheckExtensionValidation(routine.FileExtension)
	if validation {
		switch extensionType {
		case "mso":
			convertMSO(RunningPath+"/cache/"+uuid+routine.FileExtension, RunningPath+"/cache/"+uuid+".pdf")
		}
	}
}

func convertMSO(originalFilePath string, pdfPath string) {
	log.Println("[I] converting: ", originalFilePath)
	// Check whether originalFilePath exists
	if _, err := os.Stat(originalFilePath); os.IsNotExist(err) {
		log.Println("[E] ppt file not found: ", originalFilePath)
		return
	}

	// Execute script
	cmd := exec.Command("cscript", RunningPath+"/scripts/mso-convert.vbs", originalFilePath)
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
	err = os.Remove(originalFilePath)
	if err != nil {
		log.Println("[E] failed to delete original file: ", err)
	}
}
