package main

import (
	"bufio"
	"fmt"
	"golf/core"
	"os"
	"path/filepath"
)

var platform string = core.CheckOS()

func realRemove(files []core.FileInfoExt) {
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
	var result []core.FileInfoExt
	var resultAll []core.FileInfoExt
	opts := core.ParseOptions()
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
			fs, err := core.WalkDir(thisDir)
			if err != nil {
				// fmt.Println(err.Error(), thisDir)
				continue
			}
			for _, f := range fs {

				tpath := filepath.Join(thisDir, f.Name())
				t, err := os.Stat(tpath)
				if err != nil {
					// fmt.Printf("%v -> %+v\n", f, err.Error())
					continue
				}
				if t.IsDir() {
					waitDirs = append(waitDirs, tpath)
				} else {
					tt := core.FileInfoExt{t, tpath}
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
			if core.MatchFileStartswith(opts.Startswith, f.FInfo.Name()) {
				result = append(result, f)
			}
		}
	} else if opts.FileSuffix != "" {
		for _, f := range resultAll {
			if core.MatchFileSuffix(opts.FileSuffix, f.FInfo.Name()) {
				result = append(result, f)
			}
		}
	} else if opts.Name != "" {
		for _, f := range resultAll {
			if core.MatchFileName(opts.Name, f.FInfo.Name()) {
				result = append(result, f)
			}

		}
	} else {
		result = append(result, resultAll...)
	}

	var finalResult []core.FileInfoExt
	if opts.SizeFilter != "" {
		op, bytesSum, err := core.SizeFilter(opts.SizeFilter)
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, f := range result {
			if op == ">" {
				if f.FInfo.Size() > bytesSum {
					finalResult = append(finalResult, f)
				}
			} else if op == "=" {
				if f.FInfo.Size() == bytesSum {
					finalResult = append(finalResult, f)
				}
			} else if op == "<" {
				if f.FInfo.Size() < bytesSum {
					finalResult = append(finalResult, f)
				}
			}

		}
	} else {
		finalResult = result
	}

	if opts.SortBy != "" {
		core.SortBy(&finalResult, opts.SortBy)
	}
	core.InTheEnd(finalResult, opts.PrintFullPath, opts.Dir)

	if opts.ActionRemoveFile {
		promptString := "yes\n"
		if platform == "windows" {
			promptString = "yes\r\n"
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\n[DANGEROUS] delete %d files, please confirm (yes/no):", len(finalResult))
		text, _ := reader.ReadString('\n')
		if text == promptString {
			realRemove(finalResult)
		} else {
			fmt.Println("keep files, quit")
		}
	}
}
