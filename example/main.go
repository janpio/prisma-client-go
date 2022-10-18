package main

import (
	"context"
	"fmt"
	"github.com/prisma/prisma-client-go/engine"
	"io/ioutil"
)

const schmea = `
datasource db {
        // could be postgresql or mysql
        provider = "sqlite"
        url      = "file:dev.db"
    }

    generator db {
        provider = "go run github.com/prisma/prisma-client-go"
        // set the output folder and package name
        // output           = "./your-folder"
        // package          = "yourpackagename"
    }

    model Post {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        updatedAt DateTime @updatedAt
        title     String
        published Boolean
        desc      String?
    }
`

const mysqlSchema = `generator db {
  provider          = "go run github.com/prisma/prisma-client-go"
}

datasource db {
  provider = "mysql"
  url      = "mysql://root:shaoxiong123456@8.142.115.204:3306/main"
}

model oauth_user {
  id                  String    @id @db.VarChar(50)
  name                String?   @default("") @db.VarChar(50)
  nick_name           String?   @default("") @db.VarChar(50)
  user_name           String?   @unique(map: "name_index") @default("") @db.VarChar(50)
  encryption_password String?   @default("") @db.VarChar(250)
  mobile              String?   @default("") @db.VarChar(11)
  email               String?   @default("") @db.VarChar(50)
  mate_data           String?   @db.Text
  last_login_time     DateTime? @db.Timestamp(0)
  status              Int?      @default(0) @db.TinyInt
  create_time         DateTime? @default(now()) @db.Timestamp(0)
  update_time         DateTime? @db.Timestamp(0)
  is_del              Int?      @default(0) @db.UnsignedTinyInt
}
`

func main() {
	// if err := run(); err != nil {
	// 	panic(err)
	// }
	//migrationEngine := migrate.NewMigrationEngine()
	//
	//migrationEngine.Push("schema2.prisma")
	//migrationEngine.Push2("schema1.prisma")
	//migrationEngine.Push("schema2.prisma")
	//migrationEngine.Push2("schema2.prisma")

	//introspectionEngine := introspection.NewIntrospectEngine()
	//introspectionEngine.Pull("schema1.prisma")
	//ntrospectionEngine.Pull2("schema1.prisma")
	//introspectionEngine.Pull("schema1.prisma")
	//introspectionEngine.Pull2("schema1.prisma")
	ss, _ := ioutil.ReadFile("schema1.prisma")
	engine.Push(string(ss))
	//engine.Pull("schema2.prisma")
	// testDmmf()
	//engine.QueryDMMF(mysqlSchema)
	//testSdl1()
}

func testDmmf() {
	engine := engine.NewQueryEngine(schmea, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	dmmf, err := engine.IntrospectDMMF(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(dmmf.Datamodel)
}

var querySchema = `{ result:findFirstoauth_user {id name nick_name}}`

type OauthUser struct {
	ID                 string `json:"id"`                  // id
	Name               string `json:"name"`                // 姓名
	NickName           string `json:"nick_name"`           // 昵称
	UserName           string `json:"user_name"`           // 用户名
	EncryptionPassword string `json:"encryption_password"` // 加密后密码
	Mobile             string `json:"mobile"`              // 手机号
	Email              string `json:"email"`               // 邮箱
	LastLoginTime      string `json:"last_login_time"`     // 最后一次登陆时间
	Status             int64  `json:"status"`              // 状态
	MateData           string `json:"mate_data"`           // 其他信息(json字符串保存)
	CreateTime         string `json:"create_time"`         // 创建时间
	UpdateTime         string `json:"update_time"`         // 修改时间
	IsDel              int64  `json:"isDel"`               // 是否删除
}

//func testSdl1() {
//
//	queryEngine := engine.GetQueryEngineOnce(mysqlSchema)
//	ctx := context.TODO()
//	//var result OauthUser
//
//	var response OauthUser
//	payload := engine.GQLRequest{
//		Query:     querySchema,
//		Variables: map[string]interface{}{},
//	}
//	err := queryEngine.Do(ctx, payload, &response)
//	//result, err := engine.Do(ctx, querySchema)
//	//fmt.Print(result)
//	if err != nil {
//		panic(err)
//	}
//}

func testSdl() {
	engine := engine.NewQueryEngine(mysqlSchema, false)
	defer engine.Disconnect()
	if err := engine.Connect(); err != nil {
		panic(err)
	}
	ctx := context.TODO()
	var result OauthUser
	err := engine.Do(ctx, querySchema, result)
	//sdl, err := engine.IntrospectSDL(ctx)

	if err != nil {
		panic(err)
	}
	//fmt.Println(string(sdl))
}

// func run() error {
// 	client := db.NewClient()
// 	if err := client.Prisma.Connect(); err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err := client.Prisma.Disconnect(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	ctx := context.Background()

// 	// create a post
// 	createdPost, err := client.Post.CreateOne(
// 		db.Post.Title.Set("Hi from Prisma!"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, _ := json.MarshalIndent(createdPost, "", "  ")
// 	fmt.Printf("created post: %s\n", result)

// 	// find a single post
// 	post, err := client.Post.FindUnique(
// 		db.Post.ID.Equals(createdPost.ID),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, _ = json.MarshalIndent(post, "", "  ")
// 	fmt.Printf("post: %s\n", result)

// 	// for optional/nullable values, you need to check the function and create two return values
// 	// `desc` is a string, and `ok` is a bool whether the record is null or not. If it's null,
// 	// `ok` is false, and `desc` will default to Go's default values; in this case an empty string (""). Otherwise,
// 	// `ok` is true and `desc` will be "my description".
// 	desc, ok := post.Desc()
// 	if !ok {
// 		return fmt.Errorf("post's description is null")
// 	}

// 	fmt.Printf("The posts's description is: %s\n", desc)

// 	createUserA := client.Post.CreateOne(
// 		db.Post.Title.Set("2"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("222."),
// 	).Tx()

// 	createUserB := client.Post.CreateOne(
// 		db.Post.Title.Set("3"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("222."),
// 	).Tx()

// 	tx := client.Prisma.Transaction(createUserA, createUserB)
// 	if err := tx.Exec(ctx); err != nil {
// 		panic(err)
// 	}

// 	return nil
// }
