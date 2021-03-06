package impi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type VerifierTestSuite struct {
	suite.Suite
	verifier *verifier
	options  VerifyOptions
}

type verificationTestCase struct {
	name                    string
	contents                string
	expectedErrorStrings    []string
	nonExpectedErrorStrings []string
}

func (s *VerifierTestSuite) SetupTest() {
	var err error

	s.verifier, err = newVerifier()
	s.Require().NoError(err)
}

func (s *VerifierTestSuite) verify(contents string) error {
	return s.verifier.verify(strings.NewReader(contents), &s.options)
}

func (s *VerifierTestSuite) verifyTestCases(verificationTestCases []verificationTestCase) {
	for _, verificationTestCase := range verificationTestCases {
		err := s.verify(verificationTestCase.contents)

		if verificationTestCase.expectedErrorStrings == nil {
			s.Require().NoError(err, verificationTestCase.name)
			continue
		}

		for _, expectedErrorStrings := range verificationTestCase.expectedErrorStrings {
			s.Require().Error(err, verificationTestCase.name)
			s.Require().Contains(err.Error(), expectedErrorStrings, verificationTestCase.name)
		}

		for _, nonExpectedErrorStrings := range verificationTestCase.nonExpectedErrorStrings {
			if err != nil {
				s.Require().NotContains(err.Error(), nonExpectedErrorStrings, verificationTestCase.name)
			}
		}
	}
}

type IgnoreGeneratedFileTestSuite struct {
	VerifierTestSuite
}

func (s *IgnoreGeneratedFileTestSuite) SetupSuite() {
	s.options.Scheme = ImportGroupVerificationSchemeStdLocalThirdParty
	s.options.LocalPrefix = "github.com/pavius/impi"
	s.options.IgnoreGenerated = true
}

func (s *IgnoreGeneratedFileTestSuite) TestValidAllGroups() {
	verificationTestCases := []verificationTestCase{
		{
			name: "invalid order, but ignore generated files",
			contents: `// Code generated by foo; DO NOT EDIT.
// github.com/example/foo

package fixtures

import (
    "fmt"
    "os"
    "github.com/example/foo"
    "path"
)
`,
		},
	}
	s.verifyTestCases(verificationTestCases)
}

func TestIgnoreGeneratedFileTestSuite(t *testing.T) {
	suite.Run(t, new(IgnoreGeneratedFileTestSuite))
}

type NotIgnoreGeneratedFileTestSuite struct {
	VerifierTestSuite
}

func (s *NotIgnoreGeneratedFileTestSuite) SetupSuite() {
	s.options.Scheme = ImportGroupVerificationSchemeStdLocalThirdParty
	s.options.LocalPrefix = "github.com/pavius/impi"
}

func (s *NotIgnoreGeneratedFileTestSuite) TestValidAllGroups() {

	verificationTestCases := []verificationTestCase{
		{
			name: "invalid order, not ignoring generated files",
			contents: `// Code generated by foo; DO NOT EDIT.
// github.com/example/foo

package fixtures

import (
    "fmt"
    "os"
    "github.com/example/foo"
    "path"
)
`,
			expectedErrorStrings: []string{
				"Imports of different types are not allowed in the same group",
			},
		},
	}
	s.verifyTestCases(verificationTestCases)
}

func TestNotIgnoreGeneratedFileTestSuite(t *testing.T) {
	suite.Run(t, new(NotIgnoreGeneratedFileTestSuite))
}
