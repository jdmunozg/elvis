package create

import (
	"fmt"

	"github.com/cgalvisleon/elvis/utilities"
	"github.com/spf13/cobra"
)

var CmdProject = &cobra.Command{
	Use:   "micro [name author schema, schema_var]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utilities.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := prompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		author, err := prompStr("Author")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkPMicroservice(packageName, name, author, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

var CmdMicro = &cobra.Command{
	Use:   "micro [name schema, schema_var]",
	Short: "Create project base type microservice.",
	Long:  "Template project to microservice include folder cmd, deployments, pkg, rest, test and web, with files .go required for making a microservice.",
	Run: func(cmd *cobra.Command, args []string) {
		packageName, err := utilities.ModuleName()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		name, err := prompStr("Name")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMicroservice(packageName, name, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}

var CmdModelo = &cobra.Command{
	Use:   "modelo [name modelo, schema]",
	Short: "Create model to microservice.",
	Long:  "Template model to microservice include function handler model.",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := prompStr("Package")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		modelo, err := prompStr("Model")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		schema, err := prompStr("Schema")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		err = MkMolue(name, modelo, schema)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}

		fmt.Println("Remember, including the router in router.go")
	},
}

var CmdRpc = &cobra.Command{
	Use:   "rpc [name]",
	Short: "Create rpc model to microservice.",
	Long:  "Template rpc model to microservice include function handler model.",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := prompStr("Package")
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		
		err = MkRpc(name)
		if err != nil {
			fmt.Printf("Command failed %v\n", err)
			return
		}
	},
}