# rill
Simple command framwork in go that inject parameters in methods and more.

## Why another command framework ?

Are not you tired of having to parse the arguments received by other command frameworks?

I am and that's why I made this little library in Go (that i have already built for me in Java, Kotlin and Typescript).

## How it works ?

You need to first have a structure which is the `Command class`:

```go
type yeah struct {}
```

Then you need to implement rill.Command interface which is:

```go
type Command interface {
	CommandDescription() string
	CommandLabel() string
	CommandAliases() []string
	Clone() interface{}
}
```

So you got something like that:

```go
type yeah struct {}

func (yeah) Clone() interface{} {
	return &yeah{}
}

func (yeah) CommandDescription() string {
	return "yeah is a command !"
}

func (yeah) CommandLabel() string {
	return "yeah"
}

func (yeah) CommandAliases() []string {
	return []string{"y"}
}
```

Then to add a command you need to do:

```go
func (yeah) Lol(s string, i int64) {
	fmt.Println("command lol ", s, i)
}


func main() {
	commands := rill.Rill()

	commands.
		Add(yeah{}, "description full", "des").
    //  NAME   DESCRIPTION                ARGS NAME
    Cmd("lol", "lol command description", "s", "i").End()
```

And it's okay !


For more details check the main.go ;) 

