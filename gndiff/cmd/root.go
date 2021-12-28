/*
Copyright © 2021 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnames/gndiff"
	"github.com/gnames/gndiff/config"
	"github.com/gnames/gndiff/ent/output"
	"github.com/gnames/gndiff/ent/record"
	"github.com/gnames/gndiff/io/ingestio"
	"github.com/gnames/gnfmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var quiet bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gndiff source_file reference_file",
	Short: "Compares two files with scientific names.",
	Long: `
Extracts scientific names, their IDs and families the source and reference
files and prints out a match of a reference data to the source data.

The files can be in comma-separated (CSV), tab-separated (TSV) formats, or
just contain name-strings only (one per line).

TSV/CSV files must contain 'ScientificName' field, 'Family' and 'TaxonID'
fields are also ingested.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var opts []config.Option
		if versionFlag(cmd) {
			os.Exit(0)
		}

		quiet, _ = cmd.Flags().GetBool("quiet")

		opts = append(opts, fmtFlag(cmd))
		cfg := config.New(opts...)
		gnd := gndiff.New(cfg)
		if len(args) != 2 {
			log.Warn("Supply paths to two CSV/TSV files for comparison")
			_ = cmd.Help()
			os.Exit(1)
		}
		src, ref, err := readFiles(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		res, err := gnd.Compare(src, ref)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(output.MatchOutput(res, cfg.Format))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("version", "V", false,
		"shows build version and date, ignores other flags.")

	formatHelp := "Sets output format. Can be one of:\n" +
		"'csv', 'tsv', 'compact', 'pretty'\n" +
		"default is 'csv'."
	rootCmd.Flags().StringP("format", "f", "", formatHelp)

	rootCmd.Flags().BoolP("quiet", "q", false, "Do not output info and warning logs.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func versionFlag(cmd *cobra.Command) bool {
	version, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal(err)
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
			log.Warnf("Cannot recognize format '%s', keeping default 'CSV' format.", f)
		}
		return config.OptFormat(gnfmt.CSV)
	}
}

func readFiles(srcPath, refPath string) ([]record.Record, []record.Record, error) {
	cfg := config.New()
	ing := ingestio.New(cfg)

	src := filepath.Join(srcPath)
	recSrc, err := ing.Records(src)
	if err != nil {
		return nil, nil, err
	}

	ref := filepath.Join(refPath)
	recRef, err := ing.Records(ref)
	if err != nil {
		return nil, nil, err
	}

	return recSrc, recRef, nil
}
