package main

import (
	"go/build"
	"log"
	"os"
	"os/exec"
	"strings"
)

type (
	aCommand struct {
		name string
		args []string
		rm   bool
	}
)

var (
	ctx   = build.Default
	goCmd = ctx.GOROOT + "/bin/go"
	pwd   = ctx.GOPATH + "/src/github.com/corestoreio/csfw/"
)

func main() {
	cmds := []aCommand{
		aCommand{
			goCmd,
			[]string{"build", "-a", "github.com/corestoreio/csfw/tools/tableToStruct"},
			false,
		},
		aCommand{
			"rm",
			[]string{"eav/generated_tables.go", "eav/generated_eav.go"},
			false,
		},
		aCommand{
			pwd + "tableToStruct",
			[]string{"-p", "eav", "-prefixSearch", "eav", "-o", "eav/generated_tables.go", "-run"},
			true,
		},
		aCommand{
			goCmd,
			[]string{"build", "-a", "github.com/corestoreio/csfw/tools/eavToStruct"},
			false,
		},
		aCommand{
			pwd + "eavToStruct",
			[]string{"-p", "eav", "-o", "eav/generated_eav.go", "-run"},
			true,
		},
	}

	for _, cmd := range cmds {
		out, err := exec.Command(cmd.name, cmd.args...).Output()
		if cmd.rm {
			defer os.Remove(cmd.name)
		}
		if err != nil {
			log.Fatalf("Failed: %s => %s", cmd.name, err)
		}

		// @todo cmd error output
		log.Printf("%s %s:\n%s\n", cmd.name, strings.Join(cmd.args, " "), out)
	}
}
