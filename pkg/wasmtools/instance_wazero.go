// NOTE: This file was downloaded from
// https://github.com/bytecodealliance/wasm-tools/releases/tag/v1.235.0

package wasmtools

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm-tools.wasm.gz
var compressed []byte

var decompress = sync.OnceValues(func() ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
})

var buildInfo = sync.OnceValue(func() *debug.BuildInfo {
	build, _ := debug.ReadBuildInfo()
	return build
})

var versionString = sync.OnceValue(func() string {
	const nonePath string = "(none)"
	build := buildInfo()
	if build == nil {
		return nonePath
	}
	version := build.Main.Version
	var revision string
	for _, s := range build.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		}
	}
	if version == "" {
		version = nonePath
	}
	versionString := version
	if revision != "" {
		versionString += " (" + revision + ")"
	}
	return versionString
})

var compilationCache = sync.OnceValue(func() wazero.CompilationCache {
	// First try on-disk cache, so subsequent invocations can share cache
	tmp := os.TempDir()
	if tmp != "" {
		var path string
		build := buildInfo()
		if build == nil {
			path = "(none)"
		} else {
			path = build.Main.Path
		}
		rep := strings.NewReplacer(" ", "", "(", "", ")", "")
		dir := filepath.Join(tmp, rep.Replace(path+"@"+versionString()))
		c, err := wazero.NewCompilationCacheWithDir(dir)
		if err == nil {
			return c
		}
	}

	// Fall back to in-memory cache
	return wazero.NewCompilationCache()
})

// Instance is a compiled wazero instance.
type Instance struct {
	runtime wazero.Runtime
	module  wazero.CompiledModule
}

// New creates a new wazero instance.
func New(ctx context.Context) (*Instance, error) {
	c := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true).
		WithCompilationCache(compilationCache())

	r := wazero.NewRuntimeWithConfig(ctx, c)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		return nil, fmt.Errorf("error instantiating WASI: %w", err)
	}

	wasmTools, err := decompress()
	if err != nil {
		return nil, err
	}

	module, err := r.CompileModule(ctx, wasmTools)
	if err != nil {
		return nil, fmt.Errorf("error compiling wasm module: %w", err)
	}
	return &Instance{runtime: r, module: module}, nil
}

// Close closes the wazero runtime resource.
func (w *Instance) Close(ctx context.Context) error {
	return w.runtime.Close(ctx)
}

// Run runs the wasm module with the context, arguments,
// and optional stdin, stdout, stderr, and filesystem map.
// Supply a context with a timeout or other cancellation mechanism to control execution time.
// Returns an error if instantiation fails.
func (w *Instance) Run(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, fsMap map[string]fs.FS, args ...string) error {
	config := wazero.NewModuleConfig().
		WithRandSource(rand.Reader).
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithArgs(append([]string{"wasm-tools.wasm"}, args...)...)

	if stdin != nil {
		config = config.WithStdin(stdin)
	}
	if stdout != nil {
		config = config.WithStdout(stdout)
	}
	if stderr != nil {
		config = config.WithStderr(stderr)
	}

	fsConfig := wazero.NewFSConfig()
	for guestPath, guestFS := range fsMap {
		fsConfig = fsConfig.WithFSMount(guestFS, guestPath)
	}
	config = config.WithFSConfig(fsConfig)

	_, err := w.runtime.InstantiateModule(ctx, w.module, config)
	return err
}
