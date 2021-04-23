package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type FileInfoExt struct {
	FInfo        os.FileInfo
	RelativePath string
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

func MatchFileName(rg string, fname string) bool {
	// r, _ := regexp.Compile(rg)
	b, err := regexp.MatchString(rg, fname)
	if err != nil {
		fmt.Println(err.Error(), "use go regex")
		os.Exit(1)
	}
	return b
}

func MatchFileSuffix(suffix string, fname string) bool {
	suffixs := strings.Split(suffix, "|")

	ss := strings.Split(fname, ".")
	s := ss[len(ss)-1]

	for _, sf := range suffixs {
		if strings.EqualFold(strings.ToLower(sf), strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func MatchFileStartswith(s string, fname string) bool {
	return strings.HasPrefix(fname, s)
}

func WalkDir(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	return files, err
}

func SizeFilter(s string) (string, int64, error) {
	var e error
	if !strings.HasPrefix(s, ">") && !strings.HasPrefix(s, "=") && !strings.HasPrefix(s, "<") {
		return "", 0, e
	}
	var t []rune
	lastDigit := 0
	for i, n := range s {
		if unicode.IsDigit(n) {
			t = append(t, rune(n))
			lastDigit = i
		}
	}
	extra := strings.TrimSpace(s[lastDigit+1:])
	sizeString := string(t) + " " + extra
	n, _ := ParseBytes(sizeString)
	// n, err:= strconv.Atoi(string(t))
	return string(s[0]), int64(n), nil
}

func ReadFileInfoTime(finfo os.FileInfo) string {
	return fmt.Sprintf("%v", finfo.ModTime().Format("2006-01-02 15:04:05"))
}

func SortBy(files *[]FileInfoExt, key string) {
	if key == "size" {
		sort.Slice(*files, func(i, j int) bool {
			return (*files)[i].FInfo.Size() > (*files)[j].FInfo.Size()
		})
	}
}

func InTheEnd(result []FileInfoExt, fullpath bool, pdir string) {
	for _, f := range result {
		if fullpath {
			// absPath, _ := filepath.Abs(f.FInfo.Name())
			// fmt.Printf("%s\n", absPath)
			fmt.Printf("[+] %v - %v - %s\n", ReadFileInfoTime(f.FInfo), f.FInfo.Size(), f.RelativePath)
		} else {
			fmt.Printf("[+] %v - %s - %s\n", ReadFileInfoTime(f.FInfo), Bytes(uint64(f.FInfo.Size())), f.RelativePath)
		}
	}
	fmt.Printf("[*] count: %d\n", len(result))
}
