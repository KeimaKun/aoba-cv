package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"reflect"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
	red   color.Color = color.RGBA{255, 0, 0, 255}
)

// ref) http://golang.org/doc/articles/image_draw.html
func main() {
	file, err := os.Open("randCross.png")

	if err != nil {
		fmt.Println("file not found!")
		os.Exit(1)
	}

	fmt.Printf("file %T\n", file)

	img, err := png.Decode(file)

	fmt.Println(img.Bounds())

	img.At(0, 0).RGBA()
	fmt.Println("type:", reflect.TypeOf(img))

	fmt.Println(img.At(30, 40))
	fmt.Println(img.At(50, 60))

	//for i := 0; i < 5; i++ {
	//	for j := 0; j < 5; j++ {
	//		img.(*image.RGBA).Set(50+i, 60+j, red)
	//	}
	//}
	//fmt.Println(img.At(50, 60))

	DrawRSq(img.(*image.RGBA), 10, 20, 20, 10)

	w, _ := os.Create("findCross.png")
	defer w.Close()
	png.Encode(w, img) //Encode writes the Image m to w in PNG format.

	Show(w.Name())

	//---------------------------------------------------------//

	//m := image.NewRGBA(image.Rect(0, 0, 640, 480)) //*NRGBA (image.Image interface)

	// White Background
	//draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	//w, _ := os.Create("findCross.png")
	//defer w.Close()
	//png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	//Show(w.Name())
}

//Create Red Cross
func DrawRSq(m *image.RGBA, x, y, width, height int) (n int, err error) {
	//if x < m.Bounds().Min.X-12 && y < m.Bounds().Min.Y-12 && x > m.Bounds().Max.X+12 && y > m.Bounds().Max.Y+12 {
	//	return 0, err
	//}

	for i := 0; i < width; i++ {
		m.Set(x+i, y, red)
		m.Set(x+i, y+height, red)
	}
	for j := 0; j < height; j++ {
		m.Set(x, y+j, red)
		m.Set(x+width, y+j, red)
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
