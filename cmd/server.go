package cmd

import (
	"github.com/qnib/qwatch/server"
    "github.com/codegangsta/cli"
)

// ServerCmd provides the flags and the execution
var ServerCmd = cli.Command{
	Name:  "server",
	Usage: "Starts daemon to run framework",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "qwatch.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
        cli.StringFlag{
			Name:  "inputs",
			Value: "gelf",
			Usage: "Comma separated list of inputs to start",
		},
		cli.StringFlag{
			Name:  "outputs",
			Value: "log",
			Usage: "Comma separated list of outputs to start",
		},
		cli.IntFlag{
			Name:  "ticker.interval",
			Value: 15000,
			Usage: "Interval of global ticker in milliseconds",
		},
		cli.IntFlag{
			Name:  "input.docker-gelf.port",
			Value: 12201,
			Usage: "UDP port of GELF collector",
		},
        cli.StringFlag{
			Name:  "input.docker-gelf.host",
			Value: "0.0.0.0",
			Usage: "UDP host of GELF collector",
		},
		cli.StringFlag{
			Name:  "output.elasticsearch.host",
			Value: "localhost",
			Usage: "Elasticsearch host to connect the ES output to",
		},
		cli.IntFlag{
			Name:  "output.elasticsearch.port",
			Value: 9200,
			Usage: "Elasticsearch port to connect the ES output to",
		},
	},
	Action: func(c *cli.Context) error {
		return qserver.ServeQlog(c)
	},
}
