package compatibility

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type ReportData struct {
	NodeName      string
	KernelVersion string
	IsSupported   bool
}

func isSupportedKernelVersion(kernelVersion string) bool {
	// TODO: retrieve this from some external system?
	supportedKernelVersions := []string{
		"5.10.198-187.748.amzn2.x86_64",
		"5.10.199-190.747.amzn2.x86_64",
		"5.14.0-362.8.1.el9_3.x86_64",
		"5.15.0-1050-aws",
	}

	for _, supportedKernelVersion := range supportedKernelVersions {
		if supportedKernelVersion == kernelVersion {
			return true
		}
	}

	return false
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
