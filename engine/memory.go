package engine

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type FileSize int64

const (
	B FileSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

var (
	dir = "images/"
)

func Download(content string) (scheduler Scheduler) {

	scheduler = Scheduler{}

	if err := Mkdir(dir); err != nil {
		return
	}

	fileName := filepath.Base(content)

	ext := filepath.Ext(fileName)
	if len(ext) <= 0 {
		fileName += ".png"
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	request, err := http.NewRequest(http.MethodGet, content, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		scheduler.Processors = append(scheduler.Processors, Processor{
			Content: content,
			FUNC:    Download,
		})
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Dir(dir) + string(os.PathSeparator) + fileName)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	reader := bufio.NewReaderSize(resp.Body, 64*1024)
	size, err := io.Copy(writer, reader)
	if err != nil {
		return
	}
	fmt.Println(fileName, FileByte(FileSize(size)))
	return
}

// 创建目录
func Mkdir(dir string) error {
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// 计算文件大小
func FileByte(size FileSize) string {
	switch {
	default:
		return fmt.Sprintf("%.2f%s", float64(size)/float64(B), "B")
	case size >= KB && size < MB:
		return fmt.Sprintf("%.2f%s", float64(size)/float64(KB), "KB")
	case size >= MB && size < GB:
		return fmt.Sprintf("%.2f%s", float64(size)/float64(MB), "MB")
	case size >= GB && size < TB:
		return fmt.Sprintf("%.2f%s", float64(size)/float64(GB), "GB")
	}
}
