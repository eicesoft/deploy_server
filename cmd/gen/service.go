package gen

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"strings"
)

var ServiceCmd = &cobra.Command{
	Use:               "service [serviceName]",
	Short:             "gen service serviceName",
	Long:              "gen service serviceName",
	Args:              cobra.MinimumNArgs(1),
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run:               GeneratorService,
}

var (
	serviceName string
)

var outputServiceTemplate = parseTemplateOrPanic(`package service

import (
	"context"
	"deploy_server/model/{{.ModelName}}"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
)


type Wrap{{.ServiceName}}BuilderFunc func(*{{.ModelName}}.QueryBuilder) *{{.ModelName}}.QueryBuilder

type {{.ServiceName}}Service struct {
	service
}

// Get 利用warp函数进行Builder构造条件获得单条数据
func (s *{{.ServiceName}}Service) Get(ctx context.Context, warp Wrap{{.ServiceName}}BuilderFunc) (p *{{.ModelName}}.{{.ServiceName}}, err error) {
	p, err = warp({{.ModelName}}.NewQueryBuilder()).
		QueryOne(s.GetDBReader(ctx))

	return
}

// GetAll 利用warp函数进行Builder构造条件获得列表数据
func (s *{{.ServiceName}}Service) GetAll(ctx context.Context, warp Wrap{{.ServiceName}}BuilderFunc) ([]*{{.ModelName}}.{{.ServiceName}}, error) {
	list, err := warp({{.ModelName}}.NewQueryBuilder()).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// GetPages 利用warp函数进行Builder构造条件获得分页数据
func (s *{{.ServiceName}}Service) GetPages(ctx context.Context, page int, pageSize int, warp Wrap{{.ServiceName}}BuilderFunc) ([]*{{.ModelName}}.{{.ServiceName}}, error) {
	list, err := warp({{.ModelName}}.NewQueryBuilder()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		QueryAll(s.GetDBReader(ctx))

	return list, err
}

// Create 创建数据
func (s *{{.ServiceName}}Service) Create(ctx context.Context, request interface{}) (id int32, err error) {
	{{.ModelName}}Model := {{.ModelName}}.NewModel()
	core.StructCopy({{.ModelName}}Model, request)
	id, err = {{.ModelName}}Model.Create(s.GetDBWriter(ctx))

	return
}

// Update 利用warp函数进行Builder构造更新条件进行更新
func (s *{{.ServiceName}}Service) Update(ctx context.Context, data map[string]interface{}, warp Wrap{{.ServiceName}}BuilderFunc) (err error) {
	err = warp({{.ModelName}}.NewQueryBuilder()).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// UpdateById 按ID更新数据
func (s *{{.ServiceName}}Service) UpdateById(ctx context.Context, data map[string]interface{}, id int32) (err error) {
	err = {{.ModelName}}.NewQueryBuilder().
		WhereId({{.ModelName}}.EqualPredicate, id).
		Updates(s.GetDBWriter(ctx), data)

	return
}

// Delete 按ID删除数据
func (s *{{.ServiceName}}Service) Delete(ctx context.Context, id int32) (err error) {
	{{.ModelName}}Model := {{.ModelName}}.NewModel()
	{{.ModelName}}Model.Id = id
	err = {{.ModelName}}Model.Delete(s.GetDBWriter(ctx))
	return
}

// New{{.ServiceName}}Service {{.ServiceName}}Service构造函数
func New{{.ServiceName}}Service(db db.Repo) *{{.ServiceName}}Service {
	return &{{.ServiceName}}Service{
		service{
			db,
		},
	}
}
`)

type ServiceValues struct {
	ModelName   string
	ServiceName string
}

func GeneratorService(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		return
	}

	serviceName = args[0]

	values := ServiceValues{ModelName: strings.ToLower(serviceName), ServiceName: serviceName}
	buf := new(bytes.Buffer)
	err := outputServiceTemplate.Execute(buf, values)
	if err != nil {
		panic(err)
	}

	//格式化代码
	//formattedOutput, err := format.Source(buf.Bytes())
	//if err != nil {
	//	panic(err)
	//}
	//buf = bytes.NewBuffer(formattedOutput)

	outDir := "service"
	//输出文件
	filename := fmt.Sprintf("%s/%s_service.go", outDir, strings.ToLower(serviceName))
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0777); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("  └── generator service: %s, file: %s\n", serviceName, filename)
}
