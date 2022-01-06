/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   test.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2022/01/05 16:11:04 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/06 16:19:22 by hryuuta          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"io"
	"os"
)

func main() {
	srcName := os.Args[1]
	src2Name := os.Args[2]
	dstName := os.Args[3]

	src1, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer src1.Close()

	src2, err := os.Open(src2Name)
	if err != nil {
		panic(err)
	}
	defer src2.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src1)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(dst, src2)
	if err != nil {
		panic(err)
	}
}
