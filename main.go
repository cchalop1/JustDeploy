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

// Fonction for build the web part of JustDeploy
func (m *JustDeploy) BuildWeb(buildSrc *Directory) *Directory {
	return dag.Container().
		From("node:lts").
		WithMountedDirectory("/mnt", buildSrc.WithoutDirectory("node_modules")).
		WithWorkdir("/mnt").
		// WithExec([]string{"rm", "-rf", "node_modules"}).
		WithExec([]string{"npm", "i", "-g", "pnpm"}).
		WithExec([]string{"pnpm", "i"}).
		WithExec([]string{"pnpm", "run", "build"}).
		Directory("./dist")
}

// Function to build the binary from my golang code (distDir) is a build of my web part
func (m *JustDeploy) BuildGoApp(distDir *Directory, source *Directory) *Container {
	arch := runtime.GOARCH
	os := runtime.GOOS
	return dag.Container().
		From("golang:latest").
		WithWorkdir("/mnt").
		WithMountedDirectory("/mnt", source).
		WithMountedDirectory("/mnt/internal/web/dist", distDir).
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-121")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-121")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec([]string{"go", "build", "-o", "bin/justdeploy", "./cmd/just-deploy/main.go"})
}

// This function build my application
func (m *JustDeploy) Build(ctx context.Context, source *Directory) *Directory {
	fmt.Println("CI start !!!")
	distWeb := m.BuildWeb(source.Directory("web"))
	binContainer := m.BuildGoApp(distWeb, source)
	return binContainer.Directory("./bin")
	// return binContainer.WithEntrypoint([]string{"./bin/justdeploy"}).WithExposedPort(8080)
	// return binContainer
}

func (m *JustDeploy) Serve(ctx context.Context, source *Directory) *Container {
	bin := m.Build(ctx, source)
	return dag.Container().
		From("golang:alpine").
		WithMountedDirectory("/app", bin).
		WithWorkdir("/app").
		WithEntrypoint([]string{"/app/justdeploy"}).
		WithExposedPort(8080)

}
