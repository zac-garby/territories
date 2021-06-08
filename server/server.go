package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/zac-garby/territories/game"
)

type Server struct {
	s *http.Server
}

func NewServer(addr string) *Server {
	serv := &Server{}

	r := mux.NewRouter()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 4096,
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	r.PathPrefix("/ws/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			// TODO: handle this properly
			return
		}

		defer conn.Close()

		serv.connection(conn)
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	serv.s = &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return serv
}

func (s *Server) Start() {
	fmt.Printf("starting server on %s\n", s.s.Addr)
	s.s.ListenAndServe()
}

func (s *Server) connection(conn *websocket.Conn) {
	game := game.NewGame()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println(err)
			return
		}
	}
}
