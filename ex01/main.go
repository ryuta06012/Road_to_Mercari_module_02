/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2021/12/31 23:46:36 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/03 16:27:02 by hryuuta          ###   ########.fr       */
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
	fileName := path.Base(url)
	println(fileName)
	println(resp.Header.Get("Content-Length"))
	println(resp.Header)
	/* if resp.Header.Get("Accept-Ranges") != "bytes" {
		err = fmt.Errorf("Accept-Ranges = bytesではありません")
		return err
	} */
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
	if resp.Header.Get("Accept-Ranges") != "bytes" {
		err = fmt.Errorf("Accept-Ranges = bytesではありません")
		return 0, err
	}
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	return size, err
}

func sizeSplit(size int) (int, int) {
	var count int

	count = 5
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
		println("low =", low, "hight =", hight)
	}
	return low, hight
}

func SplitDownloadRun(fileUrl string) error {

	fullSize, err := sizeCheck(fileUrl)
	if err != nil {
		return err
	}

	println(fullSize)
	if err := DownloadFile("avatar.jpg", fileUrl); err != nil {
		panic(err)
	}
	return nil
}

func main() {

	/* fileUrl := "https://4.bp.blogspot.com/-2-Ny23XgrF0/Ws69gszw2jI/AAAAAAABLdU/unbzWD_U8foWBwPKWQdGP1vEDoQoYjgZwCLcBGAs/s1600/top_banner.jpg"
	if err := DownloadFile("avatar.jpg", fileUrl); err != nil {
		panic(err)
	} */
	sizeSplit(1897)
	//Run(fileUrl)
}
