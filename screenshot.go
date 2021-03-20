package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"

	"./authmiddleware"
)

func writeImage(w http.ResponseWriter, img string) {

	buffer, err := ioutil.ReadFile(img)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer)))
	if _, err := w.Write(buffer); err != nil {
		log.Println("unable to write image.")
	}
}

func capture(w http.ResponseWriter, r *http.Request) {
	img := "/tmp/screencapture.jpg"
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = captureLinux(img)
	} else {
		cmd = captureMac(img)
	}

	cmd.Run()
	writeImage(w, img)
}

func captureMac(file string) *exec.Cmd {
	return exec.Command("screencapture", "-x", file)
}

func captureLinux(file string) *exec.Cmd {
	return exec.Command("scrot", "-o", file)
}

func main() {
	http.HandleFunc("/screenshot", authmiddleware.AuthMiddleware(capture))
	http.ListenAndServe(":5051", nil)
}
