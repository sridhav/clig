package main

import "github.com/urfave/cli"

func main() {
	var commands = []cli.Command{
		cli.Command{
			Name:        "doo",
			Usage:       "do the doo",
			Description: "no really",
			Subcommands: []cli.Command{
				cli.Command{
					Name:        "doodoo",
					Usage:       "do tht doodoo",
					Description: "no no really",
					Subcommands: []cli.Command{},
				}, cli.Command{
					Name:        "doodoo2",
					Usage:       "do the doodoo2",
					Description: "no no 2 really",
					Subcommands: []cli.Command{},
				}, cli.Command{
					Name:        "doodoo3",
					Usage:       "do the doodoo3",
					Description: "no no 3 really",
					Subcommands: []cli.Command{
						cli.Command{
							Name:        "doodoodoo",
							Usage:       "do the doodoodoo",
							Description: "no no no really",
							Subcommands: []cli.Command{},
						},
					},
				},
			},
		}, cli.Command{
			Name:        "doo2",
			Usage:       "do the doo2",
			Description: "no 2 really",
			Subcommands: []cli.Command{},
		}, cli.Command{
			Name:        "doo3",
			Usage:       "do the doo3",
			Description: "no 3 really",
			Subcommands: []cli.Command{},
		},
	}

}
