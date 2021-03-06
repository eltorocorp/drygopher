# drygopher
Keep your coverage high, and your gopher dry.

[![Go Report Card](https://goreportcard.com/badge/github.com/eltorocorp/drygopher)](https://goreportcard.com/report/github.com/eltorocorp/drygopher)
[![Build Status](http://badges.awsp.eltoro.com?project=drygopher&item=build)](http://github.com/eltorocorp/drygopher)
[![Coverage](http://badges.awsp.eltoro.com?project=drygopher&item=coverage)](http://github.com/eltorocorp/drygopher)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/eltorocorp/drygopher/drygopher)

*drygopher* is yet another coverage analysis tool for go.

Its purpose is to keep your gopher dry by ensuring that that all of the code in your project is covered by tests, according to conventions you specify. 

## OK, how is this different?
"How will this keep my gopher more dry than the other go coverage tools out there", you ask?
Good question, Jimmy.

There are a lot of go coverage tools available, some better than others. Here's why *drygopher* will keep your gopher dry and safe:

* [Cross-Package Coverage](#cross-package-coverage): Calculate coverage across multiple packages (even those with no associated tests).
* [Consolidated Stats](#consolidated-stats): Consolidates cover profile data into a single file.
* [Convention Over Code](#convention-over-code): Setup conventions for excluding packages from test coverage.
* [Coverage Report](#coverage-report): Output coverage report to stdout (more friendly than raw cover profile).
* [Set Your Standard](#set-your-standard): *drygopher* assumes you want 100% coverage by default, but you can override this with ease.

## To Install

```
$ go install github.com/eltorocorp/drygopher/drygopher
```

## Basic Usage

```
Usage:
  drygopher [flags]

Examples:

drygopher provides coverage analysis for go projects. It keeps your gopher dry
by making sure everything is covered as it should be. Visit
http://github.com/eltorocorp/drygopher for more information.

Usage:
  drygopher [flags]

Examples:

Analyze coverage of all packages below the current directory, expecting 100%
coverage, and require that all packages participate in coverage analysis
(regardless of if they have/need tests). No gopher needs to be this dry.

  $ drygopher

Lower the coverage standard to 98.2% and change the name of the coverage profile
file.

  $ drygopher -s 98.2 -p coveragedata.txt

Run coverage analysis, excluding vendor and test packages, and suppress the
generation of a coverage profile.

  $ drygopher -d --suppressprofile

Run coverage analysis, excluding default packages (vendor and test), and view
the resulting coverage heatmap. Note that ';' is used rather than '&&' between
the commands to ensure that 'go tool' is run even if drygopher complains that
coverage is below standard.

  $ drygopher -d; go tool cover -html=coverage.out

Run coverage analysis, excluding vendor and test packages, and also exclude any
packages whose name ends with "service". Note that in this case, we enclose the
expression in single quotes to prevent globbing.

  $drygopher -d -e 'service$'

Run coverage analysis, excluding vendor and test packages, and packages that end
in cmd, or iface, or contain mock anywhere in the name.
The following commands are all equivalent:

  Using defaults plus a comma separated list of expressions:
  $drygopher -d -e "'cmd$','iface$',mock"

  Using defaults and explicit expressions:
  $drygopher -d -e 'cmd$' -e 'iface$' -e mock

  Using groups of explicit expressions:
  $drygopher -e "/vendor/,_test" -e "'cmd$','iface$'" -e mock

  Using defaults and a single expression:
  $drygopher -d -e "'cmd$|iface$|mock'"

  Note that when supplying a list of expressions for -e, the list must be comma
delimited. As such, literal commas cannot be used when supplying a list of
expressions for the -e flag. Generally, this shouldn't be an issue since commas
are not typically valid in package names.

Use the 'coveragepct' file, which contains the coverage percentage calculated by
drygopher, as a parameter in some other command. This example shows a trivial
'echo' command that returns 'Coverage: 100%'.

  $ drygopher -d; echo Coverage: $(cat coveragepct)%

Notes:

 - When drygopher encounters a package that has no associated unit tests, it
creates an estimate for the number of statements that the package might contain.
This estimate is the average (median) number of statements found across packages
that do have coverage. In such cases, drygopher does not go any further than
making an estimate, as it does not want to make assumptions about the code
authors' intent. By the same token, it would be generally incorrect to presume
that the package contained no statements at all. This is merely a means to an
end until the author decides to either a) cover the package with unit tests, or
b) exclude the package from coverage analysis entirely.

 - drygopher will exit with a non-zero status code (and display an informational
error message) if the calculated coverage falls below the coverage standard. To
override this behavior, either increase code coverage, or reduce the coverage
standard (using the -s flag).

 - drygopher will exit with a non-zero status code (and display an informational
error message) if any unit tests fail during analysis. Unit test failures have a
higher priority than coverage failures. Thus, if a system under test fails for
both coverage and unit test failure, drygopher will report the unit test failure
rather than the coverage failure.

 - drygopher will exit with a non-zero status code and little other information
if build errors are encountered while testing.

 - When drygopher tests each package, it applies the -race flag. This decreases
testing performance, but ensures that applicable race conditions are captured
during analysis.


Flags:
  -d, --defaultexclusions    Exclude vendor and _test packages from coverage
                             analysis. This flag can be combined with the
                             exclusions flag.
  -e, --exclusions strings   A set of regular expressions used to define
                             packages to exclude from coverage analysis. This
                             flag can be combined with the defaultexclusions
                             flag.
  -h, --help                 help for drygopher
  -p, --profilename string   The name of the coverage profile file. This flag
                             has no effect if the suppressprofile flag is also
                             set. (default "coverage.out")
  -s, --standard float       Coverage standard to use. (default 100)
      --suppresspctfile      Suppress the creation of the coverarage percentage
                             file ('coveragepct').
      --suppressprofile      Supply this flag to suppress creating the coverage
                             profile file.
```

## Details
### Cross-Package Coverage
The native go tooling ([go test](https://golang.org/cmd/go/#hdr-Test_packages)) is unable to build coverage statistics for more than one package at a time. Perhaps some day this will change. 
Other tools such as [axw/gocov](https://github.com/axw/gocov), [vieux/gocover.io](https://github.com/vieux/gocover.io),  [hay14busa/goverage](https://github.com/vieux/gocover.io), [dave/courtney](https://github.com/dave/courtney), and [others](https://github.com/search?l=Go&o=desc&p=1&q=go+coverage&s=stars&type=Repositories) all offer some form of cross-package coverage, but have limitations.

Typically, the greatest limitation in these other packages is that they will only calculate coverage for packages that already have at least one test defined. Due to the way go's test tooling works, the `go test` command will only set out to count the number of covered statements in a package if that package has at least one test defined. The drawback of this (particularly in an enterprise environment with distributed teams working on projects at high velocity) is that you can't rely on the native tooling to tell you that a package that should have tests does not have tests. This can leave your gopher damp in places where you didn't expect it to be wet. Nobody likes a damp gopher.

*drygopher* overcomes this issue with a simple, but effective, heuristic. *drygopher* by default, assumes that all packages in your project must be covered, and it builds coverage statistics even for those packages that contain no test files. With this, you can quickly identify packages that have zero tests, right along side packages that have some tests, and have confidence that you are in control over your gopher's hydration levels.

### Convention Over Code

There are situations where certain packages should, rightly, be excluded from coverage analysis. Candidates for such exclusions often include:
* service/program entrypoints
* generated code
* test packages themselves
* packages that define types but contain no executable statements
* utilities that never see production use
* etc...

It is entirely reasonable to exclude packages from code coverage when the idea of unit testing such packages doesn't make sense in the first place.

With that in mind, *drygopher* allows you to supply naming conventions that your project follows, so that *drygopher* can exclude packages that don't require coverage from the coverage standards. This concept follows the convention over code philosophy. *drygopher* will automatically know not to test certain packages by following the conventions that you have setup in your project. 

Ooo, that gopher is looking good. Well moisturized, but not saturated; a true beauty, that gopher is.

### Consolidated Stats
*drygopher* will consolidate coverage statistics for each package that it believes should have associated tests. This cover profile is formatted the same way native go profiles are, and can be seemlessly consumed by `go tool cover`.

### Coverage Report
*drygopher* outputs a human-readable coverage report to stdout while it is testing your project. This report contains information about:
* which packages are included and excluded from coverage analysis
* how many statements are covered or not covered in each package
* the coverage percentage for each package
* the statement count and coverage percentages across all tested packages.

The coverage report is very helpful when printed locally, or when included in CI build output.

### Set Your Standard
*drygopher's* default behavior is to assume you want 100% coverage. However, you can set any standard you want. This is helpful in situations where you are either fine having a slightly wet gopher, or when you're in a process of drying out a gopher over time and want to start low and work your way up to full dryness.
