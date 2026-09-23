package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type entry struct {
	name   string
	data   []byte
	digest string
}

func validName(name string) bool {
	if name == "" || strings.HasPrefix(name, "/") {
		return false
	}
	for _, part := range strings.Split(name, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
		for _, ch := range part {
			if !(ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '_' || ch == '+' || ch == '-') {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 4 {
		panic("usage: tzdata-generator <zoneinfo-directory> <expected-version> <tzdata-directory>")
	}
	source := os.Args[1]
	expected := os.Args[2]
	destination := os.Args[3]
	index, err := os.ReadFile(filepath.Join(source, "tzdata.zi"))
	if err != nil {
		panic(err)
	}
	first, _, _ := strings.Cut(string(index), "\n")
	if first != "# version "+expected {
		panic("source timezone version differs from requested version")
	}
	root, err := filepath.EvalSymlinks(source)
	if err != nil {
		panic(err)
	}
	entries := make([]entry, 0)
	err = filepath.WalkDir(source, func(path string, item fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if item.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(relative)
		if !validName(name) || name == "localtime" || name == "posixrules" || strings.HasPrefix(name, "posix/") || strings.HasPrefix(name, "right/") {
			return nil
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		inside, err := filepath.Rel(root, resolved)
		if err != nil || inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
			return fmt.Errorf("timezone entry escapes source root: %s", name)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.HasPrefix(data, []byte("TZif")) {
			return nil
		}
		if len(data) > 16777216 {
			return fmt.Errorf("timezone entry exceeds parser limit: %s", name)
		}
		entries = append(entries, entry{name: name, data: data, digest: fmt.Sprintf("%x", sha256.Sum256(data))})
		return nil
	})
	if err != nil {
		panic(err)
	}
	byName := make(map[string]entry, len(entries))
	for _, item := range entries {
		byName[item.name] = item
	}
	for _, line := range strings.Split(string(index), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 || fields[0] != "L" {
			continue
		}
		target, alias := fields[1], fields[2]
		if !validName(target) || !validName(alias) {
			panic(fmt.Errorf("timezone alias has invalid name: %s", line))
		}
		original, exists := byName[target]
		if !exists {
			panic(fmt.Errorf("timezone alias target is missing: %s", target))
		}
		if existing, exists := byName[alias]; exists {
			if existing.digest != original.digest {
				panic(fmt.Errorf("timezone alias disagrees with target: %s", alias))
			}
			continue
		}
		item := entry{name: alias, data: original.data, digest: original.digest}
		entries = append(entries, item)
		byName[alias] = item
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
	if len(entries) < 400 {
		panic("timezone source has too few TZif entries")
	}
	var checksums strings.Builder
	for _, item := range entries {
		fmt.Fprintf(&checksums, "%s  %s\n", item.digest, item.name)
	}
	var output strings.Builder
	fmt.Fprintln(&output, "package tzdata;")
	fmt.Fprintln(&output)
	fmt.Fprintf(&output, "pub const VERSION: string = %s;\n", strconv.Quote(expected))
	fmt.Fprintf(&output, "pub const ENTRY_COUNT: isize = %d;\n", len(entries))
	fmt.Fprintf(&output, "pub const MANIFEST_SHA256: string = %q;\n", fmt.Sprintf("%x", sha256.Sum256([]byte(checksums.String()))))
	fmt.Fprintf(&output, "const ENTRIES: [(string, string, string); %d] = [\n", len(entries))
	for _, item := range entries {
		fmt.Fprintf(&output, "    (%s, %s, %s),\n", strconv.Quote(item.name), strconv.Quote(base64.StdEncoding.EncodeToString(item.data)), strconv.Quote(item.digest))
	}
	fmt.Fprintln(&output, "];")
	if err := os.MkdirAll(filepath.Join(destination, "data"), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "data.gom"), []byte(output.String()), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "data", "VERSION"), []byte(expected+"\n"), 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "data", "SHA256SUMS"), []byte(checksums.String()), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("generated %d timezone entries from %s\n", len(entries), expected)
}
