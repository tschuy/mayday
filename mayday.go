package main

import (
	"archive/tar"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/mayday/mayday"
)

var (
	flagConfigFile *string
	danger     *bool
)

const (
	dirPrefix     = "/mayday"
)

type File struct {
	Name string
	Link string
}

type Command struct {
	Args []string
	Link string
}

type Config struct {
	Files    []File
	Commands []Command
}

func openConfig() (string, error) {
	configVar := os.Getenv("MAYDAY_CONFIG_FILE")
	configFile := strings.Split(configVar, "=")[0]

	if configFile == "" {
		configFile = *flagConfigFile
	}

	log.Printf("Reading configuration from %v\n", configFile)
	cbytes, err := ioutil.ReadFile(configFile)
	cstring := string(cbytes)
	return cstring, err
}

func readConfig(dat string) ([]File, []Command, error) {
	var c Config

	err := json.Unmarshal([]byte(dat), &c)
	if err != nil {
		log.Fatal(err)
	}
	return c.Files, c.Commands, nil
}

func main() {
	flagConfigFile = flag.String("config-file", "/etc/mayday.conf", "config file location")
	danger = flag.Bool("danger", false, "collect potentially private information (ex, container logs)")

	flag.Parse()

	conf, err := openConfig()
	if err != nil {
		log.Fatal(err)
	}

	files, commands, err := readConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	journals, err := mayday.ListJournals()
	if err != nil {
		log.Fatal(err)
	}

	var tarables []mayday.Tarable

	for _, f := range files {
		content, err := os.Open(f.Name)
		if err != nil {
			log.Fatal(err)
		}
		defer content.Close()

		fi, err := os.Stat(f.Name)
		if err != nil {
			log.Fatal(err)
		}

		header, err := tar.FileInfoHeader(fi, f.Name)
		header.Name = f.Name
		if err != nil {
			log.Fatal(err)
		}

		tarables = append(tarables, mayday.NewFile(content, header, f.Name, f.Link))
	}

	for _, c := range commands {
		tarables = append(tarables, mayday.NewCommand(c.Args, c.Link))
	}

	for _, j := range journals {
		tarables = append(tarables, j)
	}

	now := time.Now().Format("200601021504.999999999")
	ws := os.TempDir() + dirPrefix + now

	var t mayday.Tar
	outputFile := ws + ".tar.gz"
	tarfile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()
	t.Init(tarfile)

	mayday.Run(t, tarables)
	t.Close()

	log.Printf("Output saved in %v\n", outputFile)
	log.Printf("All done!")

	return
}
