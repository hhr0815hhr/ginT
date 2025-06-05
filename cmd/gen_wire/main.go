package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// Data 传递给模板的数据
type Data struct {
	PackageName         string
	Imports             []string
	RepoProviders       []string
	LogicProviders      []string
	ControllerProviders []string
	RouteFuncParams     []Param
	LogicStructs        []LogicInfo
}

type LogicInfo struct {
	UpName    string
	LowName   string
	LogicType string
}

// Param 代表 ProvideRoutes 函数的参数
type Param struct {
	Name string
	Type string
}

type FileInfo struct {
	Prefix string
	Suffix string
}

func checkType(name string) FileInfo {
	if strings.HasSuffix(name, "_controller.go") {
		return FileInfo{Prefix: "New", Suffix: "Controller"}
	} else if strings.HasSuffix(name, "_logic.go") {
		return FileInfo{Prefix: "New", Suffix: "Logic"}
	}
	return FileInfo{Prefix: "New", Suffix: "Repo"}
}

func main() {
	// 设置要扫描的目录
	dirs := []string{
		"./internal/controller",
		"./internal/logic",
		"./internal/database/model",
	}

	// 设置 wire.go 的输出路径
	outputFile := "./internal/wire.go"

	// 设置当前包名 (需要根据你的项目结构调整)
	packageName := "internal"

	// 导入列表
	imports := []string{
		"github.com/hhr0815hhr/gint/internal/controller",
		"github.com/hhr0815hhr/gint/internal/database/model",
		"github.com/hhr0815hhr/gint/internal/database/mysql",
		"github.com/hhr0815hhr/gint/internal/logic",
		"github.com/hhr0815hhr/gint/internal/server/http",
	}

	repoProviders := []string{}
	logicProviders := []string{}
	controllerProviders := []string{}
	routeFuncParams := []Param{}
	logicStructs := []LogicInfo{}

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path %q: %v\n", path, err)
				return nil // 忽略错误，继续下一个目录
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {

				fileType := checkType(info.Name())

				fset := token.NewFileSet()
				node, err := parser.ParseFile(fset, path, nil, 0)
				if err != nil {
					fmt.Printf("Error parsing file %q: %v\n", path, err)
					return nil // 忽略错误，继续下一个文件
				}

				// 遍历文件中的声明
				for _, decl := range node.Decls {
					if funcDecl, ok := decl.(*ast.FuncDecl); ok {
						// 查找符合特定模式的函数，例如返回单个值且名称以 "New" 开头
						if strings.HasPrefix(funcDecl.Name.Name, fileType.Prefix) && strings.HasSuffix(funcDecl.Name.Name, fileType.Suffix) {
							//
							switch fileType.Suffix {
							case "Controller":
								controllerProviders = append(controllerProviders, "controller."+funcDecl.Name.Name)
								tmp := strings.TrimPrefix(funcDecl.Name.Name, fileType.Prefix)
								tmp = strings.TrimSuffix(tmp, fileType.Suffix)
								tmp = strings.ToLower(tmp)
								routeFuncParams = append(routeFuncParams, Param{
									Name: tmp,
									Type: "controller." + strings.TrimPrefix(funcDecl.Name.Name, fileType.Prefix)},
								)
							case "Logic":
								logicProviders = append(logicProviders, "logic."+funcDecl.Name.Name)
								logicStructs = append(logicStructs, LogicInfo{
									UpName:    toUpper(strings.TrimPrefix(funcDecl.Name.Name, fileType.Prefix)),
									LowName:   toLower(strings.TrimPrefix(funcDecl.Name.Name, fileType.Prefix)),
									LogicType: toUpper(strings.TrimPrefix(funcDecl.Name.Name, fileType.Prefix)),
								})
							case "Repo":
								repoProviders = append(repoProviders, "model."+funcDecl.Name.Name)
							}
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking directory:", err)
			continue
		}
	}
	// 准备模板数据
	data := Data{
		PackageName:         packageName,
		Imports:             imports,
		RepoProviders:       repoProviders,
		LogicProviders:      logicProviders,
		ControllerProviders: controllerProviders,
		RouteFuncParams:     routeFuncParams,
		LogicStructs:        logicStructs,
	}

	// 解析模板文件
	tmpl, err := template.New("wire").Parse(wireTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// 创建输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	// 执行模板并写入文件
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Successfully generated", outputFile)
}

func toLower(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	// 将第一个字符转换为小写
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func toUpper(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	// 将第一个字符转换为小写
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// wireTemplate 定义了 wire.go 文件的模板内容 (直接在代码中定义)
const wireTemplate = `//go:build wireinject
// +build wireinject

package {{.PackageName}}

//go:generate go run github.com/google/wire/cmd/wire
import (
{{- range .Imports}}
    "{{ . }}"
{{- end}}
    "github.com/gin-gonic/gin"
    "github.com/google/wire"
)

type AppInfo struct {
	Engine    *gin.Engine
	TestLogic *logic.TestLogic
    Data      map[string]interface{}
}

func ProvideApp(
	engine *gin.Engine,
{{- range .LogicStructs}}
    {{ .LowName }} *logic.{{ .LogicType }},
{{- end}}
) *AppInfo {
	return &AppInfo{
		Engine:          engine,
{{- range .LogicStructs}}
    	{{ .UpName }}: {{ .LowName }},
{{- end}}
        Data: map[string]interface{}{},
	}
}

func InitApp() *AppInfo {
    wire.Build(DbSet, RepoSet, LogicSet, RouteSet, ProvideRoutes, http.NewHTTPServer, ProvideApp)
    return nil
}

var DbSet = wire.NewSet(mysql.ProvideDB)

var RepoSet = wire.NewSet(
{{- range .RepoProviders}}
    {{ . }},
{{- end}}
)
var LogicSet = wire.NewSet(
{{- range .LogicProviders}}
    {{ . }},
{{- end}}
)

var RouteSet = wire.NewSet(
{{- range .ControllerProviders}}
    {{ . }},
{{- end}}
)

func ProvideRoutes(
{{- range $index, $provider := .RouteFuncParams}}
    {{ $provider.Name }} *{{ $provider.Type }},
{{- end}}
) *http.HTTPRoutes {
    return &http.HTTPRoutes{
       Routers: []controller.Router{
{{- range $provider := .RouteFuncParams}}
          {{ $provider.Name }},
{{- end}}
       },
    }
}
`
