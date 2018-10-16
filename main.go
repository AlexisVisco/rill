package main

import (
	"fmt"
	"github.com/AlexisVisco/rill/rill"
)

type yeah struct {
	Pos bool `fl:"p,po" flDesc:"position of something"`
}

func (yeah) Clone() interface{} {
	return &yeah{}
}

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

func (y yeah) Empty() {
	fmt.Println("empty", y.Pos)
}

func main() {
	commands := rill.Rill()

	commands.
		Add(yeah{}, "description full", "des").
		Cmd("Empty", "Empty full desc", "Empty short desc")

	commands.Dispatch([]string{"yeah", "--pos"})
}
