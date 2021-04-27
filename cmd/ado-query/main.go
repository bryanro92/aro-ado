package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func main() {
	personalAccessToken := os.Getenv("ADO_TOKEN")
	if len(personalAccessToken) == 0 {
		log.Fatal("Token Required")
	}
	project := os.Getenv("ADO_PROJECT")
	if len(project) == 0){
		log.Fatal("Project required")
	}
	organizationUrl := "https://dev.azure.com/msazure/"

	// Create a connection to your organization
	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	ctx := context.Background()

	// Create a client to interact with the Core area
	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		log.Print("coreClient err")
		log.Fatal(err)
	}

	// Get first page of the list of team projects for your organization
	responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})

	if err != nil {
		log.Fatal(err)
	}

	index := 0
	for responseValue != nil {
		// Log the page of team project names
		for _, teamProjectReference := range (*responseValue).Value {
			// log.Printf("Name[%v] = %v", index, *teamProjectReference.Name)
			index++
			if strings.EqualFold(*teamProjectReference.Name, project) {
				log.Printf("Name[%v] = %v", index, *teamProjectReference.Name)
				log.Println("Yay, ARO")
			}
		}

		// if continuationToken has a value, then there is at least one more page of projects to get
		if responseValue.ContinuationToken != "" {
			// Get next page of team projects
			projectArgs := core.GetProjectsArgs{
				ContinuationToken: &responseValue.ContinuationToken,
			}
			responseValue, err = coreClient.GetProjects(ctx, projectArgs)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			responseValue = nil
		}
	}
	buildClient, err := build.NewClient(ctx, connection)
	// buildClient.GetBuildsArgs{
	// 	Project: &project,
	// }
	buildResponse, err := buildClient.GetBuilds(ctx, build.GetBuildsArgs{Project: &project, MinTime: &azuredevops.Time{time.Now().Add(-1 * time.Second)}})
	println(buildResponse)
	fmt.Printf("BuildResponse Size: %T, %d\n", buildResponse, unsafe.Sizeof(buildResponse))
	i := 0
	for buildResponse != nil {
		for _, buildValue := range (*buildResponse).Value {
			log.Printf("Name[%v] = %v", i, *buildValue.Result)
			i++
		}

		// if buildResponse.ContinuationToken != "" {
		// 	buildArgs := build.GetBuildArgs{
		// 		Project: &project,
		// 		MinTime: &azuredevops.Time{
		// 			time.Now().Add(-1 * time.Hour),
		// 		},
		// 		ContinuationToken: &buildResponse.ContinuationToken,
		// 	}
		// }
	}
}
