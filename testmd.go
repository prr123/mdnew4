// test program that uses the cli to parse different md files.

package main

import (
	"os"
	"fmt"
	"log"
	mdp "goDemo/mdnew4/mdparser"
	util "github.com/prr123/utility/utilLib"
)


func main() {

    numArgs := len(os.Args)
    var md []byte

    flags:=[]string{"dbg","md"}

    useStr := "/md=<markdown file> [/dbg]"
    helpStr := fmt.Sprintf("help: This program parses an md file into blocks\n")

    if numArgs > len(flags)+1 {
        fmt.Println("too many arguments in cl!")
        fmt.Printf("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    if numArgs == 2 {
        if os.Args[1] == "help" {
            fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
            fmt.Printf("%s\n", helpStr)
            os.Exit(1)
        }
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg := false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    mdFilnam := ""
    mdval, ok := flagMap["md"]
    if ok {
        if mdval.(string) == "none" {log.Fatalf("error -- no markdown file provided with /md flag!")}
        mdFilnam = mdval.(string)
        mdFullFilnam := "mdFiles/" + mdFilnam + ".md"
        md, err = os.ReadFile(mdFullFilnam)
        if err != nil {log.Fatalf("error -- cannot read md: %v", err)}
    } else {
        log.Fatalf("error -- no md flag provided!")
    }

	if dbg {
		fmt.Printf("md file: %s\n", mdFilnam)
	}

    p := mdp.InitParser(md)
    ps := mdp.InitParseState(md)

    p.PrintLines()

    err = p.Parse(ps)
    if err != nil {log.Fatal("cannot parse nd!")}

	mdp.PrintBlock(p.BlockList)
}
