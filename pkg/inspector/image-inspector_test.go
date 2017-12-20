package inspector

import (
	"context"
	"fmt"

	"testing"

	docker "github.com/fsouza/go-dockerclient"
	iiapi "github.com/openshift/image-inspector/pkg/api"
)

type FailMockScanner struct{}
type SuccMockScanner struct {
	FailMockScanner
}
type NoResMockScanner struct {
	SuccMockScanner
}
type SuccWithReportMockScanner struct {
	SuccMockScanner
}

func (ms *FailMockScanner) Scan(context.Context, string, *docker.Image, iiapi.FilesFilter) ([]iiapi.Result, interface{}, error) {
	return nil, nil, fmt.Errorf("FAIL SCANNER!")
}
func (ms *FailMockScanner) Name() string {
	return "MockScanner"
}
func (ms *SuccMockScanner) Scan(context.Context, string, *docker.Image, iiapi.FilesFilter) ([]iiapi.Result, interface{}, error) {
	return []iiapi.Result{}, nil, nil
}

func TestScanImage(t *testing.T) {
	ctx := context.Background()
	for k, v := range map[string]struct {
		ii         defaultImageInspector
		s          iiapi.Scanner
		shouldFail bool
	}{
		"Scanner fails on scan": {ii: defaultImageInspector{}, s: &FailMockScanner{}, shouldFail: true},
		"Happy Flow":            {ii: defaultImageInspector{}, s: &SuccMockScanner{}, shouldFail: false},
	} {
		v.ii.opts.DstPath = "here"
		_, _, err := v.s.Scan(ctx, v.ii.opts.DstPath, nil, nil)
		if v.shouldFail && err == nil {
			t.Errorf("%s should have failed but it didn't!", k)
		}
		if !v.shouldFail {
			if err != nil {
				t.Errorf("%s should have succeeded but failed with %v", k, err)
			}
		}
	}
}
