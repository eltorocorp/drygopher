package cmd

import (
	"errors"
	"os"

	"github.com/eltorocorp/drygopher/drygopher/coverage"
	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis"
	"github.com/eltorocorp/drygopher/drygopher/coverage/analysis/raw"
	"github.com/eltorocorp/drygopher/drygopher/coverage/host"
	"github.com/eltorocorp/drygopher/drygopher/coverage/packages"
	"github.com/eltorocorp/drygopher/drygopher/coverage/profile"
	"github.com/eltorocorp/drygopher/drygopher/coverage/report"
	wordwrap "github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var (
	coverageStandard       float64
	profileName            string
	suppressProfile        bool
	exclusionPatterns      []string
	useDefaultExclusions   bool
	suppressPercentageFile bool
)

const examples = `
Analyze coverage of all packages below the current directory, expecting 100% coverage, and require that all packages participate in coverage analysis (regardless of if they have/need tests). No gopher needs to be this dry. 

  $ drygopher

Lower the coverage standard to 98.2% and change the name of the coverage profile file.

  $ drygopher -s 98.2 -p coveragedata.txt

Run coverage analysis, excluding vendor and test packages, and suppress the generation of a coverage profile.

  $ drygopher -d --suppressprofile

Run coverage analysis, excluding vendor and test packages, and also exclude any packages whose name ends with "service". Note that in this case, we enclose the expression in single quotes to prevent globbing.

  $drygopher -d -e 'service$'

Run coverage analysis, excluding vendor and test packages, and packages that end in cmd, or iface, or contain mock anywhere in the name.
The following commands are all equivalent:

  Using defaults plus a comma separated list of expressions:
  $drygopher -d -e "'cmd$','iface$',mock"

  Using defaults and explicit expressions:
  $drygopher -d -e 'cmd$' -e 'iface$' -e mock

  Using groups of explicit expressions:
  $drygopher -e "/vendor/,_test" -e "'cmd$','iface$'" -e mock

  Using defaults and a single expression:
  $drygopher -d -e "'cmd$|iface$|mock'"

Note that when supplying a list of expressions for -e, the list must be comma delimited. As such, literal commas cannot be used when supplying a list of expressions for the -e flag. Generally, this shouldn't be an issue since commas are not typically valid in package names.
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
		rawAPI := raw.New(osioAPI, execAPI)
		analysisAPI := analysis.New(rawAPI)
		coverageAPI := coverage.New(packageAPI, analysisAPI, profileAPI, reportAPI)

		return coverageAPI.AnalyzeUnitTestCoverage(packageExclusions, coverageStandard, suppressProfile, profileName, suppressPercentageFile)
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
	rootCmd.Flags().BoolVar(&suppressPercentageFile, "suppresspctfile", false, wrap("Suppress the creation of the coverarage percentage file ('coveragepct')."))

	rootCmd.DisableFlagsInUseLine = true
	rootCmd.SilenceUsage = true
}
