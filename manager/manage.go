package manager

import "github.com/google/uuid"

var Routines []RoutineStruct

func init() {
	// Makes the map
	Routines = make([]RoutineStruct, 0)
}

func FindRoutine(uuid string) *RoutineStruct {
	for _, routine := range Routines {
		if uuid == routine.UUID {
			return &routine
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
	// Remove specific uuid
	for i, routine := range Routines {
		if routine.UUID == uuid {
			Routines = append(Routines[:i], Routines[i+1:]...)
			break
		}
	}
}
