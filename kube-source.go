package main

import (
	"log"
	"github.com/docopt/docopt-go"
)

const version = "kube-source 0.1.0"
const usage = `
Usage:
	kube-source server
	kube-source --config <config> server
	kube-source --help
	kube-source --version

Options:
	--config <config>            The kube-source config [default: /etc/kube-source/kube-source.json].
	--help                       Show this screen.
	--version                    Show version.
`

func main() {
	// Parse args
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetPrefix("[kube-source] ")
	log.Println("starting kube-source")

	if args["server"].(bool) {
		configFile := args["--config"].(string)
		log.Printf("using config: %v\n", configFile)

		err := config.LoadConfig(configFile)
		if err != nil {
			log.Fatalln(err)
		}

		// echo := app.BuildApp(config)
		server := app.BuildApp()
		log.Printf("listening on %v\n", config.ListenString)
		server.Run(fasthttp.New(config.ListenString))
	}
}