package changelog

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var githubStyleTestLog = `# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Fixed
- something broken
- some issue

### Removed
- some old stuff
- bad code

## [v0.3.0] - 2016-12-03
### Added
- This awesome feature
- More pewpew.

## [v0.2.0] - 2015-10-06
### Changed
- a thingy with some subpoints:
	- this one
	- that one
	- yay!

### Deprecated
- legacy stuff
- args of some function

## [0.1.0] - 2014-09-02
### Security
- hard coded passwords have been removed
- stack overflow issue solved!

[Unreleased]: https://github.com/myuser/myproject/compare/v0.3.0...HEAD
[v0.3.0]: https://github.com/myuser/myproject/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/myuser/myproject/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/myuser/myproject/compare/v0.0.8...v0.1.0
`

var bitbucketStyleTestLog = `# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- add infrastructure resource catalog

[Unreleased]: https://bitbucket.org/helsternware/www.helstern.org/branches/compare/HEAD%0D8765b1b1f001ea69377bf4fe38fadc0953cf248f
`

func TestParse_has_header(t *testing.T) {
	buf := bytes.NewBufferString(githubStyleTestLog)
	contents, err := Parse(buf)
	assert.Nil(t, err)
	header := `# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

`
	assert.Equal(t, header, contents.Header)
}

func TestParse_has_rest(t *testing.T) {
	buf := bytes.NewBufferString(githubStyleTestLog)
	contents, err := Parse(buf)
	assert.Nil(t, err)
	rest := `## [v0.3.0] - 2016-12-03
### Added
- This awesome feature
- More pewpew.

## [v0.2.0] - 2015-10-06
### Changed
- a thingy with some subpoints:
	- this one
	- that one
	- yay!

### Deprecated
- legacy stuff
- args of some function

## [0.1.0] - 2014-09-02
### Security
- hard coded passwords have been removed
- stack overflow issue solved!

`
	assert.Equal(t, rest, contents.Rest)
}

func TestParse_has_changes(t *testing.T) {
	buf := bytes.NewBufferString(githubStyleTestLog)
	contents, err := Parse(buf)
	assert.Nil(t, err)
	expected := []*Changes{
		&Changes{
			Tag: "Unreleased",
			Removed: `- some old stuff
- bad code`,
			Fixed: `- something broken
- some issue`,
		},
		&Changes{
			Tag: "v0.3.0",
			Added: `- This awesome feature
- More pewpew.`,
			Time: time.Date(2016, 12, 3, 0, 0, 0, 0, time.UTC),
		},
		&Changes{
			Tag: "v0.2.0",
			Changed: `- a thingy with some subpoints:
	- this one
	- that one
	- yay!`,
			Deprecated: `- legacy stuff
- args of some function`,
			Time: time.Date(2015, 10, 6, 0, 0, 0, 0, time.UTC),
		},
		&Changes{
			Tag: "0.1.0",
			Security: `- hard coded passwords have been removed
- stack overflow issue solved!`,
			Time: time.Date(2014, 9, 2, 0, 0, 0, 0, time.UTC),
		},
	}
	assert.Equal(t, expected, contents.Changes)
}

func TestParse_has_references(t *testing.T) {
	buf := bytes.NewBufferString(githubStyleTestLog)
	contents, err := Parse(buf)
	assert.Nil(t, err)
	expected := []Reference{
		{
			Type:      GITHUB_REFERENCE,
			Tag:       "Unreleased",
			Raw:       "[Unreleased]: https://github.com/myuser/myproject/compare/v0.3.0...HEAD",
			From:      "v0.3.0",
			To:        "HEAD",
			Separator: "...",
			BaseURL:   "https://github.com/myuser/myproject",
		},
		{
			Type:      GITHUB_REFERENCE,
			Tag:       "v0.3.0",
			Raw:       "[v0.3.0]: https://github.com/myuser/myproject/compare/v0.2.0...v0.3.0",
			From:      "v0.2.0",
			To:        "v0.3.0",
			Separator: "...",
			BaseURL:   "https://github.com/myuser/myproject",
		},
		{
			Type:      GITHUB_REFERENCE,
			Tag:       "v0.2.0",
			Raw:       "[v0.2.0]: https://github.com/myuser/myproject/compare/v0.1.0...v0.2.0",
			From:      "v0.1.0",
			To:        "v0.2.0",
			Separator: "...",
			BaseURL:   "https://github.com/myuser/myproject",
		},
		{
			Type:      GITHUB_REFERENCE,
			Tag:       "0.1.0",
			Raw:       "[0.1.0]: https://github.com/myuser/myproject/compare/v0.0.8...v0.1.0",
			From:      "v0.0.8",
			To:        "v0.1.0",
			Separator: "...",
			BaseURL:   "https://github.com/myuser/myproject",
		},
	}
	assert.Equal(t, expected, contents.Refs)
}

func TestRender_github(t *testing.T) {
	in := bytes.NewBufferString(githubStyleTestLog)
	contents, err := Parse(in)
	assert.Nil(t, err)

	out := bytes.NewBuffer([]byte{})
	_, err = contents.WriteTo(out)
	assert.Nil(t, err)

	assert.Equal(t, githubStyleTestLog, out.String())
}
func TestRender_bitbucket(t *testing.T) {
	in := bytes.NewBufferString(bitbucketStyleTestLog)
	contents, err := Parse(in)
	assert.Nil(t, err)

	out := bytes.NewBuffer([]byte{})
	_, err = contents.WriteTo(out)
	assert.Nil(t, err)

	assert.Equal(t, bitbucketStyleTestLog, out.String())
}
