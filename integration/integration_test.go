package integration

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	bpDir, dotnetCoreConfURI string
)

var suite = spec.New("Integration", spec.Report(report.Terminal{}))

func init() {
	suite("Integration", testIntegration)
}

func TestIntegration(t *testing.T) {
	var err error
	Expect := NewWithT(t).Expect
	bpDir, err = dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())

	dotnetCoreConfURI, err = dagger.PackageBuildpack(bpDir)
	Expect(err).ToNot(HaveOccurred())
	defer dagger.DeleteBuildpack(dotnetCoreConfURI)

	suite.Run(t)
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var Expect func(interface{}, ...interface{}) Assertion
	var Eventually func(interface{}, ...interface{}) AsyncAssertion
	it.Before(func() {
		Expect = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
	})

	when("the app is self contained", func() {
		it("builds successfully", func() {
			appRoot := filepath.Join("testdata", "self_contained_2.1")

			app, err := dagger.PackBuild(appRoot, dotnetCoreConfURI)
			Expect(err).NotTo(HaveOccurred())
			defer app.Destroy()

			Expect(app.Start()).To(Succeed())

			Eventually(func() string {
				body, _, _ := app.HTTPGet("/")
				return body
			}).Should(ContainSubstring("Hello World"))
		})
	})
}
