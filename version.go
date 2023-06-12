package binutil

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 1
	VersionMinor         = 0
	VersionPatch         = 0
	VersionReleaseLevel  = "alpha"
	VersionReleaseNumber = 1
)

// The GitVersion and BuildDate are set via ldflags
// Set the GitVersion via -ldflags="-X 'github.com/bbengfort/binutil.GitVersion=$(git rev-parse --short HEAD)'"
// Set the BuildDate via -ldflags="-X 'github.com//bbengfort/binutil.BuildDate=$(date +%F)'"
var (
	GitVersion string
	BuildDate  string
)

// Version returns the semantic version for the current build.
func Version() string {
	var versionCore string
	if VersionPatch > 0 || VersionReleaseLevel != "" {
		versionCore = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	} else {
		versionCore = fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
	}

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			versionCore = fmt.Sprintf("%s-%s.%d", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		} else {
			versionCore = fmt.Sprintf("%s-%s", versionCore, VersionReleaseLevel)
		}
	}

	if GitVersion != "" {
		if BuildDate != "" {
			versionCore = fmt.Sprintf("%s (revision %s built on %s)", versionCore, GitVersion, BuildDate)
		} else {
			versionCore = fmt.Sprintf("%s (%s)", versionCore, GitVersion)
		}
	}

	return versionCore
}
