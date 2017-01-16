package cmd

import (
	"github.com/qnib/qwatch/server"
	"github.com/urfave/cli"
)

// SeverCmd provides the flags and the execution
var ServerCmd = cli.Command{
	Name:  "server",
	Usage: "Starts daemon to run framework",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "collectors",
			Value: "gelf,docker-events",
			Usage: "Comma separated list of collectors to start",
		},
		cli.StringFlag{
			Name:  "handlers",
			Value: "log",
			Usage: "Comma separated list of handlers to start",
		},
		cli.IntFlag{
			Name:  "ticker-interval",
			Value: 15000,
			Usage: "Interval of global ticker in milliseconds",
		},
		cli.IntFlag{
			Name:  "gelf-port",
			Value: 12201,
			Usage: "UDP port of GELF collector",
		},
		cli.StringFlag{
			Name:  "es-host",
			Value: "localhost",
			Usage: "Elasticsearch host to connect the ES output to",
		},
		cli.IntFlag{
			Name:  "es-port",
			Value: 9200,
			Usage: "Elasticsearch port to connect the ES output to",
		},
	},
	Action: func(c *cli.Context) error {
		return qserver.ServeQlog(c)
	},
}
