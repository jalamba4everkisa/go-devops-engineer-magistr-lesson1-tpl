package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	var checkbox = 0
	response, err := http.Get("http://srv.msk01.gigacorp.local")
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode != 200 {
		checkbox += 1
	}
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
			bodyInt = append(bodyInt, 0)
			checkbox += 1
		} else {
			bodyInt = append(bodyInt, j)
		}
	}
	if checkbox > 3 {
		fmt.Println("Unable to fetch server statistic")
	} else {
		fmt.Println(memoryUsage(bodyInt[2], bodyInt[1]))
		fmt.Println(diskUsage(bodyInt[4], bodyInt[3]))
		fmt.Println(loadAverage(bodyInt[0]))
		fmt.Println(networkUsage(bodyInt[6], bodyInt[5]))
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
	p := float64(part)
	t := float64(total)
	temp := p / t * 100
	result := fmt.Sprintf("%.f", math.RoundToEven(temp))
	if temp > 80.0 {
		result = "Memory usage too high: " + result
	}
	return result
}

func diskUsage(part, total int) string {
	p := float64(part)
	t := float64(total)
	temp := p / t * 100
	result := fmt.Sprintf("%.f", math.RoundToEven(temp))
	if temp > 90.0 {
		space := (t - p) / 1024
		result = "Free disk space is too low: " + fmt.Sprintf("%.f", space) + " Mb left"
	}
	return result
}

func networkUsage(part, total int) string {
	p := float64(part) * 8
	t := float64(total) * 8
	temp := p / t * 100
	result := fmt.Sprintf("%.f", math.RoundToEven(temp))
	if temp > 90.0 {
		space := (t - p) / 1024
		result = "Network bandwidth usage high: " + fmt.Sprintf("%.f", space) + " Mbit/s available"
	}
	return result
}
