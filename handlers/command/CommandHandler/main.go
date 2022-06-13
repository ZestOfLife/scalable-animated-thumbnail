package main

import (
	"CommandHandler/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/reportnewjob", JobHandler)
	http.HandleFunc("/reportextract", ExtractSuccessHandler)
	http.HandleFunc("/reportresize", ResizeSuccessHandler)
	http.HandleFunc("/reportcompile", CompileSuccessHandler)
	http.HandleFunc("/reportextractfailure", ExtractFailureHandler)
	http.HandleFunc("/reportresizefailure", ResizeFailureHandler)
	http.HandleFunc("/reportcompilefailure", CompileFailureHandler)

	http.ListenAndServe(":8080", nil)
}
