package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/axw/gocov"
)

func unmarshalJSON(data []byte) (packages []*gocov.Package, err error) {
	result := &struct{ Packages []*gocov.Package }{}
	err = json.Unmarshal(data, result)
	if err == nil {
		packages = result.Packages
	}
	return
}

func packageCoverage(pkg *gocov.Package) float64 {
	var totalStatements, totalReached int
	for _, fn := range pkg.Functions {
		for _, stmt := range fn.Statements {
			totalStatements++
			if stmt.Reached > 0 {
				totalReached++
			}
		}
	}
	var stmtPercent float64
	if totalStatements > 0 {
		stmtPercent = float64(totalReached) / float64(totalStatements) * 100
	}
	return stmtPercent
}

func main() {
	flag.Parse()
	files := make([]*os.File, 0, 1)
	if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			file, err := os.Open(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file (%s): %s\n", name, err)
			} else {
				files = append(files, file)
			}
		}
	} else {
		files = append(files, os.Stdin)
	}
	var coverageSum float64
	var coverageNum uint64
	for _, file := range files {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read coverage file: %s\n", err)
			return
		}
		packages, err := unmarshalJSON(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal coverage data: %s\n", err)
			return
		}
		for _, pkg := range packages {
			pkgCoverage := packageCoverage(pkg)
			coverageSum += pkgCoverage
			coverageNum++
		}
		if file != os.Stdin {
			err := file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to close file: %s\n", err)
			}
		}
	}
	var coverage float64
	if coverageNum > 0 {
		coverage = coverageSum / float64(coverageNum)
	}
	fmt.Printf("gocov combined coverage: %.2f%%", coverage)
	return
}
