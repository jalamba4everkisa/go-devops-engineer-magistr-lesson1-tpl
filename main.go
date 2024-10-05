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
		check(memoryUsage(bodyInt[2], bodyInt[1]))
		check(diskUsage(bodyInt[4], bodyInt[3]))
		check(loadAverage(bodyInt[0]))
		check(networkUsage(bodyInt[6], bodyInt[5]))
	}
}

func loadAverage(la int) string {
	var result = strconv.Itoa(la)
	if la > 30 {
		result = "Load Average is too high: " + result
	} else {
		result = ""
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
	} else {
		result = ""
	}
	return result
}

func diskUsage(part, total int) string {
	result := ""
	p := float64(part)
	t := float64(total)
	temp := p / t * 100
	//result := fmt.Sprintf("%.f", math.RoundToEven(temp))
	if temp > 90.0 {
		diff := (t - p) / 1000000
		result = "Free disk space is too low: " + fmt.Sprintf("%.f", diff) + " Mb left"
	} else {
		result = ""
	}
	return result
}

func networkUsage(part, total int) string {
	result := ""
	p := float64(part)
	t := float64(total)
	temp := p / t * 100
	//result := fmt.Sprintf("%.f", math.RoundToEven(temp))
	if temp > 90.0 {
		diff := (t - p) * 0.000008
		result = "Network bandwidth usage high: " + fmt.Sprintf("%.f", diff) + " Mbit/s available"
	} else {
		result = ""
	}
	return result
}

func check(result string) {
	if result != "" {
		fmt.Println(result)
	}
}
