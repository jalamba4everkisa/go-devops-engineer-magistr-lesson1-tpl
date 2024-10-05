package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	response, err := http.Get("http://srv.msk01.gigacorp.local")
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		bodyStrings := strings.Split(string(body), ",")
		bodyInt := []int{}
		for _, i := range bodyStrings {
			j, err := strconv.Atoi(strings.Trim(i, " "))
			if err != nil {
				panic(err)
			}
			bodyInt = append(bodyInt, j)
		}
		fmt.Println(loadAverage(bodyInt[0]))
		fmt.Println(memoryUsage(bodyInt[2], bodyInt[1]))
		fmt.Println(diskUsage(bodyInt[4], bodyInt[3]))
		fmt.Println(networkUsage(bodyInt[6], bodyInt[5]))
	} else {
		fmt.Println("Unable to fetch server statistic")
	}
}

func loadAverage(la int) string {
	var result = strconv.Itoa(la)
	if la > 30 {
		result = "Load Average is too high: " + result
	}
	return result
}

func memoryUsage(part, total int) string {
	p := float32(part)
	t := float32(total)
	temp := p / t * 100
	result := fmt.Sprintf("%.3f", temp)
	if temp > 80.0 {
		result = "Memory usage too high: " + result
	}
	return result
}

func diskUsage(part, total int) string {
	p := float32(part)
	t := float32(total)
	temp := p / t * 100
	result := fmt.Sprintf("%.3f", temp)
	if temp > 90.0 {
		space := (t - p) / 1024
		result = "Free disk space is too low: " + fmt.Sprintf("%.3f", space) + " Mb left"
	}
	return result
}

func networkUsage(part, total int) string {
	p := float32(part) * 8
	t := float32(total) * 8
	temp := p / t * 100
	result := fmt.Sprintf("%.3f", temp)
	if temp > 90.0 {
		space := (t - p) / 1024
		result = "Network bandwidth usage high: " + fmt.Sprintf("%.3f", space) + " Mbit/s available"
	}
	return result
}
