package main

import (
	"github.com/AlexisVisco/rill/rill"
	"fmt"
)

type yeah struct {}

func (yeah) CommandDescription() string {
	return "yeah is a command"
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
func (yeah) LolX(i int64, s string) {
	fmt.Println("2", s, i)
}
func (yeah) Empty() {
	fmt.Println("empty")
}

func main() {
	commands := rill.Rill()

	commands.Add(yeah{}, "Lol", "", "")
	commands.Add(yeah{}, "LolX", "", "")
	commands.Add(yeah{}, "Empty", "", "")

	commands.Dispatch("yeah 123 yeah")
	commands.Dispatch("y yeah 123")
	commands.Dispatch("y")
}
