package genjs

import (
	"time"

	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/meta"
)

var (
	// Timeout is the request timeout before it gets aborted.
	Timeout time.Duration

	// Scheme is the URL scheme used to make requests to the API.
	Scheme string

	// Host is the API hostname.
	Host string
)

// Command is the goa application code generator command line data structure.
// It implements meta.Command.
type Command struct {
	*codegen.BaseCommand
}

// NewCommand instantiates a new command.
func NewCommand() *Command {
	base := codegen.NewBaseCommand("js", "Generate javascript client module")
	return &Command{BaseCommand: base}
}

// RegisterFlags registers the command line flags with the given registry.
func (c *Command) RegisterFlags(r codegen.FlagRegistry) {
	r.Flags().DurationVar(&Timeout, "timeout", time.Duration(20)*time.Second, `the duration before the request times out.`)
	r.Flags().StringVar(&Scheme, "scheme", "http", `the URL scheme used to make requests to the API, defaults to the scheme defined in the API design if any.`)
	r.Flags().StringVar(&Host, "host", "", `the API hostname, defaults to the hostname defined in the API design if any`)
}

// Run simply calls the meta generator.
func (c *Command) Run() ([]string, error) {
	gen := meta.NewGenerator(
		"genjs.Generate",
		[]*codegen.ImportSpec{codegen.SimpleImport("github.com/goadesign/goa/goagen/gen_js")},
		nil,
	)
	return gen.Generate()
}
