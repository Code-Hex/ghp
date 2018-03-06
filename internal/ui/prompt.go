package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type choose []string

// Choose returns you've selected item as string
func Choose(msg string, list []string, d string) string {
	var result string
	for {
		fmt.Println(msg)
		var c choose
		for i, v := range list {
			fmt.Printf("%2d) %s\n", i+1, v)
			c = append(c, fmt.Sprintf("%d", i+1))
		}
		fmt.Printf("[%s]: ", d)

		input := readline()
		if input == "" {
			return d
		}

		if c.validator(input) {
			i, err := strconv.ParseInt(input, 10, 32)
			if err != nil {
				Errorf("Failed to parseint: %s\n", err.Error())
				continue
			}
			idx := int(i - 1)
			if !(0 <= idx && idx < len(list)) {
				Errorf("Index out of range")
				continue
			}
			result = list[idx]
			break
		}
		Errorf("Invalid input: %s\n", input)
	}

	return result
}

func readline() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return scanner.Err().Error()
}

func (c choose) validator(input string) bool {
	for _, v := range c {
		if input == v {
			return true
		}
	}
	return false
}
