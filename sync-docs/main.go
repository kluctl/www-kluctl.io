package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Masterminds/semver/v3"
	yaml3 "gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
)

var repo = flag.String("repo", "", "")
var subdir = flag.String("subdir", "", "")
var dest = flag.String("dest", "", "")
var withRootReadme = flag.Bool("with-root-readme", false, "")
var ref = flag.String("ref", "", "")

func main() {
	flag.Parse()

	err := doMain()
	if err != nil {
		panic(err)
	}
}

func getLatestTag(repo string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/tags", repo))
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("get returned %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tags []map[string]any
	err = json.Unmarshal(respBytes, &tags)
	if err != nil {
		return "", err
	}

	var versions semver.Collection
	for _, tag := range tags {
		name, ok := tag["name"]
		if !ok {
			continue
		}
		nameStr, ok := name.(string)
		if !ok {
			continue
		}
		v, err := semver.NewVersion(nameStr)
		if err != nil {
			return "", err
		}
		versions = append(versions, v)
	}
	sort.Sort(versions)

	return versions[len(versions)-1].Original(), nil
}

func doMain() error {
	version, err := getLatestTag(*repo)
	if err != nil {
		return err
	}

	if *ref != "" {
		version = *ref
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "repo")

	cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s.git", *repo), repoDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "checkout", version, "--")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		return err
	}

	repoSubDir := filepath.Join(repoDir, *subdir)
	err = filepath.WalkDir(repoSubDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		relGitPath, err := filepath.Rel(repoDir, path)
		if err != nil {
			return fmt.Errorf("failed to process %s, %w", path, err)
		}
		relDestPath, err := filepath.Rel(repoSubDir, path)
		if err != nil {
			return fmt.Errorf("failed to process %s, %w", path, err)
		}

		if filepath.Base(relDestPath) == "README.md" {
			relDestPath = filepath.Join(filepath.Dir(relDestPath), "_index.md")
		}

		err = processFile(path, filepath.Join(*dest, relDestPath), repoDir, relGitPath)
		if err != nil {
			return fmt.Errorf("failed to process %s, %w", path, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if *withRootReadme {
		b, err := os.ReadFile(filepath.Join(repoDir, "README.md"))
		if err != nil {
			return err
		}
		b = bytes.ReplaceAll(b, []byte(fmt.Sprintf("./%s/", *subdir)), []byte("./"))
		err = os.WriteFile(filepath.Join(repoDir, "README.md"), b, 0600)
		if err != nil {
			return err
		}

		err = processFile(filepath.Join(repoDir, "README.md"), filepath.Join(*dest, "_index.md"), repoDir, "README.md")
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(sourcePath string, destPath string, repoDir string, relGitPath string) error {
	b, err := os.ReadFile(sourcePath)
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
		frontMatterFile := strings.TrimSuffix(sourcePath, ".md") + ".front-matter.yaml"
		if _, err := os.Stat(frontMatterFile); err == nil {
			b, err := os.ReadFile(frontMatterFile)
			if err != nil {
				return fmt.Errorf("failed to read front matter file %s: %w", frontMatterFile, err)
			}
			frontMatterStr = string(b)
		} else {
			return fmt.Errorf("front-matter not found: %s", sourcePath)
		}
	}

	var frontMatter map[string]any
	err = yaml3.Unmarshal([]byte(frontMatterStr), &frontMatter)
	if err != nil {
		os.Stderr.WriteString(frontMatterStr + "\n")
		return err
	}
	title := frontMatter["title"]

	// add github links
	frontMatter["github_repo"] = fmt.Sprintf("https://github.com/%s", *repo)
	frontMatter["path_base_for_github_subdir"] = map[string]any{
		"from": ".*",
		"to":   fmt.Sprintf("main/%s", relGitPath),
	}

	lastMod, err := getGitLastMod(repoDir, relGitPath)
	if err != nil {
		return err
	}
	frontMatter["lastmod"] = lastMod.Format(time.RFC3339)

	// remove unnecessary "# title"
	for i, l := range lines {
		if l == fmt.Sprintf("# %s", title) {
			lines[i] = ""
			break
		}
		if l == fmt.Sprintf("<h1>%s</h1>", title) {
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

	_ = os.MkdirAll(filepath.Dir(destPath), 0o755)

	b = []byte(strings.Join(lines, "\n"))
	err = os.WriteFile(destPath, b, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func getGitLastMod(repoDir string, relGitPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%ci", relGitPath)
	cmd.Dir = repoDir
	stdoutBuf := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return time.Time{}, err
	}
	lastMod, err := time.Parse("2006-01-02 15:04:05 -0700", stdoutBuf.String())
	if err != nil {
		return time.Time{}, err
	}
	return lastMod, nil
}
