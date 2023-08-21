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
var createPathParameter string

const (
	CliAnsibleDirectory   = "cli.ansible.directory"
	CliAnsibleGroupVars   = "cli.ansible.group_vars"
	CliAnsibleHostVars    = "cli.ansible.host_vars"
	CliAnsibleRolesCommon = "cli.ansible.roles.common"
	CliAnsiblePlaybooks   = "cli.ansible.playbooks"
	CliAnsibleInventories = "cli.ansible.inventories"
	CliAnsibleFiles       = "cli.ansible.files"
)

type TemplateData struct {
	Description string
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"create"},
	Short:   "Create Ansible Structure",
	Long:    `Create Ansible basic structure within the given path.`,
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ansiblePath, err := environment.ViperGetEnvVariable(CliAnsibleDirectory)

		if err != nil {
			log.Fatal("Error")
		}

		if createPathParameter != "" {
			ansibleProject = createPathParameter
		}

		if err := os.Mkdir(ansiblePath+ansibleProject, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		ansibleProjectPath := ansiblePath + ansibleProject

		groupVars, err := environment.ViperGetEnvVariableSlice(CliAnsibleGroupVars)
		if err != nil {
			log.Fatal("Error")
		}
		addVarFiles(ansibleProjectPath, "group_vars", groupVars)

		hostVars, err := environment.ViperGetEnvVariableSlice(CliAnsibleHostVars)
		if err != nil {
			log.Fatal("Error")
		}
		addVarFiles(ansibleProjectPath, "host_vars", hostVars)

		addFiles(ansibleProjectPath)

		addPlaybooks(ansibleProjectPath)

		addInventory(ansibleProjectPath)

		filePath := filepath.Join(ansibleProjectPath, "ansible.cfg")

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

func addFiles(ansibleProjectPath string) {
	files, err := environment.ViperGetEnvVariableSlice(CliAnsibleFiles)
	if err != nil {
		log.Fatal("Error")
	}

	if err := os.Mkdir(ansibleProjectPath+"/files", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if err := os.Mkdir(ansibleProjectPath+"/files/"+file, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func addPlaybooks(ansibleProjectPath string) {
	playbooks, err := environment.ViperGetEnvVariableSlice(CliAnsiblePlaybooks)

	if err != nil {
		log.Fatal("Error")
	}

	if err := os.Mkdir(ansibleProjectPath+"/playbooks", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, playbook := range playbooks {
		filePath := filepath.Join(ansibleProjectPath+"/playbooks", playbook+".yml")

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

func addInventory(ansibleProjectPath string) {
	inventories, err := environment.ViperGetEnvVariableSlice(CliAnsibleInventories)

	if err != nil {
		log.Fatal("Error")
	}

	if err := os.Mkdir(ansibleProjectPath+"/inventory", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, inventory := range inventories {
		filePath := filepath.Join(ansibleProjectPath+"/inventory", inventory+".ini")

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

func addVarFiles(ansibleProjectPath string, varDir string, vars []string) {
	if err := os.Mkdir(ansibleProjectPath+"/"+varDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, groupVar := range vars {
		filePath := filepath.Join(ansibleProjectPath+"/"+varDir+"/", groupVar+".yml")
		tmpl, err := template.New(groupVar).Parse("#add vars here")
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
	createCmd.Flags().StringVarP(&createPathParameter, "path", "p", "", "Role Path")
	err := createCmd.MarkFlagRequired("path")
	if err != nil {
		return
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
