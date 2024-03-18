// A generated module for JustDeploy functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"runtime"
)

type JustDeploy struct{}

// const (
// 	BIN_DIR,
// 	EXECUTABLE,
// )

func (m *JustDeploy) BuildWeb(buildSrc *Directory) *Directory {
	return dag.Container().
		From("node:lts").
		WithMountedDirectory("/mnt", buildSrc).
		WithWorkdir("/mnt").
		WithExec([]string{"npm", "i", "-g", "pnpm"}).
		WithExec([]string{"pnpm", "i"}).
		WithExec([]string{"pnpm", "run", "build"}).
		Directory("./dist")
}

func (m *JustDeploy) BuildGoApp(distDir *Directory, source *Directory, arch string, os string) *Container {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if os == "" {
		os = runtime.GOOS
	}
	return dag.Container().
		From("golang:latest").
		WithWorkdir("/mnt").
		WithMountedDirectory("/mnt", source).
		WithMountedDirectory("/mnt/internal/web/dist", distDir).
		WithExec([]string{"go", "build", "-o", "bin/justdeploy", "./cmd/just-deploy/main.go"})
}

func (m *JustDeploy) Build(ctx context.Context, source *Directory) {
	fmt.Println("CI start !!!")
	distWeb := m.BuildWeb(source.Directory("web"))
	binContainer := m.BuildGoApp(distWeb, source)
	binContainer.Export(ctx, "./bin")
}
