package webserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

		fs.ServeHTTP(w, r)
	}
}

func Webserver() {
	fileSystem := http.Dir("./static")
	fileServer := http.FileServer(fileSystem)

	http.Handle("/", cors(fileServer))
	http.HandleFunc("/sendvideo", getVideo)
	http.HandleFunc("/getvideo", broadcastVideo)
	fmt.Println("Server start on port: 8080")
	http.ListenAndServe(":8080", nil)
}

var b []byte
var firstchunck []byte
var flag = false //Чтобы один раз записать первый чанк
var flag2 = true //Чтобы один раз отправить первый чанк, а дальше уже самый новый отправлять

//Но в этих флагах есть проблема, что если перезапустить стрим без перезапуска сервера, то чанки не будут друг с другом работать так как они будут из разных видео

func getVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()

	b = body

	if flag == false {
		firstchunck = body
		flag = true
	}

	//Для полной записи добавить в параметры os.O_APPEND
	file, err := os.OpenFile("./static/streams/video.txt", os.O_WRONLY|os.O_CREATE, 0600)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Content-Length")

	if flag2 == true {
		w.Write(firstchunck)
		flag2 = false
	} else {
		w.Write(b)
	}
}
