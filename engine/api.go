package engine

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/engine/introspection"
	"github.com/prisma/prisma-client-go/engine/migrate"
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"
	"github.com/prisma/prisma-client-go/logger"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func NewDMFQueryEngine(schema string) (*QueryEngine, error) {
	content, err := Pull(schema)
	if err != nil {
		return nil, err
	}
	queryEngine := NewQueryEngine(content, false)
	if err := queryEngine.ConnectSDK(); err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return nil, err
	}
	return queryEngine, nil
}

var globalQueryEngine *QueryEngine

func GetQueryEngineOnce(schema string) *QueryEngine {
	if globalQueryEngine == nil {
		content, err := Pull(schema)
		if err != nil {
			logger.Debug.Printf("connect fail err : ", err)
		}
		globalQueryEngine = NewQueryEngine(content, false)
		if err := globalQueryEngine.ConnectSDK(); err != nil {
			logger.Debug.Printf("connect fail err : ", err)
		}
	}
	return globalQueryEngine
}

func ReloadQueryEngineOnce(schema string) (*QueryEngine, error) {
	// 先释放掉老的资源
	if globalQueryEngine != nil {
		globalQueryEngine.Disconnect()
		globalQueryEngine = nil
	}
	// 内省
	content, err := Pull(schema)
	if err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return nil, err
	}

	globalQueryEngine = NewQueryEngine(content, false)
	if err := globalQueryEngine.ConnectSDK(); err != nil {
		logger.Debug.Printf("connect fail err : ", err)
	}
	return globalQueryEngine, nil
}

func Push(schemaPath string) error {
	migrationEngine := migrate.NewMigrationEngine()
	return migrationEngine.Push(schemaPath)
}

func Pull(schema string) (string, error) {
	migrationEngine := introspection.NewIntrospectEngine()
	return migrationEngine.Pull(schema)
}

func InitQueryEngine(schema string) error {
	content, err := Pull(schema)
	if err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	globalQueryEngine = NewQueryEngine(content, false)
	if err := globalQueryEngine.ConnectSDK(); err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	return nil
}

func QuerySchema(querySchema string, result interface{}) error {
	ctx := context.TODO()
	payload := GQLRequest{
		Query:     querySchema,
		Variables: map[string]interface{}{},
	}
	err := globalQueryEngine.Do(ctx, payload, result)
	if err != nil {
		return err
	}
	return nil
}

func QuerySDL(dbSchemaPath, sdlSchema string, result interface{}) error {
	queryEngine := GetQueryEngineOnce(dbSchemaPath)
	ctx := context.TODO()
	payload := GQLRequest{
		Query:     sdlSchema,
		Variables: map[string]interface{}{},
	}
	err := queryEngine.Do(ctx, payload, result)
	if err != nil {
		return err
	}
	return nil
}

func QueryDMMF(dbSchemaPath string) (*dmmf.Document, error) {
	queryEngine := GetQueryEngineOnce(dbSchemaPath)
	return queryEngine.IntrospectDMMF(context.TODO())
}

func (e *QueryEngine) ensureSDK() (string, error) {
	ensureEngine := time.Now()

	dir := binaries.GlobalCacheDir()
	// 确保引擎一定下载了
	if err := binaries.FetchNative(dir); err != nil {
		return "", fmt.Errorf("could not fetch binaries: %w", err)
	}
	binariesPath := path.Join(dir, binaries.EngineVersion)
	//binaryName := platform.CheckForExtension(platform.Name(), platform.Name())
	binaryName := platform.BinaryPlatformName()

	exactBinaryName := platform.CheckForExtension(platform.Name(), platform.BinaryPlatformName())

	var file string
	forceVersion := true

	name := "prisma-query-engine-"
	localPath := path.Join("./", name+binaryName)
	localExactPath := path.Join("./", name+exactBinaryName)
	globalPath := path.Join(binariesPath, name+binaryName)
	globalExactPath := path.Join(binariesPath, name+exactBinaryName)

	logger.Debug.Printf("expecting local query engine `%s` or `%s`", localPath, localExactPath)
	logger.Debug.Printf("expecting global query engine `%s` or `%s`", globalPath, globalExactPath)

	// TODO write tests for all cases

	// first, check if the query engine binary is being overridden by PRISMA_QUERY_ENGINE_BINARY
	prismaQueryEngineBinary := os.Getenv("PRISMA_QUERY_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.Debug.Printf("PRISMA_QUERY_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)

		if _, err := os.Stat(prismaQueryEngineBinary); err != nil {
			return "", fmt.Errorf("PRISMA_QUERY_ENGINE_BINARY was provided, but no query engine was found at %s", prismaQueryEngineBinary)
		}

		file = prismaQueryEngineBinary
		forceVersion = false
	} else {
		if _, err := os.Stat(localExactPath); err == nil {
			logger.Debug.Printf("exact query engine found in working directory")
			file = localExactPath
		} else if _, err := os.Stat(localPath); err == nil {
			logger.Debug.Printf("query engine found in working directory")
			file = localPath
		}

		if _, err := os.Stat(globalExactPath); err == nil {
			logger.Debug.Printf("query engine found in global path")
			file = globalExactPath
		} else if _, err := os.Stat(globalPath); err == nil {
			logger.Debug.Printf("exact query engine found in global path")
			file = globalPath
		}
	}

	if file == "" {
		// TODO log instructions on how to fix this problem
		return "", fmt.Errorf("no binary found ")
	}

	startVersion := time.Now()
	out, err := exec.Command(file, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("version check failed: %w", err)
	}
	logger.Debug.Printf("version check took %s", time.Since(startVersion))

	if v := strings.TrimSpace(strings.Replace(string(out), "query-engine", "", 1)); binaries.EngineVersion != v {
		note := "Did you forget to run `go run github.com/prisma/prisma-client-go generate`?"
		msg := fmt.Errorf("expected query engine version `%s` but got `%s`\n%s", binaries.EngineVersion, v, note)
		if forceVersion {
			return "", msg
		}

		logger.Info.Printf("%s, ignoring since custom query engine was provided", msg)
	}

	logger.Debug.Printf("using query engine at %s", file)
	logger.Debug.Printf("ensure query engine took %s", time.Since(ensureEngine))

	return file, nil
}

func (e *QueryEngine) ConnectSDK() error {
	logger.Debug.Printf("ensure query engine binary...")

	_ = godotenv.Load("e2e.env")
	_ = godotenv.Load("db/e2e.env")
	_ = godotenv.Load("prisma/e2e.env")

	startEngine := time.Now()

	file, err := e.ensureSDK()
	if err != nil {
		return fmt.Errorf("ensure: %w", err)
	}

	if err := e.spawn(file); err != nil {
		return fmt.Errorf("spawn: %w", err)
	}

	logger.Debug.Printf("connecting took %s", time.Since(startEngine))
	logger.Debug.Printf("connected.")

	return nil
}
