package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
	red  color.Color = color.RGBA{255, 0, , 255}
	green  color.Color = color.RGBA{0, 255, , 255}
)

// ref) http://golang.org/doc/articles/image_draw.html
func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // Set Random Seed

	m := image.NewRGBA(image.Rect(0, 0, 100, 100)) //*NRGBA (image.Image interface)

	fmt.Printf("m %T\n", m)

	// White Background
	draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	DrawCross(m, rand.Intn(m.Bounds().Max.X-12)+12, rand.Intn(m.Bounds().Max.Y-12)+12)

	w, _ := os.Create("randCross.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	Show(w.Name())
}

//Create Black Cross
func DrawCross(m *image.RGBA, x, y int) (n int, err error) {
	if x < m.Bounds().Min.X-12 && y < m.Bounds().Min.Y-12 && x > m.Bounds().Max.X+12 && y > m.Bounds().Max.Y+12 {
		return 0, err
	}

	for i := -12; i < 13; i++ {
		for j := -2; j < 3; j++ {
			m.Set(x+i, y+j, black)
		}
	}
	for i := -12; i < 13; i++ {
		for j := -2; j < 3; j++ {
			m.Set(x+j, y+i, black)
		}
	}

	return 0, err
}

// show  a specified file by Preview.app for OS X(darwin)
func Show(name string) {
	command := "open"
	arg1 := "-a"
	arg2 := "/Applications/Preview.app"
	cmd := exec.Command(command, arg1, arg2, name)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
