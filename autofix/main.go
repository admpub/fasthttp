package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Replace struct {
	Old    string
	New    string
	Regexp *regexp.Regexp
}

var replaces = []*Replace{
	&Replace{`SkipBody bool`, `SkipBody bool

	// Response.Write() skips writing response (both header and body) if set
	// to true. Use it to take complete control over what's sent after
	// hijacking a RequestCtx.
	SkipResponse bool`, nil},

	&Replace{`dst.SkipBody = resp.SkipBody`, `dst.SkipBody = resp.SkipBody
	dst.SkipResponse = resp.SkipResponse`, nil},

	&Replace{`resp.SkipBody = false`, `resp.SkipBody = false
	resp.SkipResponse = false`, nil},

	&Replace{`func (resp *Response) Write(w *bufio.Writer) error {`, `func (resp *Response) Write(w *bufio.Writer) error {
	if resp.SkipResponse {
		return nil
	}`, nil},
}

func main() {
	root := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/fasthttp`)
	/*
		save := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/fasthttp`)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if info.Name() == `_tools` || strings.HasPrefix(info.Name(), `.`) {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasPrefix(info.Name(), `.`) {
				return nil
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(b)
			for _, re := range replaces {
				if re.Regexp == nil {
					content = strings.Replace(content, re.Old, re.New, -1)
				} else {
					content = re.Regexp.ReplaceAllString(content, re.New)
				}
			}
			saveAs := strings.TrimPrefix(path, root)
			saveAs = filepath.Join(save, saveAs)
			err = os.MkdirAll(filepath.Dir(saveAs), os.ModePerm)
			if err == nil {
				file, err := os.Create(saveAs)
				if err == nil {
					_, err = file.WriteString(content)
				}
			}
			if err != nil {
				return err
			}
			fmt.Println(`Autofix ` + path + `.`)
			return nil
		})
	*/
	path := filepath.Join(root, `http.go`)
	b, err := ioutil.ReadFile(path)
	if err == nil {
		content := string(b)
		for _, re := range replaces {
			if re.Regexp == nil {
				content = strings.Replace(content, re.Old, re.New, -1)
			} else {
				content = re.Regexp.ReplaceAllString(content, re.New)
			}
		}
		saveAs := path
		err = os.MkdirAll(filepath.Dir(saveAs), os.ModePerm)
		if err == nil {
			file, err := os.Create(saveAs)
			if err == nil {
				_, err = file.WriteString(content)
			}
		}
	}
	defer time.Sleep(5 * time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(`Autofix complete.`)
}
