package server

import (
	"encoding/json"
	"github.com/TurboHsu/mso-pdf-renderer/manager"
	"github.com/TurboHsu/mso-pdf-renderer/process"
	"io"
	"log"
	"net/http"
	"os"
)

var downloadSemaphore = make(chan struct{}, 1)

func init() {
	// Check whether /static folder exists
	if _, err := os.Stat(process.RunningPath + "/static"); os.IsNotExist(err) {
		log.Panicln("[E] /static folder not found, exiting")
	}
}

func ListenAndServe(addr string) {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/check", checkHandler)
	http.HandleFunc("/", http.FileServer(http.Dir(process.RunningPath+"/static")).ServeHTTP)
	log.Println("[I] listening on ", addr)
	http.ListenAndServe(addr, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var resp APIResponseStruct
	r.ParseForm()
	extension := r.FormValue("extension")
	_, extensionValidation := manager.CheckExtensionValidation("." + extension)
	if !extensionValidation {
		resp.Status = "bad"
		resp.Message = "Invalid extension"
		w.Write(marshalResponse(resp))
		return
	}
	uuid := manager.GenerateUUID()
	manager.Routines = append(manager.Routines, manager.RoutineStruct{
		UUID:          uuid,
		FileExtension: "." + extension,
	})
	resp.Status = "ok"
	resp.Message = uuid
	w.Write(marshalResponse(resp))
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	var resp APIResponseStruct
	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		resp.Status = "bad"
		resp.Message = "Invalid uuid"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(marshalResponse(resp))
		return
	}

	if _, err := os.Stat(process.RunningPath + "/cache/" + uuid + manager.FindRoutine(uuid).FileExtension); !os.IsNotExist(err) {
		resp.Status = "wait"
		resp.Message = "Not done yet"
		w.WriteHeader(http.StatusAccepted)
		w.Write(marshalResponse(resp))
		return
	}

	resp.Status = "ok"
	resp.Message = "done"
	w.Write(marshalResponse(resp))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	var resp APIResponseStruct
	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		resp.Status = "bad"
		resp.Message = "Invalid uuid"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(marshalResponse(resp))
		return
	}

	filename := uuid + ".pdf"

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Limit the number of concurrent downloads to one
	downloadSemaphore <- struct{}{}
	defer func() { <-downloadSemaphore }()

	ctx := r.Context()
	go func() {
		<-ctx.Done()
		// Do something when the client has disconnected
		manager.RemoveUUID(uuid)
		os.Remove(process.RunningPath + "/cache/" + filename)
	}()

	http.ServeFile(w, r, process.RunningPath+"/cache/"+filename)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var resp APIResponseStruct
	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		resp.Status = "bad"
		resp.Message = "Invalid uuid"
		w.WriteHeader(http.StatusBadRequest)
		w.Write(marshalResponse(resp))
		return
	}

	extension := manager.FindRoutine(uuid).FileExtension
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("[E] failed to upload file: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp.Status = "bad"
		resp.Message = "Upload failed: " + err.Error()
		w.Write(marshalResponse(resp))
		return
	}
	defer file.Close()

	fileLocal, err := os.Create(process.RunningPath + "/cache/" + uuid + extension)
	if err != nil {
		resp.Status = "bad"
		resp.Message = "Upload failed: " + err.Error()
		w.Write(marshalResponse(resp))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[E] failed to create file: ", err)
		return
	}
	defer fileLocal.Close()

	_, err = io.Copy(fileLocal, file)
	if err != nil {
		resp.Status = "bad"
		resp.Message = "Upload failed: " + err.Error()
		w.Write(marshalResponse(resp))
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[E] failed to copy file: ", err)
		return
	}

	// Wait for upload completes
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			resp.Status = "bad"
			resp.Message = "Upload failed: " + err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(marshalResponse(resp))
			log.Println("[E] failed to read file: ", err)
			return
		}
		if n == 0 {
			break
		}
	}

	resp.Status = "ok"
	resp.Message = "Upload success"
	w.Write(marshalResponse(resp))

	// Call converter to convert
	go process.Convert(uuid)
}

func marshalResponse(i APIResponseStruct) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		log.Println("[E] failed to marshal response: ", err)
		return nil
	}
	return b
}
