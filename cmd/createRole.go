/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang-ansible/pkg/environment"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// createRoleCmd represents the createRole command
var createRoleCmd = &cobra.Command{
	Use:   "createRole",
	Short: "Create Ansible Role Structure",
	Long:  `Create Ansible Role basic structure within the given path.`,
	Run: func(cmd *cobra.Command, args []string) {
		ansiblePath, err := environment.ViperGetEnvVariable(CliAnsibleDirectory)

		if err != nil {
			log.Fatal("Error1")
		}

		if len(args) >= 1 && args[0] != "" {
			ansibleProject = args[0]
		}

		if err := os.Mkdir(ansiblePath+ansibleProject, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		ansibleProjectPath := ansiblePath + ansibleProject
		addDirWithEmptyFile(ansibleProjectPath, "vars", `#empty file`)
		addDirWithEmptyFile(ansibleProjectPath, "tasks", `#empty file`)
		addDirWithEmptyFile(ansibleProjectPath, "meta",
			`dependencies: []

galaxy_info:
  role_name: 
  author: 
  description: 
  company: 
  license: 
  min_ansible_version: 
  platforms:
    - name: Platform
      version:
        - all
  galaxy_tags:
    - tags
`)
		addDirWithEmptyFile(ansibleProjectPath, "handlers", `#empty file`)
		addDirWithEmptyFile(ansibleProjectPath, "defaults", `#empty file`)

		if err := os.Mkdir(ansibleProjectPath+"/molecule", os.ModePerm); err != nil {
			log.Fatal(err)
		}

		addDirWithEmptyFile(ansibleProjectPath, "molecule/default", `#empty file`, "converge", "molecule")
		fmt.Println("createRole called")
	},
}

func addDirWithEmptyFile(ansibleProjectPath string, directory string, text string, fileNames ...string) {
	if err := os.Mkdir(ansibleProjectPath+"/"+directory, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if len(fileNames) <= 0 {
		fileNames = []string{"main"}
	}

	for _, fileName := range fileNames {
		filePath := filepath.Join(ansibleProjectPath+"/"+directory, fileName+".yml")
		tmpl, err := template.ParseFiles("./templates/template.txt")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal("Error3")
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal("Error4")
			}
		}(file)

		data := map[string]interface{}{
			"Text": text,
		}

		err = tmpl.Execute(file, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(createRoleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createRoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createRoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
