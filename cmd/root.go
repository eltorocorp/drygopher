package cmd

import (
	"errors"
	"os"

	"github.com/eltorocorp/drygopher/internal/coverage"
	"github.com/eltorocorp/drygopher/internal/coverage/analysis"
	"github.com/eltorocorp/drygopher/internal/coverage/packages"
	"github.com/eltorocorp/drygopher/internal/coverage/profile"
	"github.com/eltorocorp/drygopher/internal/coverage/report"
	"github.com/eltorocorp/drygopher/internal/host"
	wordwrap "github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var (
	coverageStandard     float64
	profileName          string
	suppressProfile      bool
	exclusionPatterns    []string
	useDefaultExclusions bool
)

const examples = `
Run the command bare with no flags or arguments. In this case, it will analyize coverage of all packages below the current directory.

  $ drygopher

Lower the coverage standard to 98.2% and change the name of the coverage profile file.

  $ drygopher -s 98.2 -p coveragedata.txt

Suppress creating the coverage profile file, and manually exclude vendored packages from coverage analysis.

  $ drygopher --suppressprofile -e '/vendor/' '_test'

Run coverage analysis excluding vendor and test packages, and also exclude any packages whose name ends with "service".

  $drygopher -d -e service$
`

var rootCmd = &cobra.Command{
	Use:   "drygopher [flags]",
	Short: "Keep your coverage high, and your gopher dry.",
	Long:  wordwrap.WrapString("\ndrygopher provides coverage analysis for go projects. It keeps your gopher dry by making sure everything is covered as it should be. Visit http://github.com/eltorocorp/drygopher for more information.", 80),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		packageExclusions := exclusionPatterns
		if useDefaultExclusions {
			packageExclusions = append(packageExclusions, "/vendor/", "_test")
		}

		execAPI := new(host.Exec)
		osioAPI := new(host.OSIO)
		packageAPI := packages.New(execAPI, osioAPI)
		profileAPI := profile.New(packageAPI, osioAPI)
		reportAPI := report.New(execAPI)
		analysisAPI := analysis.New(osioAPI)
		coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

		return coverageAPI.AnalyzeUnitTestCoverage(exclusionPatterns, coverageStandard, suppressProfile, profileName)
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if coverageStandard < 0 {
			err = errors.New("coverage standard must not be negative")
		} else if profileName == "" {
			err = errors.New("profilename must not be an empty string")
		}
		return
	},
	Example: wordwrap.WrapString(examples, 80),
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	wrap := func(s string) string {
		return wordwrap.WrapString(s, 50)
	}
	rootCmd.Flags().Float64VarP(&coverageStandard, "standard", "s", 100, wrap("Coverage standard to use."))
	rootCmd.Flags().StringVarP(&profileName, "profilename", "p", "coverage.out", wrap("The name of the coverage profile file. This flag has no effect if the suppressprofile flag is also set."))
	rootCmd.Flags().BoolVar(&suppressProfile, "suppressprofile", false, wrap("Supply this flag to suppress creating the coverage profile file."))
	rootCmd.Flags().StringSliceVarP(&exclusionPatterns, "exclusions", "e", []string{}, wrap("A set of regular expressions used to define packages to exclude from coverage analysis. This flag can be combined with the defaultexclusions flag."))
	rootCmd.Flags().BoolVarP(&useDefaultExclusions, "defaultexclusions", "d", false, wrap("Exclude vendor and _test packages from coverage analysis. This flag can be combined with the exclusions flag."))

	rootCmd.DisableFlagsInUseLine = true
	rootCmd.SilenceUsage = true
}
