package artisan

import (
	"strings"
	"time"
)

/* Command Struct */
type Command struct {
	name       string
	args       []string
	options    []string
	executedAt time.Time
}

/*
Get called command, returning string of the called command name
@return string
*/
func (c *Command) GetCommand() string {
	return c.name
}

/*
Set called command
@params string cmd
@return void
*/
func (c *Command) SetCommand(cmd string) {
	c.name = cmd
}

/* Command main struct for set name & desc of a command */
type commandMain struct {
	Name    string
	Desc    *string
	Index   int
	Args    []commandOpt
	Options []commandOpt
}

func (c *commandMain) Setter(comment string, index int) {
	f := strings.Split(comment, ",")
	for _, g := range f {
		h := strings.Split(g, ":")
		switch h[0] {
		case TagAlias:
			c.Name = h[1]
		case TagDesc:
			varToPointer := h[1]
			c.Desc = &varToPointer
		}
	}

	c.Index = index
}

func (c *commandMain) OptSetter(comment string, index int) {
	f := strings.Split(comment, ",")
	opt := commandOpt{}

	for _, g := range f {
		h := strings.Split(g, ":")
		switch h[0] {
		case TagType:
			switch h[1] {
			case TagValueArg:
				opt.Setter(comment, len(c.Args), index)
				c.Args = append(c.Args, opt)
			case TagValueOpt:
				opt.Setter(comment, len(c.Options), index)
				c.Options = append(c.Options, opt)
			}
			//catch if invalid value type
		}
		//Catch if doesnt had tags
	}
}

/* Command operator struct for set alias, default & desc of a command */
type commandOpt struct {
	Alias      []string
	Default    *string
	Desc       *string
	HasValue   bool
	IsRequired bool
	IsArray    bool
	Index      int
	IndexField int
	Value      string
}

func (c *commandOpt) Setter(comment string, index int, indexField int) {
	f := strings.Split(comment, ",")
	for _, g := range f {
		h := strings.Split(g, ":")
		switch h[0] {
		case TagAlias:
			i := strings.Split(h[1], "|")
			aliases := []string{}
			for _, alias := range i {
				aliases = append(aliases, alias)
			}
			c.Alias = aliases
		case TagDefault:
			varToPointer := h[1]
			c.Default = &varToPointer
		case TagDesc:
			varToPointer := h[1]
			c.Desc = &varToPointer
		case TagHasValue:
			c.HasValue = true
		case TagRequired:
			c.IsRequired = true
		case TagArray:
			c.IsArray = true
		}
	}

	c.Index = index
	c.IndexField = indexField
}

func (c *commandOpt) SetValue(value string) {
	c.Value = value
}

type option[T any] struct {
	Type  string
	Value T
}
