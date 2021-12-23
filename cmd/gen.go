package cmd

import (
	"bytes"
	"deploy_server/cmd/gen"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:               "gen",
	Short:             "gen controller",
	Long:              "gen controller",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run:               GeneratorController,
}

var (
	controllerName string
)

func init() {
	genCmd.AddCommand(gen.ModelCmd)

	genCmd.Flags().StringVarP(&controllerName, "cont", "", "", "[Required] The name of the controller name")

	rootCmd.AddCommand(genCmd)
}

type controllerModel struct {
	PkgName    string
	StructName string
}

func parseTemplateOrPanic(t string) *template.Template {
	tpl, err := template.New("cont_template").Parse(t)
	if err != nil {
		panic(err)
	}
	return tpl
}

var outputTemplate = parseTemplateOrPanic(`
package controller

import (
	"deploy_server/pkg/cache"
	"deploy_server/pkg/context"
	"deploy_server/pkg/core"
	"deploy_server/pkg/db"
)

type {{.PkgName}}Handler struct {
	cache cache.Repo
}

func (handle *{{.PkgName}}Handler) Test() context.HandlerFunc {
	return func(ctx context.Context) {
		ctx.Success("success")
	}
}

func (handle *{{.PkgName}}Handler) RegistryRouter(m *core.Mux) {
	authGroup := m.Group("/{{.PkgName}}")
	{
		authGroup.GET("/test", handle.Test())
	}
}

func New{{.StructName}}Controller(db db.Repo, cache cache.Repo) *projectHandler {
	return &{{.PkgName}}Handler{
		cache,
	}
}
`)

func capitalize(s string) string {
	var upperStr string
	chars := strings.Split(s, "_")
	for _, val := range chars {
		vv := []rune(val)
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				if vv[i] >= 97 && vv[i] <= 122 {
					vv[i] -= 32
					upperStr += string(vv[i])
				}
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

func GeneratorController(cmd *cobra.Command, args []string) {
	if controllerName == "" {
		_ = cmd.Usage()
		return
	}
	fmt.Printf("Generator controller " + controllerName)
	structName := capitalize(controllerName)
	values := controllerModel{controllerName, structName}
	buf := new(bytes.Buffer)
	err := outputTemplate.Execute(buf, values)
	if err != nil {
		panic(err)
	}

	//格式化代码
	formattedOutput, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	buf = bytes.NewBuffer(formattedOutput)

	outDir := "controller/"
	//输出文件
	filename := fmt.Sprintf("%s/%s_controller.go", outDir, controllerName)
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0777); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("  └── generator controller: %s, file: %s\n", controllerName, filename)
}
