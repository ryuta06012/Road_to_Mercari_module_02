/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   main.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2021/12/29 14:54:21 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/02 11:47:20 by hryuuta          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func makeCh() chan string {
	return make(chan string)
}

var t int

func init() {
	flag.IntVar(&t, "t", 30, "制限時間(分)")
	flag.Parse()
}

func input(r io.Reader) <-chan string {
	vh := makeCh()
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			vh <- s.Text()
		}
		close(vh)
	}()
	return vh
}

func getwords_from(txt_path string) ([]string, error) {
	var words []string
	sf, err := os.Open(txt_path)
	if err != nil {
		return nil, err
	} else {
		s := bufio.NewScanner(sf)
		for s.Scan() {
			words = append(words, s.Text())
		}
	}
	defer sf.Close()
	return words, err
}

func initContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	return ctx, cancel
}

func gameStart(ctx context.Context, words []string) int {
	var count int
	ch := input(os.Stdin)
	for i := 0; i < len(words); i++ {
		println(words[i])
		fmt.Print("->")
		select {
		case x, ok := <-ch:
			if !ok {
				return count
			}
			if x == words[i] {
				count++
				println("○")
			} else {
				println("x")
			}
		case <-ctx.Done():
			return count
		}
	}
	return count
}

func main() {
	ctx, cancel := initContext()
	defer cancel()
	words, err := getwords_from("text.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "No such file or directory", err)
		os.Exit(1)
	}
	println("Time Up! Score = ", gameStart(ctx, words))
}
