package main

import "flag"
import "github.com/kardianos/service"
import "log"
import "net/http"

var mode = flag.String("mode", "run", "Application mode.  Valid values are run, install, and uninstall.")
var wwwfolder = flag.String("wwwfolder", "", "Folder that contains the website.")
var wwwport = flag.String("wwwport", "80", "The port to use when serving the website.")

var slogger service.Logger

type program struct {
	configjconn *configJSONConnector
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	log.Print("Stopping")
	return nil
}

func (p *program) run() {
	p.configjconn = NewConfigJSONConnector()
	currconfig := p.configjconn.GetConfig()

	if *wwwport != "" {
		currconfig.WWWPort = *wwwport
	}

	if *wwwfolder != "" {
		currconfig.WWWFolder = *wwwfolder
	}

	p.configjconn.SetConfig(*currconfig)

	addr := ":" + p.configjconn.GetConfig().WWWPort
	http.Handle("/", http.FileServer(http.Dir(p.configjconn.GetConfig().WWWFolder)))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Print("Error starting http.  Error: ", err)
	}
}

func main() {
	flag.Parse()
	svcConfig := &service.Config{
		Name:        "GoWebServ",
		DisplayName: "Go web server",
		Description: "Go web server",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	slogger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "run":
		err = s.Run()
	case "install":
		err = s.Install()
	case "uninstall":
		err = s.Uninstall()
	default:
		err = s.Run()
	}
	if err != nil {
		slogger.Error(err)
	}
}
