package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/dariubs/percent"
)

func main() {
	var checkbox = 0
	response, err := http.Get("http://srv.msk01.gigacorp.local")
	if err != nil {
		fmt.Println(err)
		checkbox += 1
	}
	if response.StatusCode != 200 {
		checkbox += 1
	}
	body, err := io.ReadAll(response.Body)
	//defer
	response.Body.Close()
	if err != nil {
		checkbox += 1
		fmt.Println(err)
		return
	}
	bodyStrings := strings.Split(string(body), ",")
	bodyInt := []int{}
	for _, i := range bodyStrings {
		j, err := strconv.Atoi(strings.Trim(i, " "))
		if err != nil {
			bodyInt = append(bodyInt, 1)
			checkbox += 1
		} else {
			bodyInt = append(bodyInt, j)
		}
	}
	if checkbox > 2 {
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
	usagePercent := percent.PercentOf(part, total)
	result := fmt.Sprintf("%.f", math.RoundToEven(usagePercent))
	if usagePercent > 80 {
		result = "Memory usage too high: " + result + "%"
	} else {
		result = ""
	}
	return result
}

func diskUsage(part, total int) string {
	result := ""
	usagePercent := percent.PercentOf(part, total)
	if usagePercent > 90 {
		diff := float64(total-part) / 1048576
		result = "Free disk space is too low: " + fmt.Sprintf("%.f", diff) + " Mb left"
	} else {
		result = ""
	}
	return result
}

func networkUsage(part, total int) string {
	result := ""
	usagePercent := percent.PercentOf(part, total)
	if usagePercent > 90 {
		diff := float64(total-part) / 125000
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
