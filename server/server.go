package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/danielsuguimoto/api-mocker/router"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
	port int
	db map[string][]map[string]interface{}
}

func Create(port int) *Server {
	return &Server{
		app: fiber.New(),
		port: port,
	}
}

func (s *Server) LoadResources(path string) (err error) {
	var (
		file *os.File
		raw []byte
	)

	if file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	defer file.Close()

	if raw, err = io.ReadAll(file); err != nil {
		return
	}

	s.db = make(map[string][]map[string]interface{})
	err = json.Unmarshal(raw, &s.db)

	s.bindRouter(router.Create(s.db))

	return
}

func (s *Server) Listen() {
	go func() {
		_ = s.app.Listen(fmt.Sprintf(":%v", s.port))
	}()
}

func (s *Server) WatchResource(path string) {
	watcher, _ := fsnotify.NewWatcher()

	defer watcher.Close();

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						s.app.Shutdown();
						s.app = fiber.New()
						s.LoadResources(path)
						s.Listen()
					}
				}
			case err := <-watcher.Errors:
				{
					fmt.Println("ERROR", err)
				}
			}
		}
	}()

	_ = watcher.Add(path);

	<-done
}

func (s *Server) bindRouter(router *router.Router) {
	for _, route := range router.Routes {
		s.app.Get(route.Path, route.Closure)
	}
}
