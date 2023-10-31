package main

import (
	"fmt"

	"github.com/gookit/color"
)

func main() {
	result := color.Red.Sprint("Simple to use color")
	result += color.LightRed.Sprint("Light Red to use color")
	result += color.BgBlack.Sprint("Light Red to use color")
	fmt.Println(result)

	c := color.New(color.FgWhite, color.BgLightBlue, color.Bold)
	result += c.Sprint("Light Red to use color")
	fmt.Println(result)
}
