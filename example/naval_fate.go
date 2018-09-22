package main

import (
	"fmt"
	"github.com/caldwell/docopt-unmarshal"
	"github.com/docopt/docopt-go"
	"log"
)

type Options struct {
	Help     bool     `docopt:"--help"`
	Ship     bool     `docopt:"ship"`
	New      bool     `docopt:"new"`
	Name     []string `docopt:"<name>"`
	Move     bool     `docopt:"move"`
	X        int64    `docopt:"<x>"`
	Y        int64    `docopt:"<y>"`
	Speed    int32    `docopt:"--speed"`
	Shoot    bool     `docopt:"shoot"`
	Mine     bool     `docopt:"mine"`
	Set      bool     `docopt:"set"`
	Remove   bool     `docopt:"remove"`
	Moored   bool     `docopt:"--moored"`
	Drifting bool     `docopt:"--drifting"`
	Version  bool     `docopt:"--version"`
}

func main() {
	  usage := `Naval Fate.

Usage:
  naval_fate ship new <name>...
  naval_fate ship <name> move <x> <y> [--speed=<kn>]
  naval_fate ship shoot <x> <y>
  naval_fate mine (set|remove) <x> <y> [--moored|--drifting]
  naval_fate -h | --help
  naval_fate --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  --speed=<kn>  Speed in knots [default: 10].
  --moored      Moored (anchored) mine.
  --drifting    Drifting mine.`

	arguments, err := docopt.Parse(usage, nil, true, "Naval Fate 2.0", false)
        if err != nil { log.Fatal("docopt: ", err) }
	var options Options
	err = docopt_unmarshal.DocoptUnmarshal(arguments, &options)
        if err != nil { log.Fatal("options: ", err) }
	fmt.Printf("%#v\n", options)
}
