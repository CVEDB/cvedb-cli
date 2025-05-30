package display

import (
	"fmt"
	"io"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/xlab/treeprint"
)

// PrintScripts writes the scripts list in tree format to the given writer
func PrintScripts(w io.Writer, scripts []cvedb.Script) error {
	tree := treeprint.New()
	tree.SetValue("Scripts")
	for _, script := range scripts {
		branch := tree.AddBranch(script.Name)
		branch.AddNode(script.Script.Source)
	}

	_, err := fmt.Fprintln(w, tree.String())
	return err
}
