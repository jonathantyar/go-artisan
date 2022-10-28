package main

import (
	"encoding/json"
	"fmt"

	"github.com/jonathantyar/go-artisan"
)

type PrintCommand struct {
	Args    float32 `artisan:"type:arg,alias:word,default:yohoho"` //required argument : will turn error if not provided
	Args2   *string `artisan:"type:arg,alias:word_opt"`            //optional argument
	Opt     string  `artisan:"type:opt,alias:test,hasValue"`       //required option : will turn error if not provided
	OptBool bool    `artisan:"type:opt,alias:inline"`              //required option : will turn error if not provided
}

type PanicCommand struct {
	Args string `artisan:"type:arg,alias:word,default:yohoho"` //required argument : will turn error if not provided
}

type MainCommand struct {
	PrintCommand PrintCommand `artisan:"alias:print,desc:For printing a text"` //name will be the command that needs to be called
	PanicCommand PanicCommand `artisan:"alias:panic"`
}

func main() {
	cmd, _ := artisan.InitCommand(MainCommand{})

	b, _ := json.Marshal(cmd)
	fmt.Println(string(b))

	fmt.Println(artisan.RunCommand)
}
