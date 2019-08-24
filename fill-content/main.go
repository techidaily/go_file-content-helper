package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Args struct {
	Dir     string
	Suffix  string
	Content string
	File    string
}

var args = new(Args)

func init() {
	flag.StringVar(&args.Dir, "d", "", "choose a directory")
	flag.StringVar(&args.Suffix, "s", "", "file suffix")
	flag.StringVar(&args.Content, "c", "", "content to fill")
	flag.StringVar(&args.File, "f", "", "file to fill")
}

func main() {
	flag.Parse()

	if !checkArgs(args) {
		flag.Usage()
		return
	}
	args.Dir = formatFiles(args.Dir)

	var to_filled string
	if args.Content != "" {
		to_filled = args.Content
	} else if args.File != "" {
		to_filled = readFile(args.File)
	} else {
		log.Fatal("error ...")
	}
	to_filled += "\n"

	var filelist = getFileList(args.Dir, args.Suffix)
	var backup_suffix = fmt.Sprintf("%d", time.Now().Unix())
	for _, v := range filelist {
		var err = appendFileHead(v, to_filled, backup_suffix)
		if err != nil {
			fmt.Printf("file %s append failed! err: %s\n", v, err.Error())
		} else {
			fmt.Printf("file %s append success\n", v)
		}
	}

}

func checkArgs(args *Args) bool {
	return args.Dir != "" &&
		args.Suffix != "" &&
		(args.Content != "" || args.File != "")
}

func readFile(filename string) string {
	var bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("failed to open file: ", filename, "\terr: ", err.Error())
	}

	return string(bytes)
}

func formatFiles(file string) string {
	var file_info, err = os.Stat(file)
	if err != nil {
		log.Fatal("wrong file: ", file, "\terr: ", err.Error())
	}
	return file_info.Name()
}

func getFileList(dir, suffix string) []string {
	var files, err = ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("wrong dir: ", dir, "\terr: ", err.Error())
	}

	var result []string
	var sep = string(os.PathSeparator)
	for _, v := range files {
		var file = dir + sep + v.Name()
		if v.IsDir() {
			var t = getFileList(file, suffix)
			result = append(result, t...)
		} else {
			if strings.HasSuffix(v.Name(), suffix) {
				result = append(result, file)
			}
		}
	}

	return result
}

func appendFileHead(filename, content, suffix string) error {
	var bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 备份
	err = ioutil.WriteFile(filename+suffix, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err == nil {

	}
	return err
}
