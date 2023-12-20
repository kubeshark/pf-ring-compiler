package compatibility

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	SUPPORTED_KERNELS_ENDPOINT = "https://api.kubeshark.co/kernel-modules/meta/versions.json"
)

type ReportData struct {
	NodeName      string
	KernelVersion string
	IsSupported   bool
}

func isSupportedKernelVersion(kernelVersion string) (bool, error) {
	resp, err := http.Get(SUPPORTED_KERNELS_ENDPOINT)
	if err != nil {
		return false, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %w", err)
	}

	var supportedKernelVersions []string
	if err := json.Unmarshal(body, &supportedKernelVersions); err != nil {
		return false, fmt.Errorf("error unmarshaling response: %w", err)
	}

	for _, supportedKernelVersion := range supportedKernelVersions {
		if supportedKernelVersion == kernelVersion {
			return true, nil
		}
	}

	return false, nil
}

func printReportTable(items []ReportData) {
	headerFmt := color.New(color.Bold, color.Underline).SprintfFunc()

	tbl := table.New("Node", "Kernel Version", "Supported")
	tbl.WithHeaderFormatter(headerFmt)

	compatMsg := "Cluster is compatible"
	compatColor := color.New(color.Bold, color.Underline, color.FgGreen)
	for _, item := range items {
		tbl.AddRow(item.NodeName, item.KernelVersion, item.IsSupported)
		if !item.IsSupported {
			compatMsg = "Cluster is not compatible"
			compatColor = color.New(color.Bold, color.Underline, color.FgRed)
		}
	}

	tbl.Print()

	compatColor.Printf("\n%s\n", compatMsg)
}
