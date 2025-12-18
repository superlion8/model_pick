package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func main() {
	http.HandleFunc("/", serveReviewPage)
	http.HandleFunc("/api/models", handleModels)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(".."))))

	port := ":8081"
	fmt.Printf("🌐 模特评审系统已启动: http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func serveReviewPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "review.html")
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	models := []Model{}

	// 扫描所有品牌文件夹
	rootDir := ".."
	items, err := os.ReadDir(rootDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
