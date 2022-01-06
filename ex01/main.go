/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2021/12/31 23:46:36 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/06 16:36:17 by hryuuta          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Chunk struct {
	low      int
	hight    int
	filename string
}

func (p Chunk) Init() {
	p.low = 0
	p.hight = 0
	p.filename = ""
}

func DownloadFile(filepath string, url string, low, hight int) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	//defer cancel()
	println("cccc")
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	println("cccc")
	//defer resp.Body.Close()
	println("cccc")
	resp.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", low, hight))
	println("cccc")
	res, err := http.DefaultClient.Do(resp)
	println("cccc")
	out, err := os.Create(filepath)
	println("cccc")
	if err != nil {
		return err
	}
	println("cccc")
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

func sizeSplit(size int, fileName string) []Chunk {
	var count int

	count = 1000
	parts := make([]Chunk, count)
	var low, hight int
	for i := 0; i < count; i++ {
		if i == 0 {
			low = 0
		} else {
			low = hight + 1
		}
		if i == count-1 {
			hight = size - 1
		} else {
			hight = (size * (i + 1) / 1000)
		}
		fn := fileName + "_" + strconv.Itoa(i)
		part := Chunk{low: low, hight: hight, filename: fn}
		parts[i] = part
		//println("low =", low, "hight =", hight)
	}
	return parts
}

func mkDirTmp() error {
	return os.Mkdir("tmp", 0777)
}

func rmDirTmp() error {
	return os.Remove("tmp")
}

func (p Chunk) getFilePath() string {
	return "./tmp/" + p.filename
}

func merge(parts []Chunk, filename string) error {
	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}

	for _, v := range parts {
		pf, err := os.Open(v.getFilePath())
		if err != nil {
			return err
		}
		_, err = io.Copy(newfile, pf)
		defer pf.Close()
	}
	defer newfile.Close()
	return err
}

func SplitDownloadRun(fileUrl string) error {

	if err := mkDirTmp(); err != nil {
		return err
	}
	fullSize, err := sizeCheck(fileUrl)
	if err != nil {
		return err
	}
	fileName := path.Base(fileUrl)
	parts := sizeSplit(fullSize, fileName)
	for _, v := range parts {
		//fmt.Println("Split", v)
		p := v
		fmt.Println(p)
		if err := DownloadFile(v.getFilePath(), fileUrl, p.low, p.hight); err != nil {
			rmDirTmp()
			panic(err)
		}
	}

	/* if err := rmDirTmp(); err != nil {
		return err
	} */
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
	if err := SplitDownloadRun("https://releases.ubuntu.com/20.04/ubuntu-20.04.3-live-server-amd64.iso"); err != nil {
		println("error")
		return
	}
}
