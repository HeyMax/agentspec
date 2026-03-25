package check

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ScanResult struct {
	Routes  []RouteEntry
	Structs []StructEntry
}

type RouteEntry struct {
	Method string
	Path   string
	File   string
	Line   int
}

type StructEntry struct {
	Name   string
	Fields []string
	File   string
}

var (
	goRoutePatterns = []*regexp.Regexp{
		regexp.MustCompile(`\.\s*(GET|POST|PUT|PATCH|DELETE)\s*\(\s*"([^"]+)"`),
		regexp.MustCompile(`HandleFunc\s*\(\s*"(GET|POST|PUT|PATCH|DELETE)\s+([^"]+)"`),
	}
	goStructPattern = regexp.MustCompile(`type\s+(\w+)\s+struct\s*\{`)
	goFieldPattern  = regexp.MustCompile(`(\w+)\s+\w+.*` + "`" + `.*json:"(\w+)"`)
	tsFetchPattern  = regexp.MustCompile(`fetch\s*\(\s*['` + "`" + `]([^'` + "`" + `]+)['` + "`" + `].*method:\s*['"](\w+)['"]`)
)

func ScanDir(dir string) (*ScanResult, error) {
	result := &ScanResult{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if base == "node_modules" || base == ".git" || base == "vendor" || base == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(path)
		switch ext {
		case ".go":
			scanGoFile(path, result)
		case ".ts", ".tsx":
			scanTSFile(path, result)
		}
		return nil
	})
	return result, err
}

func scanGoFile(path string, result *ScanResult) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")

	for lineNum, line := range lines {
		for _, pat := range goRoutePatterns {
			matches := pat.FindStringSubmatch(line)
			if len(matches) >= 3 {
				result.Routes = append(result.Routes, RouteEntry{
					Method: matches[1], Path: matches[2], File: path, Line: lineNum + 1,
				})
			}
		}
	}

	var currentStruct *StructEntry
	for _, line := range lines {
		if matches := goStructPattern.FindStringSubmatch(line); len(matches) >= 2 {
			currentStruct = &StructEntry{Name: matches[1], File: path}
			continue
		}
		if currentStruct != nil {
			if strings.TrimSpace(line) == "}" {
				result.Structs = append(result.Structs, *currentStruct)
				currentStruct = nil
				continue
			}
			if matches := goFieldPattern.FindStringSubmatch(line); len(matches) >= 3 {
				currentStruct.Fields = append(currentStruct.Fields, matches[2])
			}
		}
	}
}

func scanTSFile(path string, result *ScanResult) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for lineNum, line := range lines {
		if matches := tsFetchPattern.FindStringSubmatch(line); len(matches) >= 3 {
			result.Routes = append(result.Routes, RouteEntry{
				Method: strings.ToUpper(matches[2]), Path: matches[1], File: path, Line: lineNum + 1,
			})
		}
	}
}
