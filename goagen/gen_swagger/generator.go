package genswagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/utils"
	"github.com/spf13/cobra"
)

// Generator is the swagger code generator.
type Generator struct{}

// Generate is the generator entry point called by the meta generator.
func Generate(roots []interface{}) (files []string, err error) {
	api := roots[0].(*design.APIDefinition)
	g := new(Generator)
	root := &cobra.Command{
		Use:   "goagen",
		Short: "Swagger generator",
		Long:  "Swagger generator",
		Run:   func(*cobra.Command, []string) { files, err = g.Generate(api) },
	}
	codegen.RegisterFlags(root)
	NewCommand().RegisterFlags(root)
	root.Execute()
	return
}

// Generate produces the skeleton main.
func (g *Generator) Generate(api *design.APIDefinition) (_ []string, err error) {
	var genfiles []string

	cleanup := func() {
		for _, f := range genfiles {
			os.Remove(f)
		}
	}

	go utils.Catch(nil, cleanup)

	defer func() {
		if err != nil {
			cleanup()
		}
	}()

	s, err := New(api)
	if err != nil {
		return
	}
	b, err := json.Marshal(s)
	if err != nil {
		return
	}
	swaggerDir := filepath.Join(codegen.OutputDir, "swagger")
	os.RemoveAll(swaggerDir)
	if err = os.MkdirAll(swaggerDir, 0755); err != nil {
		return
	}
	genfiles = append(genfiles, swaggerDir)
	swaggerFile := filepath.Join(swaggerDir, "swagger.json")
	err = ioutil.WriteFile(swaggerFile, b, 0644)
	if err != nil {
		return
	}
	genfiles = append(genfiles, swaggerFile)
	controllerFile := filepath.Join(swaggerDir, "swagger.go")
	genfiles = append(genfiles, controllerFile)
	file, err := codegen.SourceFileFor(controllerFile)
	if err != nil {
		return
	}
	imports := []*codegen.ImportSpec{
		codegen.SimpleImport("github.com/goadesign/goa"),
	}
	file.WriteHeader(fmt.Sprintf("%s Swagger Spec", api.Name), "swagger", imports)
	file.Write([]byte(swagger))
	if err = file.FormatCode(); err != nil {
		return
	}

	return genfiles, nil
}

const swagger = `
// MountController mounts the swagger spec controller under "/swagger.json".
func MountController(service goa.Service) {
	service.ServeFiles("/swagger.json", "swagger/swagger.json")
}
`
