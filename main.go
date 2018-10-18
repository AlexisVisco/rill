package main

import (
	"errors"
	"fmt"
	"github.com/AlexisVisco/rill/rill"
	"os"
	"reflect"
)

//custom type
type all string

func fromString(s string) (interface{}, error) {
	if s == "all" {
		return all("all"), nil
	}
	return nil, errors.New("not convertible to the type all")
}

func init() {
	rill.CustomTypes[reflect.TypeOf(all(""))] = fromString
}

type yeah struct {
	Pos bool `fl:"pos,P" flDesc:"position of something"`
}

func (yeah) Clone() interface{} {
	return &yeah{}
}

func (yeah) CommandDescription() string {
	return "yeah is a command of fhwebfew"
}

func (yeah) CommandLabel() string {
	return "yeah"
}

func (yeah) CommandAliases() []string {
	return []string{"y"}
}

func (yeah) Lol(s string, i int64) {
	fmt.Println("1", s, i)
}
func (yeah) All(_ all) {
	fmt.Println("yeah all")
}

func (y yeah) Empty() {
	fmt.Println("empty", y.Pos)
}

type Y struct{}

func (Y) CommandDescription() string {
	return "full desc"
}

func (Y) CommandLabel() string {
	return "zbc"
}

func (Y) CommandAliases() []string {
	return []string{"l"}
}

func (Y) Clone() interface{} {
	return &Y{}
}

func main() {
	commands := rill.Rill()

	commands.
		Add(Y{}, "description full", "des").
		End().
		Add(yeah{}, "description full", "des").
		Cmd("Empty", "Empty short desc").
		Cmd("All", "All short desc", "*all")

	commands.Dispatch(os.Args[1:])
	//rill.HelpCommand(*cmd)
}
