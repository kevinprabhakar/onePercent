package main

import(
	"onePercent/util"
	"io/ioutil"
	"os"
	"net/http"
	"fmt"
)

var ServerLogger = util.NewLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

func main(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello World!")
	})

	http.ListenAndServe(":3000", nil)


}