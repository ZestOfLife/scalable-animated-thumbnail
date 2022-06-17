package main

import (
	"CommandHandler/handlers"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting up Command Handler")

	http.HandleFunc("/reportnewjob", handlers.JobHandler)
	http.HandleFunc("/reportextract", handlers.ExtractSuccessHandler)
	http.HandleFunc("/reportresize", handlers.ResizeSuccessHandler)
	http.HandleFunc("/reportcompile", handlers.CompileSuccessHandler)
	http.HandleFunc("/reportextractfailure", handlers.ExtractFailureHandler)
	http.HandleFunc("/reportresizefailure", handlers.ResizeFailureHandler)
	http.HandleFunc("/reportcompilefailure", handlers.CompileFailureHandler)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
