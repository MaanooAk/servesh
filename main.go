package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var port int
var script string
var paramPrefix string
var allowIcos bool

func main() {
	var noParamPrefix bool

	flag.IntVar(&port, "port", 8080, "Server port number")
	flag.StringVar(&script, "script", "", "Script to handle requests")
	flag.StringVar(&paramPrefix, "prefix", "--", "S params prefix")
	flag.BoolVar(&noParamPrefix, "no-prefix", false, "Disable script params prefix")
	flag.BoolVar(&allowIcos, "icos", false, "Allow request of *.ico")
	flag.Parse()

	if noParamPrefix {
		paramPrefix = ""
	}
	if script == "" && len(flag.Args()) > 0 {
		script = flag.Args()[0]
	}

	if script == "" {
		fmt.Printf("servesh: script not provided\n")
		os.Exit(3)
	}

	if !strings.Contains(script, "{}") && !strings.HasSuffix(script, ";") {
		script += " {}"
	}

	fmt.Printf(" Script: %s\n", script)
	fmt.Printf("Address: http://localhost:%d/\n", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	if !allowIcos && strings.HasSuffix(r.URL.Path, ".ico") {
		return
	}

	commandText := GenerateCommand(script, r.URL)

	fmt.Printf("%s - %s\n", r.RemoteAddr, commandText)
	command := exec.Command("sh", "-c", commandText)

	command.Stdout = w
	command.Stderr = os.Stderr

	command.Start()
	err := command.Wait()

	if err != nil {
		fmt.Printf("%s - %s - %s\n", r.RemoteAddr, commandText, err)
	}
}

func GenerateCommand(script string, url *url.URL) string {

	params := UrlToParams(url)
	return strings.ReplaceAll(script, "{}", params)
}

func UrlToParams(url *url.URL) string {
	var sb strings.Builder

	path := url.Path[1:]
	if strings.Contains(path, " ") {
		sb.WriteString("\"")
		sb.WriteString(path)
		sb.WriteString("\"")
	} else {
		sb.WriteString(path)
	}

	for key, values := range url.Query() {
		sb.WriteString(" ")

		for _, value := range values {

			sb.WriteString(paramPrefix)
			sb.WriteString(key)

			if value != "" {
				sb.WriteString(" ")
				sb.WriteString(value)
			}
		}
	}

	return sb.String()
}

