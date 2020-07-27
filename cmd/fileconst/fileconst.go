package main

import (
	"flag"
	"os"
	"strings"

	"github.com/PennState/fileconst"
	log "github.com/sirupsen/logrus"
)

type multiString []string

func (a *multiString) Set(v string) error {
	*a = append(*a, v)
	return nil
}

func (a multiString) String() string {
	return strings.Join(a, "|")
}

func main() {
	const (
		shortHand = " (short-hand)"

		dirUsage = "directory containing the files to be processed - " +
			"may be used more than once for multiple directories.  Both " +
			"relative and absolute paths may be used but beware that " +
			"relative paths are relative to the directory with the " +
			"go generator directive."
		extUsage = "extensions of files to process - may be used more " +
			"than once if multiple files types should be processed" +
			" together.  At least one extension must be specified."
		helpUsage = "displays this help screen."
	)

	var help bool

	dirs := multiString{}
	exts := multiString{}

	flag.Var(&dirs, "dir", dirUsage)
	flag.Var(&dirs, "d", dirUsage+shortHand)
	flag.Var(&exts, "ext", extUsage)
	flag.Var(&exts, "e", extUsage+shortHand)
	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage+shortHand)

	flag.Parse()

	if help {
		showHelp()
		os.Exit(0)
	}

	if len(exts) == 0 {
		showHelp()
		log.Fatal("At least one file extension must be specfied")
	}

	log.Info("ext: ", exts)

	err := fileconst.Run(dirs, toSet(exts))
	if err != nil {
		log.WithError(err).Fatal("fileconst failed - see associated err")
	}

	os.Exit(0)
}

func showHelp() {
	flag.PrintDefaults()
}

func toSet(a []string) map[string]bool {
	m := map[string]bool{}

	for _, s := range a {
		m[s] = true
	}

	return m
}
