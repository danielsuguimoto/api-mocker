package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/danielsuguimoto/api-mocker/router"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
	db map[string][]map[string]interface{}
}

func Create() *Server {
	return &Server{
		app: fiber.New(),
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

func (s *Server) Listen(port int) error {
	return s.app.Listen(fmt.Sprintf(":%v", port))
}

func (s *Server) bindRouter(router *router.Router) {
	for _, route := range router.Routes {
		s.app.Get(route.Path, route.Closure)
	}
}
