package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func checkAndConfigureAWS (c *cli.Context) {
	// Checking the "--region" flag,
	// then the AWS_DEFAULT_REGION env variable,
	// then the AWS_REGION env variable.
	// Defaults to "us-east-1" if nothing else is set

	region := c.GlobalString("region")

	if region != "" {
		defaults.DefaultConfig.Region = aws.String(c.GlobalString("region"))
	} else if os.Getenv("AWS_DEFAULT_REGION") != "" {
		defaults.DefaultConfig.Region = aws.String(os.Getenv("AWS_DEFAULT_REGION"))
	} else if os.Getenv("AWS_REGION") != "" {
		defaults.DefaultConfig.Region = aws.String(os.Getenv("AWS_REGION"))
	} else {
		defaults.DefaultConfig.Region = aws.String("us-east-1")
	}
}

func listServices (c *cli.Context) {
	fmt.Println(c.String("timeout"))
}

func listTaskDefs  (c *cli.Context) {
	checkAndConfigureAWS (c)

	svc := ecs.New(nil)

params := &ecs.ListTaskDefinitionsInput{
    Sort:         aws.String("DESC"),
    Status:       aws.String("ACTIVE"),
}
resp, err := svc.ListTaskDefinitions(params)

if err != nil {
    // Print the error, cast err to awserr.Error to get the Code and
    // Message from an error.
    fmt.Println(err.Error())
    return
}

// Pretty-print the response data.
fmt.Println(resp)


}

func main() {
	app := cli.NewApp()
	app.Name = "ecs-deploy"
	app.Usage = "Blue-Green deployments against AWS ECS"
	app.Version = "0.0.0"
	app.Action = listTaskDefs

	app.Commands = []cli.Command {
		cli.Command{
			Name: "services",
			ShortName: "s",
			Usage: "List ECS Service defined for this region and cluster",
			Action: listServices,
		},
		cli.Command{
			Name: "taskdefs",
			ShortName: "t",
			Usage: "List ECS Task Definitions defined for this region",
			Action: listTaskDefs,
		},
	}

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "k, aws-access-key",
			Usage: "AWS Access Key ID [$AWS_ACCESS_KEY_ID]",

			// Not used, the aws SDK automatically retrives this value
			//EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name: "s, aws-secret-key",
			Usage: "AWS Secret Access Key [$AWS_SECRET_ACCESS_KEY]",

			// Not used, the SDK automatically retrives this value
			//EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name: "r, region",
			Usage: "AWS Region Name. Defaults to \"us-east-1\"",

			// Not used. Extra logic needed for backward compatibility
			//  See above: checkAndConfigureAWS
			//EnvVar: "AWS_DEFAULT_REGION",
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
			Name: "e, tag-env-var",
			Usage: "Get image tag name from environment variable. If provided this will override value specified in image name argument.",
			EnvVar: "TAG_ENV_VAR",
		},
		cli.BoolFlag{
			Name: "V, verbose",
			Usage: "Verbose output",
		},
	}

	app.Run(os.Args)
}
