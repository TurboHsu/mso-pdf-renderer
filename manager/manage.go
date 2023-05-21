package manager

import (
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Routines []RoutineStruct
var RunningPath string

func init() {
	// Makes the map
	Routines = make([]RoutineStruct, 0)
}

func FindRoutine(uuid string) *RoutineStruct {
	for i := 0; i < len(Routines); i++ {
		if uuid == Routines[i].UUID {
			return &Routines[i]
		}
	}
	return nil
}

func DoesUUIDExist(uuid string) bool {
	for _, routine := range Routines {
		if uuid == routine.UUID {
			return true
		}
	}
	return false
}

func GenerateUUID() string {
	// Generate random uuid
	return uuid.New().String()
}

func RemoveUUID(uuid string) {
	// If any file related to this uuid exists, remove it
	relatedFiles, err := filepath.Glob(RunningPath + "/cache/" + uuid + ".*")
	if err != nil {
		log.Println("[E] failed to check files with uuid ", uuid, ", error: ", err)
	}

	for _, file := range relatedFiles {
		if err := os.Remove(file); err != nil {
			log.Println("[E] failed to remove file ", file, ", error: ", err)
		}
	}

	// Remove specific uuid
	for i, routine := range Routines {
		if routine.UUID == uuid {
			Routines = append(Routines[:i], Routines[i+1:]...)
			break
		}
	}
}

func CheckExtensionValidation(extension string) (string, bool) {
	// Check whether extension is valid
	switch extension {
	case ".ppt", ".pptx", ".doc", ".docx", ".xls", ".xlsx":
		return "mso", true
	default:
		return "", false
	}
}

func CacheLifeCycleRoutine(lifecycle int64, interval int64) {
	// Event loop
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		// Check whether cache is expired
		for _, routine := range Routines {
			if time.Now().Unix()-routine.LifeCycleStart > lifecycle {
				log.Println("[I] removing expired uuid: ", routine.UUID)
				RemoveUUID(routine.UUID)
			}
		}
	}
}
