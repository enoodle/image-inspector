package inspector_test

import (
	. "github.com/openshift/image-inspector/pkg/inspector"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/fsouza/go-dockerclient"
	iiapi "github.com/openshift/image-inspector/pkg/api"
	iicmd "github.com/openshift/image-inspector/pkg/cmd"
	openscap "github.com/openshift/image-inspector/pkg/openscap"
)

type happyEmptyImageAcquirer struct{}

func (heia *happyEmptyImageAcquirer) Acquire(source string) (string, docker.Image, iiapi.ScanResult, iiapi.FilesFilter, error) {
	scanResults := iiapi.ScanResult{
		APIVersion: iiapi.DefaultResultsAPIVersion,
		ImageName:  source,
		Results:    []iiapi.Result{},
	}
	return "", docker.Image{}, scanResults, nil, nil
}

type happyEmptyImageServer struct{}

func (heis *happyEmptyImageServer) ServeImage(meta *iiapi.InspectorMetadata,
	ImageServeURL string,
	results iiapi.ScanResult,
	scanReport []byte,
	htmlScanReport []byte) error {
	return nil
}

type happyEmptyOscapScanner struct{}

func (hes *happyEmptyOscapScanner) Scan(ctx context.Context, path string, image *docker.Image, filter iiapi.FilesFilter) ([]iiapi.Result, interface{}, error) {
	return []iiapi.Result{}, openscap.OpenSCAPReport{}, nil
}

func (hes *happyEmptyOscapScanner) Name() string {
	return "happyEmptyOscapScanner"
}

type happyEmptyScannerFactory struct{}

func (hesf *happyEmptyScannerFactory) CreateScanner(string) (iiapi.Scanner, error) {
	return &happyEmptyOscapScanner{}, nil
}

var _ = Describe("ImageInspector", func() {
	var (
		ii         ImageInspector
		opts       *iicmd.ImageInspectorOptions
		serve      = "localhost:8088"
		validToken = "w599voG89897rGVDmdp12WA681r9E5948c1CJTPi8g4HGc4NWaz62k6k1K0FMxHW40H8yOO3Hoe"
		err        error
	)
	BeforeEach(func() {
		opts = iicmd.NewDefaultImageInspectorOptions()
		opts.Serve = serve
		opts.AuthToken = validToken
		opts.Image = "registry.access.redhat.com/rhel7:latest"
		opts.ScanType = "openscap"
		opts.DstPath = ""
	})
	Describe("Inspect()", func() {
		It("Simple Sanity with empty implementations", func() {
			opts.ImageAcquirer = &happyEmptyImageAcquirer{}
			opts.ImageServer = &happyEmptyImageServer{}
			opts.ScannerFactory = &happyEmptyScannerFactory{}
			ii = NewDefaultImageInspector(*opts)
			err = ii.Inspect()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
