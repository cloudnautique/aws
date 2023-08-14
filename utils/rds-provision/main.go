package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/acorn-io/aws/utils/rds-provision/pkg/aws/cloudformation"
	"github.com/sirupsen/logrus"
)

func main() {
	stackName := os.Getenv("ACORN_EXTERNAL_ID")
	event := os.Getenv("ACORN_EVENT")

	if event == "create" || event == "update" {
		//if err := GenerateTemplateFile("cfn.yaml"); err != nil {
		//logrus.Fatal(err)
		//}

		if err := applyCfnTemplateFile("cfn.yaml", stackName); err != nil {
			logrus.Fatal(err)
		}
	}

	if event == "delete" {
		if err := deleteStack(stackName); err != nil {
			logrus.Fatal(err)
		}
	}
}

func applyCfnTemplateFile(inputFile, stackName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client, err := cloudformation.NewClient(ctx)
	if err != nil {
		return err
	}

	templateBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	if err := cloudformation.DeployStack(client, stackName, string(templateBytes)); err != nil {
		return err
	}

	if err := cloudformation.DumpOutputsToFile(client, stackName, "outputs.json"); err != nil {
		return err
	}
	return runServiceAcornRenderExec("./service.sh")
}

func deleteStack(stackName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client, err := cloudformation.NewClient(ctx)
	if err != nil {
		return err
	}
	return cloudformation.Delete(client, stackName)
}

func GenerateTemplateFile(outputFile string) error {
	cmd := exec.Command("cdk", "synth", "--path-metadata", "false", "--lookups", "false")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running cdk synth: %v, %v", err, stderr.String())
	}

	return writeCFNTemplate(out.String(), outputFile)
}

func writeCFNTemplate(content, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func runServiceAcornRenderExec(executable string) error {
	cmd := exec.Command(executable)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &out
	return cmd.Run()

}
