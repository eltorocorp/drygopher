// Package cmd is responsible for the CLI's user interface.
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

Run coverage analysis, excluding default packages (vendor and test), and view the resulting coverage heatmap. Note that ';' is used rather than '&&' between the commands to ensure that 'go tool' is run even if drygopher complains that coverage is below standard.

  $ drygopher -d; go tool cover -html=coverage.out

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

Use the 'coveragepct' file, which contains the coverage percentage calculated by drygopher, as a parameter in some other command. This example shows a trivial 'echo' command that returns 'Coverage: 100%'.

  $ drygopher -d; echo Coverage: $(cat coveragepct)%

Notes:

 - When drygopher encounters a package that has no associated unit tests, it creates an estimate for the number of statements that the package might contain. This estimate is the average (median) number of statements found across packages that do have coverage. In such cases, drygopher does not go any further than making an estimate, as it does not want to make assumptions about the code authors' intent. By the same token, it would be generally incorrect to presume that the package contained no statements at all. This is merely a means to an end until the author decides to either a) cover the package with unit tests, or b) exclude the package from coverage analysis entirely.

 - drygopher will exit with a non-zero status code (and display an informational error message) if the calculated coverage falls below the coverage standard. To override this behavior, either increase code coverage, or reduce the coverage standard (using the -s flag).
 
 - drygopher will exit with a non-zero status code (and display an informational error message) if any unit tests fail during analysis. Unit test failures have a higher priority than coverage failures. Thus, if a system under test fails for both coverage and unit test failure, drygopher will report the unit test failure rather than the coverage failure.

 - drygopher will exit with a non-zero status code and little other information if build errors are encountered while testing.
 
 - When drygopher tests each package, it applies the -race flag. This decreases testing performance, but ensures that applicable race conditions are captured during analysis.
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
