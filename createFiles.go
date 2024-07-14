package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func main() {
	err := CopyFile("key.txt", "nf.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	var splitNum int
	splitNum, err = GetSplitNum("key.txt", 8000)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = CreateSplitFile("nf", splitNum)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Executed successfully!")
}

func CopyFile(src string, dst string) (err error) {
	var srcFile *os.File
	srcFile, err = os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer func(srcFile *os.File) {
		err = srcFile.Close()
		if err != nil {
			return
		}
	}(srcFile)
	var info os.FileInfo
	info, err = srcFile.Stat()
	if err != nil {
		return
	}
	if info.IsDir() {
		err = fmt.Errorf("source is a directory")
		return
	}
	var dstFile *os.File
	dstFile, err = createFile(dst, info)
	if err != nil {
		return
	}
	defer func(dstFile *os.File) {
		err = dstFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(dstFile)
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return
	}
	return
}

func createFile(dst string, info os.FileInfo) (*os.File, error) {
	dir := filepath.Dir(dst)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	var dstFile *os.File
	dstFile, err = os.Create(dst)
	if err != nil {
		return nil, err
	}
	err = dstFile.Chmod(info.Mode())
	if err != nil {
		err = dstFile.Close()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, err
	}
	return dstFile, nil
}
func CreateSplitFile(pathName string, splitNum int) (err error) {
	dirPath := path.Join("../", pathName)
	err = os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		return
	}
	for i := range splitNum {
		fileName := pathName + "_" + strconv.Itoa(i) + ".txt"
		filePath := path.Join(dirPath, fileName)
		_, err = os.Create(filePath)
		if err != nil {
			return
		}
	}
	return nil
}

func GetFileLineNum(filename string) (lineCount int, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	return lineCount, nil
}

func GetSplitNum(filename string, lineNumPerFile int) (splitNum int, err error) {
	var line int
	line, err = GetFileLineNum(filename)
	if err != nil {
		return 0, err
	}
	return int(math.Ceil(float64(line) / float64(lineNumPerFile))), nil
}
