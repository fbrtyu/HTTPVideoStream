package webserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

		fs.ServeHTTP(w, r)
	}
}

func Webserver() {
	//go fileinfo()
	fileSystem := http.Dir("./static")

	fileServer := http.FileServer(fileSystem)

	http.Handle("/", cors(fileServer))
	http.HandleFunc("/sendvideo", getVideo)
	http.HandleFunc("/getvideo", broadcastVideo)
	http.HandleFunc("/getinfo", sendinfo)
	fmt.Println("Server start on port: 8080")
	http.ListenAndServe(":8080", nil)
}

var timeidet = 0

func fileinfo() {
	for {
		f, err := os.Open("./static/streams/video.txt")
		if err != nil {
			panic(err)
		}
		fi, err := f.Stat()
		if err != nil {
			panic(err)
		}
		fmt.Println(fi.ModTime().Second())
		timeidet = fi.ModTime().Second()
	}
}

var b []byte
var info = false

func sendinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	// if info == true {
	// 	w.Write([]byte("true"))
	// } else {
	// 	w.Write([]byte("false"))
	// }

	w.Write([]byte(strconv.Itoa(timeidet)))
}

func getVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	b = body
	//Для полной записи добавить в параметры os.O_APPEND
	file, err := os.OpenFile("./static/streams/video.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	if _, err = file.WriteString(string(body)); err != nil {
		panic(err)
	}
	defer file.Close()
}

func broadcastVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	w.Write(b)
}
