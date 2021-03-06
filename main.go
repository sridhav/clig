package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/packr"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Overwrite   string
	Name        string
	Version     string
	Author      string
	Description string
	Usage       string
	VCSHost     string
	Folder      string
	Flags       []Flag
	Commands    []Command
	CommandMap  string
	License     License
	Imports     []Import
}

type License struct {
	Header    string
	Copyright string
	Text      string
}

type Command struct {
	Name        string
	Usage       string
	Description string
	Flags       []Flag
	Commands    []Command
	Buffer      string
	Package     string
	FuncPkg     string
	Debug       string
	Header      string
	Copyright   string
}

type Import struct {
	Name string
}
type Flag struct {
	Name    string
	Type    string
	Default string
	Usage   string
}

var box packr.Box

var overwrite bool

func initBox() {
	box = packr.NewBox("./templates")
}

func main() {

	if len(os.Args) != 2 {
		usage()
	}
	initBox()
	filename := os.Args[1]
	config := new(Config)
	source, err := ioutil.ReadFile(filename)
	checkErr(err)
	err = yaml.Unmarshal(source, &config)
	checkErr(err)
	validation(config)
	addLicense(config)
	updateLicense(config)
	var buf bytes.Buffer
	checkErr(err)

	//recursive commands update
	appPath := createAppPath(config.VCSHost, config.Author, config.Name, config.Folder)
	commandPath := appPath + "/command"
	os.MkdirAll(commandPath, 0755)
	imports := make([]Import, 0)
	buf = recursiveUpdate(config.Commands, &config.Commands[0], commandPath, commandPath, &imports, config.License)
	str := strings.Replace(buf.String(), "#", "\"", -1)
	config.CommandMap = str
	for i, val := range imports {
		out := strings.Split(val.Name, config.VCSHost)
		if len(out) == 2 {
			imports[i].Name = config.VCSHost + out[1]
		}
	}
	config.Imports = imports

	// Add commands.go file
	var tempbuf bytes.Buffer
	execBufTemplate("commands.go.tmpl", &tempbuf, config)
	str = html.UnescapeString(tempbuf.String())
	file, err := os.Create(appPath + "/commands.go")
	file.WriteString(str)
	file.Close()
	checkErr(err)

	//Add main.go file
	execTemplate("main.go.tmpl", appPath+"/main.go", config)

	runGoFormat(config.VCSHost, config.Author, config.Name, config.Folder)
}

// runs go format on all generated go files
func runGoFormat(VCSHost string, user string, app string, folder string) {
	gopath := VCSHost + "/" + user + "/" + app
	if len(folder) > 1 {
		gopath = VCSHost + "/" + user + "/" + app + "/" + folder
	}
	_, _ = exec.Command("go", "fmt", gopath).Output()
	// checkErr(err)
}

// recursively updates commands buffer
func recursiveUpdate(commands []Command, callback *Command, directory string, commandPath string, imports *[]Import, license License) bytes.Buffer {
	var buf bytes.Buffer
	currDirectory := directory
	for _, element := range commands {
		element.Package = path.Base(currDirectory)
		element.FuncPkg = ""
		element.Copyright = license.Copyright
		element.Header = license.Header
		funcpkg := currDirectory
		out := strings.Split(funcpkg, commandPath)
		if len(out) == 2 {
			funcpkg = out[1]
			splits := strings.Split(funcpkg, "/")
			funcpkg = camelCase(splits)
			element.FuncPkg = funcpkg + ""
		}
		createCommandFile(currDirectory+"/"+element.Name+".go", element)
		if element.Commands != nil {
			directory = currDirectory + "/" + element.Name
			os.MkdirAll(directory, 0755)
			imp := Import{Name: directory}
			*imports = append(*imports, imp)
			recursiveUpdate(element.Commands, &element, directory, commandPath, imports, license)
		}

		execBufTemplate("command.arr.go.tmpl", &buf, element)
		callback.Buffer = buf.String()
	}
	return buf
}

// Generates Camel case string
func camelCase(splits []string) string {
	for index, element := range splits {
		splits[index] = strings.Title(element)
	}
	return strings.Join(splits, "")
}

// Generates App Path
func createAppPath(VCSHost string, user string, appname string, folder string) string {
	gopath := os.Getenv("GOPATH")
	if len(gopath) < 1 {
		gopath = userHomeDir() + "/go"
	}
	apppath := gopath + "/src/" + VCSHost + "/" + user + "/" + appname

	if len(folder) > 1 {
		apppath = apppath + "/" + folder
	}
	os.MkdirAll(apppath, 0755)
	return apppath
}

// returns user home directory
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func username() string {
	user, _ := user.Current()
	return user.Username
}

// creates individual command files
func createCommandFile(filename string, command Command) {
	path, _ := filepath.Abs(filename)
	execTemplate("commands/command.go.tmpl", path, command)
}

// writes templates to the writer
func execBufTemplate(file string, wr io.Writer, data interface{}) {
	dat, err := box.Find(file)
	tmpl, err := template.New("test").Funcs(funcMap()).Parse(string(dat))
	checkErr(err)
	err = tmpl.Execute(wr, data)
	checkErr(err)
}

// writes templates to the writer
func execTemplate(file string, outfile string, data interface{}) {
	if !fileExists(outfile) || overwrite {
		wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0644)
		checkErr(err)
		dat, err := box.Find(file)
		checkErr(err)
		tmpl, err := template.New("test").Funcs(funcMap()).Parse(string(dat))
		checkErr(err)
		err = tmpl.Execute(wr, data)
		checkErr(err)
		wr.Close()
	}
}

// checks error and panics
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// function map for templates
func funcMap() template.FuncMap {
	return template.FuncMap{
		"title":   strings.Title,
		"toUpper": strings.ToUpper,
	}
}

func addLicense(config *Config) {
	path := createAppPath(config.VCSHost, config.Author, config.Name, "")
	execTemplate("LICENSE.tmpl", path+"/LICENSE", config.License)
}

func updateLicense(config *Config) {
	commentHeader(&config.License.Header)
	commentHeader(&config.License.Copyright)
}

func commentHeader(variable *string) {
	if len(*variable) > 1 {
		*variable = "// " + *variable
	}
}

func validation(config *Config) {
	randomString := "myapp"
	requiredVariable(&config.VCSHost, "config.vcshost", "github.com")
	requiredVariable(&config.Author, "config.author", username())
	requiredVariable(&config.Name, "config.name", randomString)
	requiredVariable(&config.Folder, "config.folder", "")
	overwrite = false
	if config.Overwrite == "true" {
		overwrite = true
	}
}

func requiredVariable(variable *string, name string, def string) {
	if len(*variable) < 1 {
		fmt.Printf("WARN : Variable %s not set in yml document. Using default: %s\n", name, def)
		*variable = def
	}
}

func fileExists(file string) bool {
	ret := true
	if _, err := os.Stat(file); os.IsNotExist(err) {
		ret = false
	}
	return ret
}

func usage() {
	fmt.Println("A simplified go cli generator")
	fmt.Println("Usage: ")
	fmt.Println("clig <yml-document>")
	fmt.Println("clig clig.yml")
	os.Exit(0)
}
