package core

import (
	"flag"
	"fmt"
	"os"
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

func ParseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.IsRecur, "recur", false, "是否递归")
	flag.StringVar(&options.Dir, "d", "", "扫描目录, 如果不提供则扫描当前工作目录")
	// flag.StringVar(&options.Type, "t", "", "-t=f or -t=d")
	flag.StringVar(&options.SortBy, "s", "", "时间time 文件大小size")
	flag.StringVar(&options.Name, "name", "", "-name=foo")
	flag.StringVar(&options.FileSuffix, "suffix", "", "-suffix=zip 注意 程序只取split('.')[-1]")
	flag.BoolVar(&options.PrintFullPath, "bytes", false, "size with bytes")
	flag.BoolVar(&options.ActionRemoveFile, "rm", false, "[!] 删除文件")
	flag.StringVar(&options.Startswith, "startswith", "", "startswith string")
	flag.StringVar(&options.SizeFilter, "size", "", "-size \">10k\"")
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
