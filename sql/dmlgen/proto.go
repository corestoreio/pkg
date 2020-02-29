package dmlgen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corestoreio/errors"
)

// ProtocOptions allows to modify the protoc CLI command.
type ProtocOptions struct {
	BuildTags          []string
	WorkingDirectory   string
	ProtoGen           string // default gofast, options: gogo, gogofast, gogofaster
	GRPC               bool
	GRPCGatewayOutMap  []string // GRPC must be enabled in the above field
	GRPCGatewayOutPath string   // GRPC must be enabled in the above field
	ProtoPath          []string
	GoGoOutPath        string
	GoGoOutMap         []string
	SwaggerOutPath     string
	CustomArgs         []string
	// TODO add validation plugin, either
	//  https://github.com/mwitkow/go-proto-validators as used in github.com/gogo/grpc-example/proto/example.proto
	//  This github.com/mwitkow/go-proto-validators seems dead.
	//  or https://github.com/envoyproxy/protoc-gen-validate
	//  Requirement: error messages must be translatable and maybe use an errors.Kind type
}

var defaultProtoPaths = make([]string, 0, 8)

func init() {
	preDefinedPaths := [...]string{
		build.Default.GOPATH + "/src/",
		build.Default.GOPATH + "/src/github.com/gogo/protobuf/protobuf/",
		build.Default.GOPATH + "/src/github.com/gogo/googleapis/",
		"vendor/github.com/grpc-ecosystem/grpc-gateway/",
		"vendor/github.com/gogo/googleapis/",
		"vendor/",
		".",
	}
	for _, pdp := range preDefinedPaths {
		if _, err := os.Stat(pdp); !os.IsNotExist(err) {
			defaultProtoPaths = append(defaultProtoPaths, pdp)
		}
	}
}

func (po *ProtocOptions) toArgs() []string {
	gogoOut := make([]string, 0, 4)
	if po.GRPC {
		gogoOut = append(gogoOut, "plugins=grpc")
		if po.GRPCGatewayOutMap == nil {
			po.GRPCGatewayOutMap = []string{
				"allow_patch_feature=false",
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
				"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
				"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			}
		}
		if po.GRPCGatewayOutPath == "" {
			po.GRPCGatewayOutPath = "."
		}
	}
	if po.GoGoOutPath == "" {
		po.GoGoOutPath = "."
	}
	if po.GoGoOutMap == nil {
		po.GoGoOutMap = []string{
			"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
			"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
		}
	}
	gogoOut = append(gogoOut, po.GoGoOutMap...)

	if po.ProtoPath == nil {
		po.ProtoPath = append(po.ProtoPath, defaultProtoPaths...)
	}

	if po.ProtoGen == "" {
		po.ProtoGen = "gofast"
	} else {
		switch po.ProtoGen {
		case "gofast", "gogo", "gogofast", "gogofaster":
			// ok
		default:
			panic(fmt.Sprintf("[dmlgen] ProtoGen CLI command %q not supported, allowed: gofast, gogo, gogofast, gogofaster", po.ProtoGen))
		}
	}

	// To generate PHP Code replace `gogo_out` with `php_out`.
	// Java bit similar. Java has ~15k LOC, Go ~3.7k
	args := []string{
		"--" + po.ProtoGen + "_out", fmt.Sprintf("%s:%s", strings.Join(gogoOut, ","), po.GoGoOutPath),
		"--proto_path", strings.Join(po.ProtoPath, ":"),
	}
	if po.GRPC && len(po.GRPCGatewayOutMap) > 0 {
		args = append(args, "--grpc-gateway_out="+strings.Join(po.GRPCGatewayOutMap, ",")+":"+po.GRPCGatewayOutPath)
	}
	if po.SwaggerOutPath != "" {
		args = append(args, "--swagger_out="+po.SwaggerOutPath)
	}
	return append(args, po.CustomArgs...)
}

func (po *ProtocOptions) chdir() (deferred func(), _ error) {
	deferred = func() {}
	if po.WorkingDirectory != "" {
		oldWD, err := os.Getwd()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := os.Chdir(po.WorkingDirectory); err != nil {
			return nil, errors.Wrapf(err, "[dmlgen] Failed to chdir to %q", po.WorkingDirectory)
		}
		deferred = func() {
			_ = os.Chdir(oldWD)
		}
	}
	return deferred, nil
}

// GenerateProto searches all *.proto files in the given path and calls protoc
// to generate the Go source code.
func GenerateProto(protoFilesPath string, po *ProtocOptions) error {
	restoreFn, err := po.chdir()
	if err != nil {
		return errors.WithStack(err)
	}
	defer restoreFn()

	protoFilesPath = filepath.Clean(protoFilesPath)
	if ps := string(os.PathSeparator); !strings.HasSuffix(protoFilesPath, ps) {
		protoFilesPath += ps
	}

	protoFiles, err := filepath.Glob(protoFilesPath + "*.proto")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access proto files in path %q", protoFilesPath)
	}

	cmd := exec.Command("protoc", append(po.toArgs(), protoFiles...)...)
	cmdStr := fmt.Sprintf("\ncd %s && %s\n\n", po.WorkingDirectory, cmd)
	if isDebug() {
		if po.WorkingDirectory == "" {
			po.WorkingDirectory = "."
		}
		print(cmdStr)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] %s%s", out, cmdStr)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.Contains(text, "WARNING") {
			return errors.WriteFailed.Newf("[dmlgen] protoc Error: %s", text)
		}
	}

	// what a hack: find all *.pb.go files and remove `import null
	// "github.com/corestoreio/pkg/storage/null"` because no other way to get
	// rid of the unused import or reference that import somehow in the
	// generated file :-( Once there's a better solution, remove this code.
	pbGoFiles, err := filepath.Glob(protoFilesPath + "*.pb.*go")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access pb.go files in path %q", protoFilesPath)
	}

	removeImports := [][]byte{
		[]byte("import null \"github.com/corestoreio/pkg/storage/null\"\n"),
		[]byte("null \"github.com/corestoreio/pkg/storage/null\"\n"),
	}
	for _, file := range pbGoFiles {
		fContent, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, ri := range removeImports {
			fContent = bytes.Replace(fContent, ri, nil, -1)
		}

		var buf bytes.Buffer
		for _, bt := range po.BuildTags {
			fmt.Fprintf(&buf, "// +build %s\n", bt)
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
			buf.Write(fContent)
			fContent = buf.Bytes()
		}

		if err := ioutil.WriteFile(file, fContent, 0644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
