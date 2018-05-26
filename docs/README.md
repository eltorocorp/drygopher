# drygopher
Keep your coverage high, and your gopher dry.

[![Go Report Card](https://goreportcard.com/badge/github.com/eltorocorp/drygopher)](https://goreportcard.com/report/github.com/eltorocorp/drygopher)

*drygopher* is yet another coverage analysis tool for go.

Its purpose is to keep your gopher dry by ensuring that that all of the code in your project is covered by tests, according to conventions you specify. 

## OK, how is this different?
"How will this keep my gopher more dry than the other go coverage tools out there", you ask?
Good question, Jimmy.

There are a lot of go coverage tools available, some better than others. Here's why *drygopher* will keep your gopher dry and safe:

* Cross-Package Coverage: Calculate coverage across multiple packages (even those with no associated tests).
* Consolidated Stats: Consolidates cover profile data into a single file.
* Convention Over Code: Setup conventions for excluding packages from test coverage.
* Coverage Report: Output coverage report to stdout (more friendly than raw cover profile).
* Set Your Standard: *drygopher* assumes you want 100% coverage by default, but you can override this with ease.

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
