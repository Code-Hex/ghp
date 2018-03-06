package ghp

import (
	"bytes"
	"fmt"
	"os"

	"reflect"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

const (
	name   = "ghp"
	indent = "        "
)

// Options struct for parse command line arguments
type Options struct {
	Help bool `short:"h" long:"help" description:"show this message"`

	WithLicense bool `long:"with-license" description:"create project with a license"`
	StackTrace  bool `long:"trace" description:"display detail error messages"`
}

func (opts *Options) parse(argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.None)
	args, err := p.ParseArgs(argv)
	if err != nil {
		os.Stderr.WriteString(opts.usage())
		return nil, errors.Wrap(err, "invalid command line options")
	}
	return args, nil
}

func (opts Options) usage() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `Usage: %s [options] [PROJECT]
Options:
`, name)

	t := reflect.TypeOf(opts)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		desc := tag.Get("description")
		var o string
		if s := tag.Get("short"); s != "" {
			o = fmt.Sprintf("-%s, --%s", tag.Get("short"), tag.Get("long"))
		} else {
			o = fmt.Sprintf("--%s", tag.Get("long"))
		}
		fmt.Fprintf(&buf, "  %-21s %s\n", o, desc)

		if deflt := tag.Get("default"); deflt != "" {
			fmt.Fprintf(&buf, "  %-21s default: --%s='%s'\n", indent, tag.Get("long"), deflt)
		}
	}

	return buf.String()
}
