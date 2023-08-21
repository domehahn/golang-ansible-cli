/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"golang-ansible/pkg/environment"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ansibleProject = "ansible-playground"

const (
	CliAnsibleDirectory   = "cli.ansible.directory"
	CliAnsibleVars        = "cli.ansible.vars"
	CliAnsiblePlaybook    = "cli.ansible.playbook"
	CliAnsibleRolesCommon = "cli.ansible.roles.common"
)

type TemplateData struct {
	Description string
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Ansible Structure",
	Long:  `Create Ansible basic structure within the given path.`,
	Run: func(cmd *cobra.Command, args []string) {
		ansiblePath, err := environment.ViperGetEnvVariable(CliAnsibleDirectory)

		if err != nil {
			log.Fatal("Error")
		}

		if len(args) >= 1 && args[0] != "" {
			ansibleProject = args[0]
		}

		if err := os.Mkdir(ansiblePath+ansibleProject, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		ansibleProjectPath := ansiblePath + ansibleProject

		addVarFiles(ansibleProjectPath)

		addPlaybook(ansibleProjectPath)

		if err := os.Mkdir(ansibleProjectPath+"/roles", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		if err := os.Mkdir(ansibleProjectPath+"/roles/common", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		commons, err := environment.ViperGetEnvVariableSlice(CliAnsibleRolesCommon)

		if err != nil {
			log.Fatal("Error")
		}

		for _, common := range commons {
			if err := os.Mkdir(ansibleProjectPath+"/roles/common/"+common, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("create called")
	},
}

func addPlaybook(ansibleProjectPath string) {
	playbooks, err := environment.ViperGetEnvVariableSlice(CliAnsiblePlaybook)

	if err != nil {
		log.Fatal("Error")
	}

	for _, playbook := range playbooks {
		filePath := filepath.Join(ansibleProjectPath, playbook+".yml")

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("Error")
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal("Error")
			}
		}(file)
	}
}

func addVarFiles(ansibleProjectPath string) {
	vars, err := environment.ViperGetEnvVariableSlice(CliAnsibleVars)

	if err != nil {
		log.Fatal("Error")
	}

	for _, variable := range vars {
		if err := os.Mkdir(ansibleProjectPath+"/"+variable, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		filePath := filepath.Join(ansibleProjectPath+"/"+variable, variable+".yml")
		tmpl, err := template.New(variable).Parse("#add vars here")
		if err != nil {
			log.Fatal("Error")
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("Error")
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal("Error")
			}
		}(file)

		err = tmpl.Execute(file, nil)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
