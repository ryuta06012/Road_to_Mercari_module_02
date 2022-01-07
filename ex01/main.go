/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2021/12/31 23:46:36 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/07 22:53:50 by hryuuta          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Option struct {
	TargetURL string // ダウンロードの対象URL
	PCount    int    // 分割数
	OutputDir string // 結合後のファイルの格納場所
}

type Chunk struct {
	low      int
	hight    int
	filename string
	dirname  string
}

func (p Chunk) Init() {
	p.low = 0
	p.hight = 0
	p.filename = ""
	p.dirname = ""
}

func (o *Option) Init() {
	o.TargetURL = ""
	o.PCount = 1000
	o.OutputDir = "./new/"
}

func DownloadFile(filepath, url string, low, hight int) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	//defer cancel()
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", low, hight))
	res, err := http.DefaultClient.Do(resp)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	return err
}

func sizeCheck(url string) (int, error) {

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		err = fmt.Errorf("Accept-Ranges = bytesではありません")
		return 0, err
	}
	println(resp.Header.Get("Content-Length"))
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	return size, err
}

func sizeSplit(size int, fileName string, option Option) []Chunk {
	parts := make([]Chunk, option.PCount)
	var low, hight int
	for i := 0; i < option.PCount; i++ {
		if i == 0 {
			low = 0
		} else {
			low = hight + 1
		}
		if i == option.PCount-1 {
			hight = size - 1
		} else {
			hight = (size * (i + 1) / option.PCount)
		}
		fn := fileName + "_" + strconv.Itoa(i)
		part := Chunk{low: low, hight: hight, filename: fn}
		parts[i] = part
	}
	return parts
}

func mkDirTmp(dirname string) error {
	return os.Mkdir(dirname, 0777)
}

func rmDirTmp(dirname string) error {
	return os.Remove(dirname)
}

func getNewFilePath(dirname, filename string) string {
	return dirname + filename
}

func (p Chunk) getFilePath() string {
	return "./tmp/" + p.filename
}

func merge(parts []Chunk, filename string) error {
	newfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	for _, part := range parts {
		pf, err := os.Open(part.getFilePath())
		if err != nil {
			return err
		}
		_, err = io.Copy(newfile, pf)
		if err != nil {
			return err
		}
		defer pf.Close()
	}
	defer newfile.Close()
	return nil
}

func removePartFile(parts []Chunk) error {
	for _, part := range parts {
		if err := os.Remove(part.getFilePath()); err != nil {
			return err
		}
	}
	return nil
}

func SplitDownloadRun(fileUrl string, option Option) error {

	if err := mkDirTmp("tmp"); err != nil {
		return err
	}
	fullSize, err := sizeCheck(fileUrl)
	if err != nil {
		return err
	}
	fileName := path.Base(fileUrl)
	parts := sizeSplit(fullSize, fileName, option)
	for _, v := range parts {
		p := v
		fmt.Println(p)
		if err := DownloadFile(v.getFilePath(), fileUrl, p.low, p.hight); err != nil {
			rmDirTmp("tmp")
			panic(err)
		}
	}
	if err := mkDirTmp("new"); err != nil {
		return err
	}
	if err := merge(parts, getNewFilePath(option.OutputDir, fileName)); err != nil {
		return err
	}
	if err := removePartFile(parts); err != nil {
		return err
	}
	if err := rmDirTmp("tmp"); err != nil {
		return err
	}
	return nil
}

func main() {

	/* fileUrl := "https://4.bp.blogspot.com/-2-Ny23XgrF0/Ws69gszw2jI/AAAAAAABLdU/unbzWD_U8foWBwPKWQdGP1vEDoQoYjgZwCLcBGAs/s1600/top_banner.jpg"
	if err := DownloadFile("avatar.jpg", fileUrl); err != nil {
		panic(err)
	} */
	//mkDirTmp()
	/* fileName := path.Base(fileUrl)
	parts := sizeSplit(2000, fileName)
	println(parts)
	for _, v := range parts {
		//part := v
		fmt.Println(v)
	} */
	//fmt.Println(sizeSplit(1897, "avatar.jpg"))
	//Run(fileUrl)
	option := Option{}
	option.Init()
	//fmt.Println("aaaaa", os.Args[0])
	//option.TargetURL = os.Args[1]
	if err := SplitDownloadRun("https://releases.ubuntu.com/20.04/ubuntu-20.04.3-live-server-amd64.iso", option); err != nil {
		log.Fatal(err)
		return
	}
}
