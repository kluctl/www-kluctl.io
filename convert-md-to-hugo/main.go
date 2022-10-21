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

	var frontMatterStr string
outer2:
	for i, l := range lines {
		if l == "---" {
			for j := i + 1; j < len(lines); j++ {
				if lines[j] == "---" {
					frontMatterStr = strings.Join(lines[i+1:j], "\n")
					lines = lines[j+1:]
					break outer2
				}
			}
		}
	}
	if frontMatterStr == "" {
		return fmt.Errorf("front-matter not found: %s", path)
	}

	var frontMatter map[string]any
	err = yaml3.Unmarshal([]byte(frontMatterStr), &frontMatter)
	if err != nil {
		return err
	}
	title := frontMatter["title"]

	// add github links
	frontMatter["github_repo"] = "https://github.com/kluctl/kluctl"
	if strings.HasSuffix(path, "_index.md") {
		frontMatter["path_base_for_github_subdir"] = map[string]any{
			"from": "content/en/docs/(.*/?)_index.md",
			"to":   "docs/${1}README.md",
		}
	} else {
		frontMatter["path_base_for_github_subdir"] = map[string]any{
			"from": "content/en/docs/(.*)",
			"to":   "docs/$1",
		}
	}

	// remove unnecessary "# title"
	for i, l := range lines {
		if l == fmt.Sprintf("# %s", title) {
			lines[i] = ""
			break
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

	frontMatterBytes, err := yaml3.Marshal(&frontMatter)
	if err != nil {
		return err
	}
	frontMatterStr = fmt.Sprintf("---\n%s---\n", string(frontMatterBytes))
	lines = append(strings.Split(frontMatterStr, "\n"), lines...)

	b = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(path, b, 0o600)
	if err != nil {
		return err
	}

	return nil
}
