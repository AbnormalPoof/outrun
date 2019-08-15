package main

import (
	"net/http"

	"github.com/fluofoxxo/outrun/muxhandlers/muxobj"
)

const UNKNOWN_REQUEST_DIRECTORY = "logging/unknown_requests/"

const (
	LogExecutionTime = true
)

func OutputUnknownRequest(w http.ResponseWriter, r *http.Request) {
	recv := cryption.GetReceivedMessage(r)
	// make a new logging path
	timeStr := strconv.Itoa(int(time.Now().Unix()))
	os.MkdirAll(UNKNOWN_REQUEST_DIRECTORY, 0644)
	normalizedReq := strings.ReplaceAll(r.URL.Path, "/", "-")
	path := UNKNOWN_REQUEST_DIRECTORY + normalizedReq + "_" + timeStr + ".txt"
	err := ioutil.WriteFile(path, recv, 0644)
	if err != nil {
		log.Println("[OUT] UNABLE TO WRITE UNKNOWN REQUEST: " + err.Error())
		w.Write([]byte(""))
		return
	}
	log.Println("[OUT] !!!!!!!!!!!! Unknown request, output to " + path)
	w.Write([]byte(""))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	h := muxobj.Handle
	mux := http.NewServeMux()
	// Login
	mux.HandleFunc("/Login/login/", h(muxhandlers.Login, LogExecutionTime))

	mux.HandleFunc("/", OutputUnknownRequest)
	log.Println("Starting server on port 9001")
	panic(http.ListenAndServe(":9001", mux))
}
