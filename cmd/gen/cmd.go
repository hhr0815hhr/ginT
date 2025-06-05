package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const controllerTemplate = `package controller

import (
	"github.com/gin-gonic/gin"
)

type {{.ControllerName}} struct {}

func New{{.ControllerName}}() *{{.ControllerName}} {
	return &{{.ControllerName}}{}
}

var _ Router = (*{{.ControllerName}})(nil)

func (c *{{.ControllerName}}) RegisterRoute(r *gin.Engine) {
	// 在这里定义你的路由
}
`

const logicTemplate = `package logic

type {{.LogicName}} struct {
}

func New{{.LogicName}}() *{{.LogicName}} {
	return &{{.LogicName}}{}
}`

const modelTemplate = `package model

import (
	"github.com/hhr0815hhr/gint/internal/database"
	"gorm.io/gorm"
)

type {{.ModelName}} struct {
}

type {{.ModelName}}Repo struct {
	*database.BaseRepository[{{.ModelName}}]
}

func New{{.ModelName}}Repo(db *gorm.DB) *{{.ModelName}}Repo {
	return &{{.ModelName}}Repo{
		BaseRepository: database.NewBaseRepository[{{.ModelName}}](db),
	}
}`

var GenCmd = &cobra.Command{
	Use:   "gen [type] [name]",
	Short: "生成模版",
	Long:  "一键生成模版",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[1]
		switch args[0] {
		case "controller":
			controllerName := strings.Title(name) + "Controller"
			fileName := strings.ToLower(name) + "_controller.go"
			data := struct {
				ControllerName string
			}{
				ControllerName: controllerName,
			}
			tmpl, err := template.New("controller").Parse(controllerTemplate)
			if err != nil {
				fmt.Println("Error parsing template:", err)
				os.Exit(1)
			}

			outputDir := "internal/controller" // 默认在当前目录生成，可以添加 flag 控制输出目录
			outputPath := filepath.Join(outputDir, fileName)

			file, err := os.Create(outputPath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				fmt.Println("Error executing template:", err)
				os.Exit(1)
			}
			fmt.Printf("Generated controller: %s\n", outputPath)
		case "logic":
			logicName := strings.Title(name) + "Logic"
			fileName := strings.ToLower(name) + "_logic.go"
			data := struct {
				LogicName string
			}{
				LogicName: logicName,
			}
			tmpl, err := template.New("logic").Parse(logicTemplate)
			if err != nil {
				fmt.Println("Error parsing template:", err)
				os.Exit(1)
			}

			outputDir := "internal/logic" // 默认在当前目录生成，可以添加 flag 控制输出目录
			outputPath := filepath.Join(outputDir, fileName)

			file, err := os.Create(outputPath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				fmt.Println("Error executing template:", err)
				os.Exit(1)
			}
			fmt.Printf("Generated logic: %s\n", outputPath)
		case "model":
			modelName := strings.Title(name)
			fileName := modelName + ".go"
			data := struct {
				ModelName string
			}{
				ModelName: modelName,
			}
			tmpl, err := template.New("model").Parse(modelTemplate)
			if err != nil {
				fmt.Println("Error parsing template:", err)
				os.Exit(1)
			}

			outputDir := "internal/database/model" // 默认在当前目录生成，可以添加 flag 控制输出目录
			outputPath := filepath.Join(outputDir, fileName)

			file, err := os.Create(outputPath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				fmt.Println("Error executing template:", err)
				os.Exit(1)
			}
			fmt.Printf("Generated model: %s\n", outputPath)
		default:
			cobra.CheckErr("请输入正确的类型")
		}
	},
}
