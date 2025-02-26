package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eagleql/xray-core/common"
	"github.com/eagleql/xray-core/core"
)

var directory = flag.String("pwd", "", "Working directory of Xray vprotogen.")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of vprotogen:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if !filepath.IsAbs(*directory) {
		pwd, wdErr := os.Getwd()
		if wdErr != nil {
			fmt.Println("Can not get current working directory.")
			os.Exit(1)
		}
		*directory = filepath.Join(pwd, *directory)
	}

	pwd := *directory
	GOBIN := common.GetGOBIN()
	binPath := os.Getenv("PATH")
	pathSlice := []string{binPath, GOBIN, pwd}
	binPath = strings.Join(pathSlice, string(os.PathListSeparator))
	os.Setenv("PATH", binPath)

	EXE := ""
	if runtime.GOOS == "windows" {
		EXE = ".exe"
	}
	protoc := "protoc" + EXE

	if path, err := exec.LookPath(protoc); err != nil {
		fmt.Println("Make sure that you have `" + protoc + "` in your system path or current path. To download it, please visit https://github.com/protocolbuffers/protobuf/releases")
		os.Exit(1)
	} else {
		protoc = path
	}

	protoFilesMap := make(map[string][]string)
	walkErr := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".proto") {
			path = path[len(pwd)+1:]
			protoFilesMap[dir] = append(protoFilesMap[dir], path)
		}

		return nil
	})
	if walkErr != nil {
		fmt.Println(walkErr)
		os.Exit(1)
	}

	for _, files := range protoFilesMap {
		for _, relProtoFile := range files {
			var args []string
			if core.ProtoFilesUsingProtocGenGoFast[relProtoFile] {
				args = []string{"--gofast_out", pwd, "--gofast_opt", "paths=source_relative", "--plugin", "protoc-gen-gofast=" + GOBIN + "/protoc-gen-gofast" + EXE}
			} else {
				args = []string{"--go_out", pwd, "--go_opt", "paths=source_relative", "--go-grpc_out", pwd, "--go-grpc_opt", "paths=source_relative", "--plugin", "protoc-gen-go=" + GOBIN + "/protoc-gen-go" + EXE, "--plugin", "protoc-gen-go-grpc=" + GOBIN + "/protoc-gen-go-grpc" + EXE}
			}
			args = append(args, relProtoFile)
			cmd := exec.Command(protoc, args...)
			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Dir = pwd
			output, cmdErr := cmd.CombinedOutput()
			if len(output) > 0 {
				fmt.Println(string(output))
			}
			if cmdErr != nil {
				fmt.Println(cmdErr)
				os.Exit(1)
			}
		}
	}
}
