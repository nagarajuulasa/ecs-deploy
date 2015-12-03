package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/credentials"
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

	// The aws go library automatically looks for credentials in various places,
	// including the ENV variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.
	//
	// We just need to override this behavior if the user explictly passes credentials

	if c.GlobalString("aws-access-key") != "" && c.GlobalString("aws-secret-key") != "" {
		provider := credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID: c.GlobalString("aws-access-key"),
				SecretAccessKey: c.GlobalString("aws-secret-key"),
			},
		}

		defaults.DefaultConfig.Credentials = credentials.NewCredentials(&provider)
	}
}

func deploy (c *cli.Context) {
	service := c.String("service-name")
	if service == "" {
		fmt.Printf("Service not specified\n\n")
		cli.ShowCommandHelp (c, "deploy")
		return
	}

	image   := c.String("image")
	if image == "" {
		fmt.Printf("Image not specified\n\n")
		cli.ShowCommandHelp (c, "deploy")
		return
	}

	checkAndConfigureAWS (c)
	//ECS := ecs.New(nil)

	fmt.Println("Deployment is work in progress...")
}

func listServices (c *cli.Context) {
	fmt.Println(c.String("timeout"))
}

func listTaskDefs  (c *cli.Context) {
	checkAndConfigureAWS (c)
	ECS := ecs.New(nil)

	// Hard Coded Value
	var maxres int64 = 10

	input := ecs.ListTaskDefinitionFamiliesInput{
		MaxResults: &maxres,
	}

	resp, err := ECS.ListTaskDefinitionFamilies(&input)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var families []*string

	families = append(families, resp.Families...)

	// The aws library only returns up to MaxResults task definitions per call, and so
	// must be iterated to ensure that all families are obtained.

	for resp.NextToken != nil {
		input = ecs.ListTaskDefinitionFamiliesInput{
			MaxResults: &maxres,
			NextToken: resp.NextToken,
		}
		
		resp, err = ECS.ListTaskDefinitionFamilies(&input)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		families = append(families, resp.Families...)
	}

	// ...and now for the pretty printing
	for i := 0; i < len(families); i++{
		fmt.Println(*families[i])

		resp, err := ECS.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
				FamilyPrefix: families[i],
				Sort: aws.String("DESC"),
			})

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var taskdefs []*string
		taskdefs = append(taskdefs, resp.TaskDefinitionArns...)

		for resp.NextToken != nil {
			
			resp, err = ECS.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
					FamilyPrefix: families[i],
					Sort: aws.String("DESC"),
					NextToken: resp.NextToken,
				})

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			taskdefs = append(taskdefs, resp.TaskDefinitionArns...)
		}

		// Whew, we now this family's list of task definitions
		// ...pretty-print time!

		for j := 0; j < len(taskdefs); j++ {
			fmt.Printf("  %v\n", *taskdefs[j])
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ecs-deploy"
	app.Usage = "Blue-Green deployments against AWS ECS"
	app.Version = "0.0.0"
	app.Action = deploy

	app.Commands = []cli.Command {
		{
			Name: "services",
			ShortName: "s",
			Usage: "List ECS Service defined for this region and cluster",
			Action: listServices,
		},
		{
			Name: "taskdefs",
			ShortName: "t",
			Usage: "List ECS Task Definitions defined for this region",
			Action: listTaskDefs,
		},
		{
			Name: "deploy",
			ShortName: "d",
			Usage: "Initiate an Blue/Green AWS deployment",
			Action: deploy,
			Flags: []cli.Flag {
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
					Usage: "Maximum number of seconds to wait until considering the task definition as failed to start",
					EnvVar: "DEPLOY_TIMEOUT",
				},
				cli.StringFlag{
					Name: "e, tag-env-var",
					Usage: "Get image tag name from environment variable. If provided this will override value specified in image name argument.",
					EnvVar: "TAG_ENV_VAR",
				},
			},
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
		cli.BoolFlag{
			Name: "V, verbose",
			Usage: "Verbose output",
		},
	}

	app.Run(os.Args)
}
