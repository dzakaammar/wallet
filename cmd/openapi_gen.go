package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"wallet/internal/rest"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var openapiGenCmd = &cobra.Command{
	RunE: runOpenAPIGen,
	Use:  "openapi-gen",
}

func init() {
	openapiGenCmd.PersistentFlags().StringP("path", "d", "", "path of directory")
	rootCmd.AddCommand(openapiGenCmd)
}

func runOpenAPIGen(cmd *cobra.Command, _ []string) error {
	pathDir, _ := cmd.PersistentFlags().GetString("path")
	if pathDir == "" {
		return fmt.Errorf("path is required")
	}

	swagger := rest.NewOpenAPI3()

	// openapi3.json
	data, err := json.Marshal(&swagger)
	if err != nil {
		return fmt.Errorf("couldn't marshal json: %s", err)
	}

	if err := os.WriteFile(path.Join(pathDir, "openapi3.json"), data, 0600); err != nil {
		return fmt.Errorf("couldn't write json: %s", err)
	}

	// openapi3.yaml
	data, err = yaml.Marshal(&swagger)
	if err != nil {
		return fmt.Errorf("couldn't marshal json: %s", err)
	}

	if err := os.WriteFile(path.Join(pathDir, "openapi3.yaml"), data, 0600); err != nil {
		return fmt.Errorf("couldn't write json: %s", err)
	}

	return nil
}
