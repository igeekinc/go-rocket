package rocket

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

/*
This is the HTTP server that runs on the rocket.  It is used for control over WiFi, and file download
 */

// RunRocketHTTPServer creates and run the HTTP server.  It does not return unless the HTTP server has an error
func RunRocketHTTPServer(httpRoot string, httpPort int) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, path.Join(httpRoot, request.URL.Path[1:]))
	})

	requestMap := map[string]func(writer http.ResponseWriter, request *http.Request){

	}
	for request, function := range requestMap {
		http.HandleFunc("/api/"+request, function)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",httpPort), nil))
}