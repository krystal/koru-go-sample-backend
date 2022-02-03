package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type response struct {
	ReqCount   int      `json:"reqCount"`
	DirEntries []string `json:"dirEntries"`
}

func handleReqErr(w http.ResponseWriter, err error) {
	log.Printf("err occured in req processing: %s", err)
	w.WriteHeader(500)
}

func main() {
	lisAddr := ":9000"
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		lisAddr = ":" + portEnv
	}

	dirEnv := os.Getenv("TARGET_DIR")
	dir := "/app"
	if dirEnv != "" {
		dir = dirEnv
	}

	reqCount := 0
	httpServ := http.Server{
		Addr: lisAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount++
			res := response{
				ReqCount:   reqCount,
				DirEntries: []string{},
			}

			dirEntries, err := os.ReadDir(dir)
			if err != nil {
				handleReqErr(w, err)
				return
			}
			for _, dirEntry := range dirEntries {
				res.DirEntries = append(res.DirEntries, dirEntry.Name())
			}

			w.WriteHeader(200)
			bytes, err := json.Marshal(res)
			if err != nil {
				handleReqErr(w, err)
				return
			}
			if _, err := w.Write(bytes); err != nil {
				log.Printf("failed to write bytes: %s", err)
			}
		}),
	}

	log.Printf("Listening on %s, serving directory %s", lisAddr, dir)
	err := httpServ.ListenAndServe()
	if err != nil {
		log.Fatalf("listening failed: %s", err)
	}
}
