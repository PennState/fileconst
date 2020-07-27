package fileconst // nolint:testpackage

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/PennState/proctor/pkg/goldenfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectory(t *testing.T) {
	exp := os.Getenv("PWD")
	require.NotEmpty(t, exp)

	os.Setenv("GOFILE", "./const_test.go")
	defer os.Unsetenv("GOFILE")

	act, err := directory()
	require.NoError(t, err)
	assert.Equal(t, exp, act)
}

func TestFiles(t *testing.T) {
	tests := []struct {
		name string
		inp  []string
		len  int
	}{
		{
			name: "Single directory",
			inp:  []string{"./testdata/"},
			len:  3,
		},
		{
			name: "Two directories",
			inp:  []string{"./", "./testdata/"},
			len:  10,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			fis, err := files(test.inp)
			assert.NoError(t, err)
			assert.Len(t, fis, test.len)
		})
	}
}

func TestGenerate(t *testing.T) {
	spec := Spec{
		Package: "something",
		Files: []FileSpec{
			{
				Name:    "MergeNodeQuery",
				Comment: "// MergeNodeQuery does nothing",
				Content: "`MERGE (n:Node {name: 'name'))`",
			},
		},
	}

	sb := strings.Builder{}
	writer := NOPWriteCloser(&sb)
	err := generate(&spec, writer)
	assert.NoError(t, err)
	goldenfile.AssertStringEq(t, "testdata/cypher_gen.go.golden", sb.String())
}

func TestProcess(t *testing.T) {
	exts := map[string]bool{
		"cypher": true,
		"sql":    true,
	}

	paths, err := files([]string{"./testdata/"})
	require.NoError(t, err)
	require.Len(t, paths, 3)

	spec, err := process("something", paths, exts)
	assert.NoError(t, err)
	assert.Len(t, spec.Files, 2)
}

type nopWriteCloser struct {
	io.Writer
}

func NOPWriteCloser(writer io.Writer) io.WriteCloser {
	return nopWriteCloser{Writer: writer}
}

func (w nopWriteCloser) Close() error {
	return nil
}
