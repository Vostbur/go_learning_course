package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return formatTree(out, path, printFiles, "")
}

func formatTree(out io.Writer, path string, printFiles bool, prefix string) error {
	// рекурсивно дерево каталогов с символами псевдографики

	// получаем список файлов по пути path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// если не указан в командной строке ключ -f оставляем в списке файлов только каталоги
	if !printFiles {
		files = filterFiles(files, func(file os.FileInfo) bool { return file.IsDir() })
	}

	for count, file := range files {
		isLastFile := count == len(files)-1
		fmt.Fprintln(out, formatOutput(file, prefix, isLastFile))

		if file.IsDir() {
			devideSimbol := "│"
			if isLastFile {
				devideSimbol = ""
			}
			_ = formatTree(out, path+string(os.PathSeparator)+file.Name(), printFiles, prefix+devideSimbol+"\t")
		}
	}

	return nil
}

func filterFiles(files []os.FileInfo, compareFunc func(os.FileInfo) bool) []os.FileInfo {
	// фильтруем список файлов по заданному параметру
	filesFiltered := make([]os.FileInfo, 0)

	for _, file := range files {
		if compareFunc(file) {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return filesFiltered
}

func formatOutput(file os.FileInfo, prefix string, isLastFile bool) string {
	// форматируем строку вывода символами псевдографики
	leftSimbol := "├"
	if isLastFile {
		leftSimbol = "└"
	}
	outputResult := prefix + leftSimbol + "───" + file.Name()

	if !file.IsDir() {
		fileSize := file.Size()

		if fileSize == 0 {
			outputResult += " (empty)"
		} else {
			outputResult += fmt.Sprintf(" (%vb)", fileSize)
		}
	}

	return outputResult
}
