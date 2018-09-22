package docopt_unmarshal

import (
	"github.com/stretchr/testify/assert"
	"github.com/docopt/docopt-go"
	"reflect"
	"strconv"
	"strings"
	. "testing"
	"time"
	"net/url"
)

func TestBasics(t *T) {
	type NavalFate struct {
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

	naval_fate := `Naval Fate.

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

	parse_unmarshal_assert(t, "basic",   naval_fate, []string{"ship", "A", "move", "10", "20", "--speed=30"}, &NavalFate{}, &NavalFate{ Ship:true, Name:[]string{"A"}, Move:true, X:10, Y:20, Speed:30})
	parse_unmarshal_assert(t, "default", naval_fate, []string{"ship", "A", "move", "10", "20"              }, &NavalFate{}, &NavalFate{ Ship:true, Name:[]string{"A"}, Move:true, X:10, Y:20, Speed:10})
}

func TestComplexTypes(t *T) {
	type Opts struct {
		Duration time.Duration `docopt:"<duration>"`
		Float float32 `docopt:"<float>"`
	}
	parse_unmarshal_assert(t, "duration+float", `Usage: test <duration> <float>`, []string{"20s", "3.14"}, &Opts{}, &Opts{20*time.Second, 3.14})
}

func TestNestedStructs(t *T) {
	type C 	  struct {     X bool    `docopt:"x"` }
	type B 	  struct {C C; Y int32   `docopt:"<y>"`}
	type A 	  struct {B B; Z string  `docopt:"<z>"`}
	type Opts struct {A A; P float32 `docopt:"<p>"`}
	parse_unmarshal_assert(t, "nested structs", `Usage: test x <y> <z> <p>`, []string{"x", "27", "Hello", "2.718"}, &Opts{}, &Opts{A:A{B:B{C:C{X:true},Y:27},Z:"Hello"},P:2.718})
}

func TestComplicatedAPI(t *T) {
	type Opts struct {
		Duration time.Duration `docopt:"<duration>"`
		Float float32 `docopt:"<float>"`
	}
	parse_unmarshal_hook_assert(t, "duration+float", `Usage: test <duration> <float>`, []string{"20s", "3.14"}, &Opts{}, &Opts{20*time.Second, 3.14}, nil)
}

func TestTypeHook(t *T) {
	type Negative int
	type Opts struct {
		Negative Negative `docopt:"<int>"`
	}
	parse_unmarshal_hook_assert(t, "type hook", `Usage: test <int>`, []string{"10"}, &Opts{}, &Opts{-10}, map[string]Hook{
		"docopt_unmarshal.Negative": func(f_val reflect.Value, arg string) error {
			iv, err := strconv.ParseInt(arg, 10, 64)
			if err == nil {
				f_val.SetInt(-iv)
			}
			return nil
		},
	})
}

func TestURL(t *T) {
	type Opts struct {
		Url *url.URL `docopt:"<url>"`
	}
	u := "http://example.com/a/path"
	up, _ := url.Parse(u)
	parse_unmarshal_hook_assert(t, "type hook", `Usage: test <url>`, []string{u}, &Opts{}, &Opts{up}, nil)
}


func parse_unmarshal_assert(t *T, name string, doc string, argv []string, structure interface{}, assertion interface{}) {
	t.Run(name, func(t *T) {
		arguments, err := docopt.ParseArgs(doc, argv, "Naval Fate 2.0")
		assert.Nil(t, err, `Docopt Parse`)
		err = DocoptUnmarshal(arguments, structure)
		assert.Nil(t, err, `DocoptUnmarshal`)
		assert.Equal(t, structure, assertion, strings.Join(argv, " "))
	})
}

func parse_unmarshal_hook_assert(t *T, name string, doc string, argv []string, structure interface{}, assertion interface{}, hooks map[string]Hook) {
	t.Run(name, func(t *T) {
		arguments, err := docopt.ParseArgs(doc, argv, "Naval Fate 2.0")
		assert.Nil(t, err, `Docopt Parse`)
		um := New()
		um.AddHooks(hooks)
		err = um.Unmarshal(arguments, structure)
		assert.Nil(t, err, `DocoptUnmarshal`)
		assert.Equal(t, structure, assertion, strings.Join(argv, " "))
	})
}
