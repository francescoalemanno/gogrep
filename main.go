package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func eatErr[T any](x T, e error) T {
	handleErr(e)
	return x
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(1)
	}
}

func main() {
	cwd := eatErr(os.Getwd())
	flag_dir := flag.String("d", cwd, "directory to work on.")
	flag_cmd := flag.String("c", "", "interpolate into command.")
	flag_casesensitive := flag.Bool("I", false, "enable case sensitivity.")
	flag.Parse()
	tofind := flag.Args()
	if !*flag_casesensitive {
		for i := range tofind {
			tofind[i] = strings.ToLower(tofind[i])
		}
	}

	dir := eatErr(os.ReadDir(*flag_dir))
	for i := range dir {
		if !dir[i].Type().IsRegular() || dir[i].IsDir() {
			continue
		}
		path := filepath.Join(*flag_dir, dir[i].Name())
		content := string(eatErr(os.ReadFile(path)))
		if !*flag_casesensitive {
			content = strings.ToLower(content)
		}
		found := true
		for j := range tofind {
			if !strings.Contains(content, tofind[j]) {
				found = false
				break
			}
		}
		if found {
			if *flag_cmd != "" {
				cmd := exec.Command("bash", "-c", strings.ReplaceAll(*flag_cmd, "{}", path))
				fmt.Print(string(eatErr(cmd.Output())))
			} else {
				fmt.Println(dir[i].Name())
			}
		}
	}
}
