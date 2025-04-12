package display

import (
	"fmt"
	"io"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
	"github.com/xlab/treeprint"
)

const (
	dateFormat = "2006-01-02 15:04:05"
)

// PrintFiles writes the list of files in tree format to the given writer
func PrintFiles(w io.Writer, files []cvedb.File) error {
	tree := treeprint.New()
	tree.SetValue("Files")
	for _, file := range files {
		fileSubBranch := tree.AddBranch(fileEmoji + " " + file.Name)
		fileSubBranch.AddNode(sizeEmoji + " " + file.PrettySize)
		fileSubBranch.AddNode(dateEmoji + " " + file.ModifiedDate.Format(dateFormat))
	}

	_, err := fmt.Fprintln(w, tree.String())
	return err
}
