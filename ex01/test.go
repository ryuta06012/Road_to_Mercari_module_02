/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   test.go                                            :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: hryuuta <hryuuta@student.42tokyo.jp>       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2022/01/05 16:11:04 by hryuuta           #+#    #+#             */
/*   Updated: 2022/01/05 16:37:18 by hryuuta          ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "https://4.bp.blogspot.com/-2-Ny23XgrF0/Ws69gszw2jI/AAAAAAABLdU/unbzWD_U8foWBwPKWQdGP1vEDoQoYjgZwCLcBGAs/s1600/top_banner.jpg"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header.Get("Authorization"))
	resp.Header.Set("Authorization", "Bearer access-token")
	fmt.Println(resp.Header.Get("Authorization"))
	resp.Header.Set("Authorization", "BBBBBB access-token")
	fmt.Println(resp.Header.Get("Authorization"))
	//	fmt.Println(resp)

}
