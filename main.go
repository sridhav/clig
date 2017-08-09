package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name         string
	Version      string
	Author       string
	Description  string
	Usage        string
	VCSHost      string
	FunctionName string
	GlobalFlags  []Flag
	Commands     []Command
	CommandMap   string
}

type Command struct {
	Name        string
	Usage       string
	Description string
	Flags       []Flag
	Commands    []Command
	Buffer      string
	Package     string
}

type Flag struct {
	Name        string
	Type        string
	Description string
}

func main() {
	filename := os.Args[1]
	config := new(Config)
	source, err := ioutil.ReadFile(filename)
	checkErr(err)
	err = yaml.Unmarshal(source, &config)
	checkErr(err)
	var buf bytes.Buffer
	checkErr(err)
	// execTemplate("./templates/command.arr.go.tmpl", &buf, config)
	os.MkdirAll("./command", 0755)
	buf = recursiveUpdate(config.Commands, &config.Commands[0], "./command")
	str := strings.Replace(buf.String(), "#", "\"", -1)
	config.CommandMap = str
	//

	// var tempbuf bytes.Buffer
	// execTemplate("./templates/commands.go.tmpl", &tempbuf, config)
	// str = html.UnescapeString(tempbuf.String())
	// file, err := os.Create("/tmp/dat2")
	// file.WriteString(str)
	// file.Close()
	// // _, err = exec.Command("go", "fmt", "/tmp/dat2").Output()
	// checkErr(err)
}

// func recursiveUpdate(commands []Command, buf *bytes.Buffer) {
// 	mybuffer := new(bytes.Buffer)
// 	for _, element := range commands {
// 		if element.Commands != nil {
// 			recursiveUpdate(element.Commands, mybuffer)
// 		}
// 		element.Buffer = buf.String()
// 		execTemplate("./templates/command.arr.go.tmpl", buf, element)
// 		// fmt.Printf("element: %#v\n", element)
// 	}
// }

func recursiveUpdate(commands []Command, callback *Command, directory string) bytes.Buffer {
	var buf bytes.Buffer
	currDirectory := directory
	for _, element := range commands {
		element.Package = path.Base(currDirectory)
		createCommandFile(currDirectory+"/"+element.Name+".go", element)
		if element.Commands != nil {
			directory = directory + "/" + element.Name
			os.MkdirAll(directory, 0755)
			recursiveUpdate(element.Commands, &element, directory)
			fmt.Println(directory)
		}
		execTemplate("./templates/command.arr.go.tmpl", &buf, element)
		callback.Buffer = buf.String()
	}
	return buf
}

func createCommandFile(filename string, command Command) {
	path, _ := filepath.Abs(filename)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	checkErr(err)
	execTemplate("./templates/commands/command.go.tmpl", file, command)
	file.Close()
}

func execTemplate(file string, wr io.Writer, data interface{}) {
	dat, err := ioutil.ReadFile(file)
	tmpl, err := template.New("test").Parse(string(dat))
	checkErr(err)
	err = tmpl.Execute(wr, data)
	checkErr(err)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
