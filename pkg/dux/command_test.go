package dux

import (
	"fmt"
	"testing"
)

func TestRoute_AddRoute(t *testing.T) {
	rt := Command{Name: "route"}
	testRoute := &Command{Name: "addedRoute"}
	rt.AddCommand(testRoute)
	if rt.commands["addedRoute"] != testRoute {
		t.Fail()
	}
}

func TestRoute_GetRoute(t *testing.T) {
	rt := Command{Name: "route"}
	wanted := &Command{Name: "addedRoute"}
	rt.AddCommand(wanted)

	t.Run("GetRoute returns desired route", func(t *testing.T) {
		got, ok := rt.GetRoute("addedRoute")
		if !ok {
			fmt.Printf("could not find route addedRoute\n")
			t.Fail()
		}
		if got != wanted {
			fmt.Printf("wrong route found.\nGot: %v\nWanted: %v\n", got, wanted)
			t.Fail()
		}
	})
	t.Run("GetRoute returns false when route was not found", func(t *testing.T) {
		got, ok := rt.GetRoute("notFound")
		if ok {
			fmt.Printf("route was found: %v\n", got)
			t.Fail()
		}
	})
}

func TestRoute_DeepFind(t *testing.T) {
	testRoute := Command{
		Name: "bot",
		commands: map[string]*Command{
			"sub1": {Name: "sub1", commands: map[string]*Command{}},
			"sub2": {Name: "sub2", commands: map[string]*Command{
				"sub2sub1": {Name: "sub2sub1", commands: map[string]*Command{}},
				"sub2sub2": {Name: "sub2sub2", commands: map[string]*Command{
					"sub2sub2sub1": {Name: "sub2sub2sub1", commands: map[string]*Command{}},
				}},
			}},
		},
	}

	argsToCheck := [][]string{
		{"sub1", "thoseAreArgs"},
		{"sub2", "sub2sub1", "thoseAreArgs"},
		{"sub2", "sub2sub2", "sub2sub2sub1", "thoseAreArgs"},
	}
	for _, testArgs := range argsToCheck {
		rt, args := testRoute.DeepFind(testArgs)
		if rt.Name != testArgs[len(testArgs)-2] || args[0] != "thoseAreArgs" {
			fmt.Printf("wrong results\nroute - %v\nargs - %v", rt, args)
			t.Fail()
		}
	}
}
