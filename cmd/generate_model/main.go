package main

import (
	"io"
	"os"
	"strings"

	"github.com/real-web-world/go-api/pkg/bdk"
)

const (
	tplPath = "models/tpl.go"
	outPath = "models/"
)

func main() {
	modelNameList := []string{"Article", "ArticleProfilePicture", "ArticleTag"}
	tplFile, err := os.Open(tplPath)
	if err != nil {
		panic("tpl don't exist")
	}
	defer tplFile.Close()
	tplBytes, err := io.ReadAll(tplFile)
	if err != nil {
		panic("read tpl failed")
	}
	tplStr := bdk.Bytes2Str(tplBytes)
	for _, modelName := range modelNameList {
		modelNameBytes := ([]byte)(modelName)
		modelNameBytes[0] += 32
		modelFileName := string(modelNameBytes) + ".go"
		modelFilePath := outPath + modelFileName
		if bdk.IsFile(modelFilePath) {
			panic(modelName + "is exist model file")
		}
		currModel, err := os.OpenFile(modelFilePath, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			panic("create" + modelName + " model failed")
		}
		defer currModel.Close()
		modelStr := ""
		modelStr = strings.ReplaceAll(tplStr, "// +build ignore\n\n", "")
		modelStr = strings.ReplaceAll(modelStr, "Tpl", modelName)
		modelStr = strings.ReplaceAll(modelStr, "_ *ginApp.App", "app *ginApp.App")
		_, err = currModel.WriteString(modelStr)
		if err != nil {
			panic("write" + modelName + " model failed")
		}
	}
}
