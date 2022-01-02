// env2dockerenv
//
// Convert a .env file to a starting point usable
// as a .dockerenv
//
// jum@anubis.han.de
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type kv struct {
	k, v string
}

var (
	infile  = flag.String("infile", ".env", "input .env file")
	outfile = flag.String("outfile", "-", "output .dockerenv file")
	ignore  = flag.String("ignore", "SERVE_STATIC_ASSETS,FIRESTORE_EMULATOR_HOST", "variables to ignore")
)

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}
	var (
		ifd, ofd    *os.File
		err         error
		ignore_list = strings.Split(*ignore, ",")
	)
	if *infile == "-" {
		ifd = os.Stdin
	} else {
		ifd, err = os.Open(*infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
	scanner := bufio.NewScanner(ifd)
	if *outfile == "-" {
		ofd = os.Stdout
	} else {
		ofd, err = os.Create(*outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
	var env_vars []kv
	for scanner.Scan() {
		before, after, found := strings.Cut(scanner.Text(), "=")
		if found {
			after = strings.Trim(after, "\"'")
			exp, v, found := strings.Cut(before, " ")
			if !found {
				panic(before)
			}
			if exp != "export" {
				panic(exp)
			}
			skip := false
			for _, iv := range ignore_list {
				if v == iv {
					skip = true
					break
				}
			}
			if !skip {
				env_vars = append(env_vars, kv{k: v, v: after})
			}
		} else {
			before, k, found := strings.Cut(before, " ")
			if !found {
				panic(before)
			}
			if before != "unset" {
				panic(before)
			}
			for i := range env_vars {
				if env_vars[i].k == k {
					env_vars[i].k = "X" + k
				}
			}
		}
	}
	ifd.Close()
	for _, kv := range env_vars {
		fmt.Fprintf(ofd, "%v=%v\n", kv.k, kv.v)
	}
}
