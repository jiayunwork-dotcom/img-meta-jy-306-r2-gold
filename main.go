package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"img-meta/internal/meta"
	"img-meta/internal/rename"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  img-meta scan --dir DIR [--template T] [--apply]")
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 || args[0] != "scan" {
		usage(stderr)
		return fmt.Errorf("missing command")
	}
	var dir, template string
	apply := false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--dir":
			if i+1 < len(args) {
				dir = args[i+1]
				i++
			}
		case "--template":
			if i+1 < len(args) {
				template = args[i+1]
				i++
			}
		case "--apply":
			apply = true
		}
	}
	if dir == "" {
		usage(stderr)
		return fmt.Errorf("scan requires --dir")
	}
	paths, err := imageFiles(dir)
	if err != nil {
		return err
	}
	var metas []meta.Meta
	for _, p := range paths {
		m, err := meta.Extract(p)
		if err != nil {
			fmt.Fprintf(stderr, "skip %s: %v\n", p, err)
			continue
		}
		metas = append(metas, m)
	}
	if template == "" {
		template = "{aspect}-{w}x{h}-{i}"
	}
	plans := rename.Plan(metas, template)
	for _, p := range plans {
		fmt.Fprintf(stdout, "%s -> %s\n", p.From, p.To)
	}
	if apply {
		if err := rename.Apply(plans); err != nil {
			return err
		}
	}
	return nil
}

func imageFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out, nil
}
