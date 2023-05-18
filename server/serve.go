package server

import (
	"io"
	"log"
	"mso-pdf-renderer/manager"
	"mso-pdf-renderer/process"
	"net/http"
	"os"
)

var downloadSemaphore = make(chan struct{}, 1)

func ListenAndServe(addr string) {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/check", checkHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Reads web/index.html
		f, err := os.Open(process.RunningPath + "/web/index.html")
		if err != nil {
			log.Println("[E] failed to open index.html: ", err)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})
	http.ListenAndServe(addr, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	extension := r.FormValue("extension")
	if extension == "" {
		w.Write([]byte("Invalid extension"))
		return
	}
	uuid := manager.GenerateUUID()
	manager.Routines = append(manager.Routines, manager.RoutineStruct{
		UUID:          uuid,
		FileExtension: "." + extension,
	})
	w.Write([]byte(uuid))

}

func checkHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		w.Write([]byte("Invalid uuid"))
		return
	}

	if _, err := os.Stat(process.RunningPath + "/cache/" + uuid + manager.FindRoutine(uuid).FileExtension); os.IsNotExist(err) {
		w.Write([]byte("Not done yet"))
		return
	}

	w.Write([]byte("Done"))

}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		w.Write([]byte("Invalid uuid"))
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
	r.ParseForm()
	uuid := r.FormValue("uuid")
	if !manager.DoesUUIDExist(uuid) {
		w.Write([]byte("Invalid uuid"))
		return
	}

	extension := r.FormValue("extension")
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		log.Println("[E] failed to upload file: ", err)
		return
	}
	defer file.Close()

	fileLocal, err := os.Create(process.RunningPath + "/cache/" + uuid + extension)
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		log.Println("[E] failed to create file: ", err)
		return
	}
	defer fileLocal.Close()

	_, err = io.Copy(fileLocal, file)
	if err != nil {
		w.Write([]byte("Upload failed: " + err.Error()))
		log.Println("[E] failed to copy file: ", err)
		return
	}

	// Wait for upload completes
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			w.Write([]byte("Upload failed: " + err.Error()))
			log.Println("[E] failed to read file: ", err)
			return
		}
		if n == 0 {
			break
		}
	}

	w.Write([]byte("Upload success"))

	// Call converter to convert
	go process.ConvertPPT(process.RunningPath+"/cache/"+uuid+extension, process.RunningPath+"/cache/"+uuid+".pdf")
}
