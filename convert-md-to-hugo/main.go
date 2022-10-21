package main

import (
	"flag"
	"fmt"
	yaml3 "gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var docsDir = flag.String("docs-dir", "", "Path to documentation")

type FrontMatter struct {
	Title string `yaml:"title"`
}

func main() {
	flag.Parse()

	err := filepath.WalkDir(*docsDir, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		err = processFile(path)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func processFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
outer:
	for i, l := range lines {
		if l == "<!-- This comment is uncommented when auto-synced to www-kluctl.io" {
			for j := i + 1; j < len(lines); j++ {
				if lines[j] == "-->" {
					lines[i] = ""
					lines[j] = ""
					break outer
				}
			}
		}
	}

	var frontMatter string
outer2:
	for i, l := range lines {
		if l == "---" {
			for j := i + 1; j < len(lines); j++ {
				if lines[j] == "---" {
					frontMatter = strings.Join(lines[i:j], "\n")
					break outer2
				}
			}
		}
	}
	if frontMatter != "" {
		d := yaml3.NewDecoder(strings.NewReader(frontMatter))

		var fm FrontMatter
		err = d.Decode(&fm)
		if err != nil {
			return err
		}

		// remove unnecessary "# title"
		for i, l := range lines {
			if l == fmt.Sprintf("# %s", fm.Title) {
				lines[i] = ""
				break
			}
		}
	}

	// remove "# Table of Contents"
outer3:
	for i, l := range lines {
		if strings.HasPrefix(l, "#") && strings.Index(strings.ToLower(l), "table of contents") != -1 {
			for j := i + 1; j < len(lines); j++ {
				firstChar := ""
				if strings.TrimSpace(lines[j]) != "" {
					firstChar = string(lines[j][0])
				}
				if (firstChar != "" && !unicode.IsNumber(rune(lines[j][0]))) || j == len(lines)-1 {
					lines = append(lines[:i], lines[j:]...)
					break outer3
				}
			}
		}
	}

	b = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(path, b, 0o600)
	if err != nil {
		return err
	}

	return nil
}
