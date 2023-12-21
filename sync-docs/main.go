package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v57/github"
	yaml3 "gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

var argRepo = flag.String("repo", "", "")
var argMinVersion = flag.String("min-version", "", "")
var argSubdir = flag.String("subdir", "", "")
var argDest = flag.String("dest", "", "")
var argDestSubdir = flag.String("dest-subdir", "", "")
var argWithRootReadme = flag.Bool("with-root-readme", false, "")

func main() {
	flag.Parse()

	err := doMain(context.Background())
	if err != nil {
		panic(err)
	}
}

func getReleaseTags(ctx context.Context, owner string, repo string, minVersion string) (semver.Collection, error) {
	client := github.NewClient(nil)

	opt := github.ListOptions{
		PerPage: 100,
	}

	minVersion2, err := semver.NewVersion(minVersion)
	if err != nil {
		return nil, err
	}

	var versions semver.Collection
	for {
		r, resp, err := client.Repositories.ListReleases(ctx, owner, repo, &opt)
		if err != nil {
			return nil, err
		}
		for _, x := range r {
			v, err := semver.NewVersion(*x.TagName)
			if err != nil {
				continue
			}
			if v.LessThan(minVersion2) {
				continue
			}
			versions = append(versions, v)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	sort.Sort(versions)

	return versions, nil
}

func doMain(ctx context.Context) error {
	s := strings.Split(*argRepo, "/")
	if len(s) != 2 {
		return fmt.Errorf("unexpected repo %s", *argRepo)
	}

	versions, err := getReleaseTags(ctx, s[0], s[1], *argMinVersion)
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	repoDir := filepath.Join(tmpDir, "repo")

	cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s.git", *argRepo), repoDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "config", "advice.detachedHead", "false")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir
	err = cmd.Run()
	if err != nil {
		return err
	}

	for _, version := range versions {
		err = processRef(ctx, repoDir, *argSubdir, "v"+version.String(), filepath.Join(*argDest, "v"+version.String(), *argDestSubdir), *argWithRootReadme)
		if err != nil {
			return err
		}
	}

	err = processRef(ctx, repoDir, *argSubdir, "v"+versions[len(versions)-1].String(), filepath.Join(*argDest, "latest", *argDestSubdir), *argWithRootReadme)
	if err != nil {
		return err
	}
	err = processRef(ctx, repoDir, *argSubdir, "main", filepath.Join(*argDest, "devel", *argDestSubdir), *argWithRootReadme)
	if err != nil {
		return err
	}
	return nil
}

func processRef(ctx context.Context, repoDir string, subdir string, ref string, dest string, withRootReadme bool) error {
	cmd := exec.Command("git", "checkout", ref, "--")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		return err
	}

	repoSubDir := filepath.Join(repoDir, subdir)
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

		err = processFile(path, filepath.Join(dest, relDestPath), repoDir, relGitPath)
		if err != nil {
			return fmt.Errorf("failed to process %s, %w", path, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if withRootReadme {
		b, err := os.ReadFile(filepath.Join(repoDir, "README.md"))
		if err != nil {
			return err
		}
		b = bytes.ReplaceAll(b, []byte(fmt.Sprintf("./%s/", subdir)), []byte("./"))
		err = os.WriteFile(filepath.Join(repoDir, "README.md"), b, 0600)
		if err != nil {
			return err
		}

		err = processFile(filepath.Join(repoDir, "README.md"), filepath.Join(dest, "_index.md"), repoDir, "README.md")
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
	frontMatter["github_repo"] = fmt.Sprintf("https://github.com/%s", *argRepo)
	frontMatter["path_base_for_github_subdir"] = map[string]any{
		"from": ".*",
		"to":   fmt.Sprintf("main/%s", relGitPath),
	}

	lastMod, err := getGitLastMod(repoDir, relGitPath)
	if err != nil {
		return err
	}
	frontMatter["lastmod"] = lastMod.Format(time.RFC3339)

	lines = append([]string{
		"<!-- WARNING WARNING WARNING -->",
		fmt.Sprintf("<!-- DO NOT EDIT THIS FILE, IT IS AUTO SYNCED FROM github.com/%s -->", *argRepo),
		"<!-- WARNING WARNING WARNING -->",
	}, lines...)

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

	r := regexp.MustCompile(`]\((\..*)/README\.md(#.*)?\)`)
	for i, l := range lines {
		lines[i] = r.ReplaceAllString(l, `]($1/$2)`)
	}

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
		return time.Time{}, nil
	}
	return lastMod, nil
}
