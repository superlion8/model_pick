package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Model struct {
	ID        string `json:"id"`
	Brand     string `json:"brand"`
	Name      string `json:"name"`
	Original  string `json:"original"`
	Generated string `json:"generated"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	models := []Model{}

	// 扫描所有品牌文件夹
	rootDir := "../.."
	items, err := os.ReadDir(rootDir)
	if err != nil {
		// 如果在Vercel环境中无法访问文件系统，返回空数组
		json.NewEncoder(w).Encode(models)
		return
	}

	for _, item := range items {
		if !item.IsDir() || item.Name() == "batch_processor" || item.Name() == "all_models" || strings.HasPrefix(item.Name(), ".") {
			continue
		}

		brandName := item.Name()
		brandDir := filepath.Join(rootDir, brandName)

		// 扫描品牌文件夹中的子文件夹
		subdirs, err := os.ReadDir(brandDir)
		if err != nil {
			continue
		}

		for _, subdir := range subdirs {
			if !subdir.IsDir() {
				continue
			}

			modelDir := filepath.Join(brandDir, subdir.Name())
			originalPath := filepath.Join(modelDir, "original_model.png")
			generatedPath := filepath.Join(modelDir, "model_new.png")

			// 检查两个文件是否都存在
			if _, err := os.Stat(originalPath); err == nil {
				if _, err := os.Stat(generatedPath); err == nil {
					models = append(models, Model{
						ID:        fmt.Sprintf("%s_%s", brandName, subdir.Name()),
						Brand:     brandName,
						Name:      subdir.Name(),
						Original:  fmt.Sprintf("/images/%s/%s/original_model.png", brandName, subdir.Name()),
						Generated: fmt.Sprintf("/images/%s/%s/model_new.png", brandName, subdir.Name()),
					})
				}
			}
		}
	}

	json.NewEncoder(w).Encode(models)
}
