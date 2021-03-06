package uconfig

import (
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/file"
	"github.com/omeid/uconfig/plugins/flag"
)

// Files represents a set of file paths and the appropriate
// unmarshal function for the given file.
type Files map[string]file.Unmarshal

// Classic creates a uconfig manager with defaults,environment variables,
// and flags (in that order) and optionally file loaders based on the provided
// Files map and parses them right away.
func Classic(conf interface{}, files Files) (Config, error) {
	c, err := New(conf,
		defaults.New(),
	)

	if err != nil {
		return nil, err
	}

	for path, unmarshal := range files {
		err = c.Walker(file.New(path, unmarshal))
		if err != nil {
			return nil, err
		}
	}

	err = c.Visitor(env.New())
	if err != nil {
		return nil, err
	}
	err = c.Visitor(flag.Standard())
	if err != nil {
		return c, err
	}
	return c, c.Parse()
}
