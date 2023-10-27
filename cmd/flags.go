package cmd

import (
	"fmt"
	"log/slog"
	"os"

	gndiff "github.com/gnames/gndiff/pkg"
	"github.com/gnames/gndiff/pkg/config"
	"github.com/gnames/gnfmt"
	"github.com/spf13/cobra"
)

func versionFlag(cmd *cobra.Command) bool {
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if version {
		fmt.Printf("\nversion: %s\n\nbuild:   %s\n\n",
			gndiff.Version, gndiff.Build)
		return true
	}
	return false
}

func fmtFlag(cmd *cobra.Command) config.Option {
	f, err := cmd.Flags().GetString("format")
	if f == "" || err != nil {
		return config.OptFormat(gnfmt.CSV)
	}

	switch f {
	case "csv":
		return config.OptFormat(gnfmt.CSV)
	case "tsv":
		return config.OptFormat(gnfmt.TSV)
	case "compact":
		return config.OptFormat(gnfmt.CompactJSON)
	case "pretty":
		return config.OptFormat(gnfmt.PrettyJSON)
	default:
		if !quiet {
			slog.Warn(
				"Cannot recognize format-string, keeping default 'CSV' format.",
				slog.String("format-string", f),
			)
		}
		return config.OptFormat(gnfmt.CSV)
	}
}
