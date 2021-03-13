package main

import (
	"flag"
	"fmt"
	"golf/core"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"bufio"
)

type Options struct {
	Dir              string
	IsRecur          bool
	Type             string
	Depth            uint8
	SizeFilter       string
	SortBy           string
	Name             string
	FileSuffix       string
	PrintFullPath    bool
	ActionRemoveFile bool
	ActionRenameFile bool
	RenameSedRule    string
	Startswith       string
}

type FileInfoExt struct {
	FInfo        os.FileInfo
	RelativePath string
}

var platform string = core.CheckOS()

func sizeFilterValidater() (msg string, ret bool) {
	return "", true
}

func checkDirExists(d string) (msg string, ret bool) {
	ds, err := os.Stat(d)
	if err != nil {
		return fmt.Sprintf("no such file or directory: %s", d), false
	}
	if !ds.IsDir() {
		return fmt.Sprintf("should be a directory: %s", d), false
	}
	return "", true
}

func matchFileName(rg string, fname string) bool {
	// r, _ := regexp.Compile(rg)
	b, err := regexp.MatchString(rg, fname)
	if err != nil {
		fmt.Println(err.Error(), "use go regex")
		os.Exit(1)
	}
	return b
}

func matchFileSuffix(suffix string, fname string) bool {
	ss := strings.Split(fname, ".")
	s := ss[len(ss)-1]
	return strings.EqualFold(strings.ToLower(suffix), strings.ToLower(s))
}

func matchFileStartswith(s string, fname string) bool {
	return strings.HasPrefix(fname, s)
}


func parseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.IsRecur, "recur", false, "是否递归")
	flag.StringVar(&options.Dir, "d", "", "扫描目录, 如果不提供则扫描当前工作目录")
	// flag.StringVar(&options.Type, "t", "", "-t=f or -t=d")
	flag.StringVar(&options.SortBy, "s", "", "-s=time -s=size")
	flag.StringVar(&options.Name, "name", "", "-name=foo")
	flag.StringVar(&options.FileSuffix, "suffix", "", "-suffix=zip 注意 程序只取split('.')[-1]")
	flag.BoolVar(&options.PrintFullPath, "fullpath", false, "输出绝对路径")
	flag.BoolVar(&options.ActionRemoveFile, "rm", false, "[!] 删除文件")
	flag.StringVar(&options.Startswith, "startswith", "", "startswith string")
	flag.Parse()
	if options.Dir == "" {
		options.Dir, _ = os.Getwd()
	}
	m, r := checkDirExists(options.Dir)
	if !r {
		fmt.Printf("%s\n", m)
		os.Exit(1)
	}
	return options
}

func sortBy(files *[]FileInfoExt, key string) {
	if key == "size" {
		sort.Slice(*files, func(i, j int) bool {
			return (*files)[i].FInfo.Size() > (*files)[j].FInfo.Size()
		})
	}
}

func sizeHumanRead(b int64) string {
	return ""
}

func walkDir(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error(), dir)
		os.Exit(1)
	}
	return files
}

func readFileInfoTime(finfo os.FileInfo) string {
	return fmt.Sprintf("%v", finfo.ModTime().Format("2006-01-02 15:04:05"))
}

func end(result []FileInfoExt, fullpath bool, pdir string) {
	for _, f := range result {
		if fullpath {
			// absPath, _ := filepath.Abs(f.FInfo.Name())
			// fmt.Printf("%s\n", absPath)
			fmt.Printf("%s\n", f.RelativePath)
		} else {
			fmt.Printf("[+] %v - %s - %s\n", readFileInfoTime(f.FInfo), core.Bytes(uint64(f.FInfo.Size())), f.RelativePath)
		}

	}
}

func realRemove(files []FileInfoExt) {
	for _, f := range files {
		err := os.Remove(f.RelativePath)
		if err != nil {
			fmt.Printf("[*] error when remove %s\n", f.RelativePath)
		}
		fmt.Printf("[-] %s removed\n", f.RelativePath)
	}
}

func main() {
	// colorRed := "\033[31m"
	var result []FileInfoExt
	var resultAll []FileInfoExt
	opts := parseOptions()
	var waitDirs []string

	waitDirs = append(waitDirs, opts.Dir)
	round := 0
	for {
		round += 1
		if round >= 10000 {
			fmt.Println("too many directory, quit")
			os.Exit(1)
		}
		if len(waitDirs) >= 1 {
			thisDir := waitDirs[0]
			if len(waitDirs) > 1 {
				waitDirs = waitDirs[1:]
			} else if len(waitDirs) == 1 {
				waitDirs = waitDirs[:0]
			}
			fs := walkDir(thisDir)
			for _, f := range fs {
				tpath := filepath.Join(thisDir, f.Name())
				t, _ := os.Stat(tpath)
				if t.IsDir() {
					waitDirs = append(waitDirs, tpath)
				} else {
					tt := FileInfoExt{t, tpath}
					resultAll = append(resultAll, tt)
				}
			}
		} else {
			break
		}
		if !opts.IsRecur {
			break
		}
	}
	if opts.Startswith != "" {
		for _, f := range resultAll {
			if matchFileStartswith(opts.Startswith, f.FInfo.Name()) {
				result = append(result, f)
			}
		}
	} else if opts.FileSuffix != "" {
		for _, f := range resultAll {
			if matchFileSuffix(opts.FileSuffix, f.FInfo.Name()) {
				result = append(result, f)
			}
		}
	} else if opts.Name != "" {
		for _, f := range resultAll {
			if matchFileName(opts.Name, f.FInfo.Name()) {
				result = append(result, f)
			}

		}
	} else {
		result = append(result, resultAll...)
	}

	if opts.SortBy != "" {
		fmt.Println(opts.SortBy)
		sortBy(&result, opts.SortBy)
	}
	end(result, opts.PrintFullPath, opts.Dir)
	
	if opts.ActionRemoveFile {
		promptString := "yes\n"
		if platform == "windows" {
			promptString = "yes\r\n"
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("[DANGEROUS] delete %d files, please confirm (yes/no):", len(result))
		text, _ := reader.ReadString('\n')
		if text == promptString {
			realRemove(result)
		} else {
			fmt.Println("keep files, quit")
		}
	}
}
