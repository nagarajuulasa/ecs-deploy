package main

import (
	"os"
	"github.com/codegangsta/cli"
)

func listServices (c *cli.Context) {
}

func main() {
	app := cli.NewApp()
	app.Name = "ecs-deploy"
	app.Usage = "Blue-Green deployments against AWS ECS"
	app.Version = "0.0.0"
	app.Action = listServices

	app.Commands = []cli.Command {
		cli.Command{
			Name: "deploy",
			ShortName: "d",
			Usage: "List ECS Service defined for this region",
			Action: listServices,
		},
	}

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "k, aws-access-key",
			Usage: "AWS Access Key ID",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name: "s, aws-secret-key",
			Usage: "AWS Secret Access Key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name: "r, region",
			Usage: "AWS Region Name",
			Value: "us-east-1",
			EnvVar: "AWS_DEFAULT_REGION",
		},
		cli.StringFlag{
			Name: "c, cluster",
			Usage: "Name of ECS cluster",
			EnvVar: "AWS_ECS_CLUSTER",
		},
		cli.StringFlag{
			Name: "n, service-name",
			Usage: "Name of service to deploy",
			EnvVar: "AWS_ECS_SERVICE",
		},
		cli.StringFlag{
			Name: "i, image",
			Usage: "Name of Docker image to run, ex: mariadb:latest",
			EnvVar: "DEPLOY_IMAGE",
		},
		cli.IntFlag{
			Name: "t, timeout",
			Value: 90,
			Usage: "Maximum number of seconds to wait until the new task definition is running",
			EnvVar: "DEPLOY_TIMEOUT",
		},
		cli.StringFlag{
			Name: "g, timestamp",
			EnvVar: "CI_TIMESTAMP",
		},
	}

	app.Run(os.Args)
}
