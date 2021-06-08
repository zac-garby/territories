package server

import (
	"bytes"
	"encoding/json"
	"errors"
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
	g := &game.Game{}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := s.handleMessage(msg, g, conn); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) handleMessage(msg []byte, g *game.Game, conn *websocket.Conn) error {
	var cmd, param []byte

	if len(msg) > 4 {
		cmd = msg[:3]
		param = msg[4:]
		if msg[3] != ' ' {
			if err := conn.WriteMessage(websocket.TextMessage, RESP_INVALID); err != nil {
				return err
			}

			return errors.New("invalid format - missing fourth character space")
		}
	} else {
		cmd = msg[:3]
	}

	fmt.Println(string(cmd), string(param))

	if bytes.Equal(cmd, CMD_GENERATE) {
		// generate a new world
		*g = *game.NewGame(600, 600, 40, 10)
		log.Println("a new game has been created")
		return conn.WriteMessage(websocket.TextMessage, RESP_GEN)
	} else if bytes.Equal(cmd, CMD_POLYGON) {
		// get the polygons for the current world
		if g == nil {
			if err := conn.WriteMessage(websocket.TextMessage, RESP_NOGAME); err != nil {
				return err
			}

			return errors.New("no game has been created")
		}

		regions := g.World.Regions
		points := make([][]float64, len(regions))
		for i, reg := range regions {
			points[i] = make([]float64, len(reg)*2)
			for j, p := range reg {
				points[i][2*j] = p.X
				points[i][2*j+1] = p.Y
			}
		}

		regionsJson, _ := json.Marshal(points)
		return conn.WriteMessage(websocket.TextMessage, append(RESP_POLYGON, regionsJson...))
	}

	return nil
}
