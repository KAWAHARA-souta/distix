package command

import (
	"fmt"
	"os"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/distix-pj/distix/data/extractor"
	// "github.com/distix-pj/distix/data/model"
	"github.com/distix-pj/distix/format"
)


type DistSystemRunner struct {
	RpmDb string
}
var distSysRpmDb string

func (r *DistSystemRunner) Setup() error {
	absRpmDbPath, err := filepath.Abs(distSysRpmDb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	r.RpmDb = absRpmDbPath

	slog.Debug("PackageRunner Options: ",
		"options", r,
		"RpmDb", distSysRpmDb,
	)
	return nil
}

func (r *DistSystemRunner) Run() error {
	ext := extractor.NewRpmdbExtractor(r.RpmDb)
	sysData, err := ext.Extract()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	doc, subdocs, err := sysData.Convert2MultiProtobomDocument()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	w, err := format.NewProtobomWriter(doc, RootOpts.SbomType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	ws := make([]*format.ProtobomWriter, len(subdocs))
	for i, subdoc := range subdocs {
		ws[i], err = format.NewProtobomWriter(subdoc, RootOpts.SbomType)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
	}

	if w.Write(RootOpts.OutputFile) != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}


func NewDistSystemCmd() *cobra.Command {
	runner := &DistSystemRunner{}
	cmd := &cobra.Command{
		Use:   "distsystem",
		Short: "Generate SBOM for system using external references to each package SBOM.",
		// Long: `A longer description here..`
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runner.Setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Run()
		},
	}
	cmd.Flags().StringVarP(&distSysRpmDb, "rpmdb", "r", DEF_RPMDB_PATH, "Path to target RPM Package")
	return cmd
}

func init() {
	RegisterSubCommand(NewDistSystemCmd())
}
