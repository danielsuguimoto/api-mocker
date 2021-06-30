package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RouteClosure func (c *fiber.Ctx) error

type Route struct {
	Path string
	Closure RouteClosure
}

type Router struct {
	Routes []Route
}

func Create(db map[string][]map[string]interface{}) (router *Router) {
	router = &Router{}

	for resource, data := range db {
		router.addGetAll(resource, data)
		router.addGetById(resource, data)
	}

	return router
}

func (r *Router) AddRoute(path string, closure RouteClosure) {
	r.Routes = append(r.Routes, Route{
		Path: path,
		Closure: closure,
	})
}

func (r *Router) addGetAll(resource string, data []map[string]interface{}) {
	r.AddRoute(
		fmt.Sprintf("/%s", resource),
		func(c *fiber.Ctx) error {
			return c.JSON(data)
		},
	)
}

func (r *Router) addGetById(resource string, data []map[string]interface{}) {
	r.AddRoute(
		fmt.Sprintf("/%s/:id", resource),
		func(c *fiber.Ctx) error {
			id := c.Params("id")

			for _, item := range data {
				if itemId, ok := item["id"]; ok && fmt.Sprintf("%v", itemId) == id {
					return c.JSON(item)
				}
			}

			return c.SendStatus(404)
		},
	)
}

