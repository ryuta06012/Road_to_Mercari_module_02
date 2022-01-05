/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2021/12/31 23:46:36 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/05 16:10:32 by hryuuta          ###   ########.fr       */
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

func DownloadFile(filepath string, url string) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	//defer cancel()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//fileName := path.Base(url)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func sizeCheck(url string) (int, error) {

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	/* if resp.Header.Get("Accept-Ranges") != "bytes" {
		err = fmt.Errorf("Accept-Ranges = bytesではありません")
		return 0, err
	} */
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	return size, err
}

func sizeSplit(size int, fileName string) []Chunk {
	var count int

	count = 5
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
			hight = (size * (i + 1) / 5)
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
	return "/tmp/" + p.filename
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
		fmt.Println("Split", v)
	}
	/* if err := DownloadFile("avatar.jpg", fileUrl); err != nil {
		panic(err)
	} */
	if err := rmDirTmp(); err != nil {
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
}
