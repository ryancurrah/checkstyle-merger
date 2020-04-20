package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
)

// CheckstyleReport represents a report in the checkstyle format
type CheckstyleReport struct {
	XMLName xml.Name              `xml:"checkstyle"`
	Version string                `xml:"version,attr"`
	File    CheckstyleReportFiles `xml:"file"`
}

// fileExists checks if error file exists, if so returns the index and true
func (r *CheckstyleReport) fileExists(filename string) (int, bool) {
	for n := range r.File {
		if r.File[n].Name == filename {
			return n, true
		}
	}

	return 0, false
}

// AddFiles to checkstyle report ensuring no duplicate file names
func (r *CheckstyleReport) AddFiles(files CheckstyleReportFiles) {
	for n := range files {
		index, exists := r.fileExists(files[n].Name)
		if !exists {
			r.File = append(r.File, files[n])
			continue
		}

		r.File[index].Error = append(r.File[index].Error, files[n].Error...)
	}
}

// CheckstyleReportFiles is a list of reported files with errors
type CheckstyleReportFiles []struct {
	Name  string                     `xml:"name,attr"`
	Error CheckstyleReportFileErrors `xml:"error"`
}

// CheckstyleReportFileErrors is a list of reported errors
type CheckstyleReportFileErrors []struct {
	Line     string `xml:"line,attr"`
	Column   string `xml:"column,attr"`
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
}

var usage = `Usage: checkstyle-merger [options] [files]
Options:
  -o  Merged report filename`

func main() {
	flag.Usage = func() { fmt.Println(usage) }

	outputFileName := flag.String("o", "", "merged report filename")

	flag.Parse()

	checkstyleFiles := flag.Args()

	if len(checkstyleFiles) == 0 {
		flag.Usage()
		return
	}

	mergedReport := CheckstyleReport{XMLName: xml.Name{Local: "checkstyle"}, Version: "1.0.0"}

	for n := range checkstyleFiles {
		var report CheckstyleReport

		checkstyleFile, err := ioutil.ReadFile(checkstyleFiles[n])
		if err != nil {
			panic(err)
		}

		err = xml.Unmarshal(checkstyleFile, &report)
		if err != nil {
			panic(err)
		}

		mergedReport.AddFiles(report.File)
	}

	mergedOutput, _ := xml.MarshalIndent(&mergedReport, "", "    ")

	mergedReportFile := fmt.Sprintf("%s%s", xml.Header, string(mergedOutput))

	err := ioutil.WriteFile(*outputFileName, []byte(mergedReportFile), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Merged %d reports to %s\n", len(checkstyleFiles), *outputFileName)
}
