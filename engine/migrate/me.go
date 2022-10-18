package migrate

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
	"github.com/prisma/prisma-client-go/logger"
)

func NewMigrationEngine() *MigrationEngine {
	// TODO:这里可以设置默认值
	engine := &MigrationEngine{
		// path: path,
	}
	file, err := engine.ensure() //确保引擎一定安装了
	if err != nil {
		panic(err)
	}
	engine.path = file
	return engine
}

type MigrationEngine struct {
	path string
}

// func (e *MigrationEngine) Name() string {
// 	return "migration-engine"
// }

func (e *MigrationEngine) ensure() (string, error) {
	ensureEngine := time.Now()

	dir := binaries.GlobalCacheDir()
	// 确保引擎一定下载了
	if err := binaries.FetchNative(dir); err != nil {
		return "", fmt.Errorf("could not fetch binaries: %w", err)
	}
	// check for darwin/windows/linux first
	//binaryName := platform.CheckForExtension(platform.Name(), platform.Name())
	binaryName := platform.BinaryPlatformName()

	var file string
	// forceVersion saves whether a version check should be done, which should be disabled
	// when providing a custom query engine value
	// forceVersion := true

	name := "prisma-migration-engine-"
	globalPath := path.Join(dir, binaries.EngineVersion, name+binaryName)

	logger.Debug.Printf("expecting global migration engine `%s` ", globalPath)

	// TODO write tests for all cases

	// first, check if the query engine binary is being overridden by PRISMA_MIGRATION_ENGINE_BINARY
	prismaQueryEngineBinary := os.Getenv("PRISMA_MIGRATION_ENGINE_BINARY")
	if prismaQueryEngineBinary != "" {
		logger.Debug.Printf("PRISMA_MIGRATION_ENGINE_BINARY is defined, using %s", prismaQueryEngineBinary)

		if _, err := os.Stat(prismaQueryEngineBinary); err != nil {
			return "", fmt.Errorf("PRISMA_MIGRATION_ENGINE_BINARY was provided, but no query engine was found at %s", prismaQueryEngineBinary)
		}

		file = prismaQueryEngineBinary
		// forceVersion = false
	} else {
		if _, err := os.Stat(globalPath); err == nil {
			logger.Debug.Printf("exact query engine found in global path")
			file = globalPath
		}
	}

	if file == "" {
		// TODO log instructions on how to fix this problem
		return "", fmt.Errorf("no binary found ")
	}

	logger.Debug.Printf("using migration engine at %s", file)
	logger.Debug.Printf("ensure migration engine took %s", time.Since(ensureEngine))

	return file, nil
}

func (e *MigrationEngine) Push(schemaPath string) error {
	startParse := time.Now()
	// 可以缓存到改引擎中？
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		err = fmt.Errorf("load prisma schema: %s", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.path, "--datamodel", schemaPath)

	pipe, err := cmd.StdinPipe() // 标准输入流
	if err != nil {
		err = fmt.Errorf("migration engine std in pipe: %s", err.Error())
		return err
		// return "", err
	}
	defer pipe.Close()
	// 构建一个json-rpc 请求参数
	req := MigrationRequest{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  "schemaPush",
		Params: MigrationRequestParams{
			Force:  true,
			Schema: string(schema),
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	// 入参追加到管道中
	_, err = pipe.Write(append(data, []byte("\n")...))
	if err != nil {
		// return "", err
		return err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("migration std out pipe: %s", err.Error())
	}
	r := bufio.NewReader(out)

	// 开始执行
	err = cmd.Start()
	if err != nil {
		return err
	}

	var response MigrationResponse

	outBuf := &bytes.Buffer{}
	// {\"jsonrpc\":\"2.0\",\"result\":{\"executedSteps\":1,\"unexecutable\":[],\"warnings\":[]},\"id\":1}\n
	// 这一段的意思是，每100ms读取一次结果，直到超时或有结果
	for {
		// 等待100 ms
		//time.Sleep(time.Millisecond * 100)
		b, err := r.ReadByte()
		if err != nil {
			err = fmt.Errorf("migration ReadByte: %s", err.Error())
		}
		err = outBuf.WriteByte(b)
		if err != nil {
			err = fmt.Errorf("migration writeByte: %s", err.Error())
		}

		if b == '\n' {
			// 解析响应结果
			err = json.Unmarshal(outBuf.Bytes(), &response)
			if err != nil {
				return err
			}
			if response.Error == nil {
				log.Println("Migration successful")
			}
			fmt.Print("ende ")
			break
		}
		// 如果超时了？跳出读取？
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	log.Printf("[timing] migrate took %s", time.Since(startParse))
	if response.Error != nil {
		return fmt.Errorf("migrate error: %s", response.Error.Data.Message)
	}
	return nil
}
