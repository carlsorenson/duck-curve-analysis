package duckCurve

import (
	"net/http"
)

func init() {
	http.HandleFunc("/api/", apiHandler)
	http.Handle("/", http.FileServer(http.Dir("output")))
	//http.ListenAndServe(":80", http.FileServer(http.Dir("output")))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, world again 3!")
// }
