package server

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

// Запускает файловый сервер на порту port, считая корневой директорией htmlDir
// Открывает путь для веб-сокета
func Run(rootDir string, port string) {
	router := mux.NewRouter()

	router.Methods("GET").Path("/websocket").HandlerFunc(createWSHandler())
	router.Methods("GET").Path("/websocket-addition").HandlerFunc(createAdditionWSHandler())
	router.Methods("GET").Path("/").HandlerFunc(createPageHandler(rootDir))
	router.Methods("GET").Path("/{file}").HandlerFunc(createPageHandler(rootDir))
	router.Methods("GET").Path("/{dir}/{path}").HandlerFunc(createPathHandler(rootDir))

	http.Handle("/", router)

	err := http.ListenAndServe("127.0.0.1:" + port, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
