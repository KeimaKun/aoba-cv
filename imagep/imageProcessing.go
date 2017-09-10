package imagep

import (
	"aoba-cv/colorscheme"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
)

var (
	Gx = [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	Gy = [3][3]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}
	//Gaussian blur constant, need norm by 159
	H = [5][5]int{{2, 4, 5, 4, 2}, {4, 9, 12, 9, 4}, {5, 12, 15, 12, 5}, {4, 9, 12, 9, 4}, {2, 4, 5, 4, 2}}
)

func Main() {
	infile, err := os.Open("test03.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer infile.Close()

	imgDecoded, _, err := image.Decode(infile)
	if err != nil {
		log.Fatalln(err)
	}

	img := &Converted{imgDecoded, color.RGBAModel}

	//Gx := [3][3]int{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}}
	//Gy := [3][3]int{{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}}

	bd := img.Bounds()

	//from original to gray
	imgGray := image.NewRGBA(bd)
	for y := 0; y < bd.Max.Y; y++ {
		for x := 0; x < bd.Max.X; x++ {
			//oldPixel := img.At(x, y)
			//r, g, b, _ := oldPixel.RGBA()
			//yGr := 0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256)
			//newPixel := color.Gray{uint8(yGr)}
			newPixel := color.Gray{img.GrayCode(x, y)}
			imgGray.Set(x, y, newPixel)
		}
	}

	// from original to blur
	imgGuassianBlur := image.NewRGBA(bd)
	for y := 0; y < bd.Max.Y; y++ {
		for x := 0; x < bd.Max.X; x++ {
			if x == 0 || y == 0 || x == 1 || y == 1 || x == bd.Max.X-1 || y == bd.Max.Y-1 || x == bd.Max.X || y == bd.Max.Y {
				imgGuassianBlur.Set(x, y, color.Gray{img.GrayCode(x, y)})
			} else {
				blurValue := float64(0)
				for i := 0; i < 5; i++ {
					for j := 0; j < 5; j++ {
						blurValue += float64(H[j][i]) * float64(img.GrayCode(x+i-2, y+j-2)) / float64(159)
					}
				}
				imgGuassianBlur.Set(x, y, color.Gray{uint8(blurValue)})
			}
		}
	}
	// from original to edge
	imgGradient := image.NewRGBA(bd)
	for y := 0; y < bd.Max.Y; y++ {
		for x := 0; x < bd.Max.X; x++ {
			if x == 0 || y == 0 || x == bd.Max.X || y == bd.Max.Y {
				imgGradient.Set(x, y, color.Gray{0})
			} else {
				imgGradient.Set(x, y, color.Gray{img.GradientTotal(x, y)})
			}
		}
	}

	//from blur to edge
	imgBlured := &Converted{imgGuassianBlur, color.RGBAModel}
	imgBlurGradient := image.NewRGBA(bd)
	for y := 0; y < bd.Max.Y; y++ {
		for x := 0; x < bd.Max.X; x++ {
			if x == 0 || y == 0 || x == bd.Max.X || y == bd.Max.Y {
				imgBlurGradient.Set(x, y, color.Gray{0})
			} else {
				//gradientX := 0
				//gradientY := 0
				//for i := 0; i < 3; i++ {
				//	for j := 0; j < 3; j++ {
				//		gradientX += Gx[j][i] * int(imgBlured.GrayCode(x+i-1, y+j-1))
				//		gradientY += Gy[j][i] * int(imgBlured.GrayCode(x+i-1, y+j-1))
				//	}
				//}
				//gradient := math.Sqrt(float64(gradientX)*float64(gradientX) + float64(gradientY)*float64(gradientY))
				imgBlurGradient.Set(x, y, color.Gray{imgBlured.GradientTotal(x, y)})
			}
		}
	}

	imgBlurGradientAngle := image.NewRGBA(bd)
	for y := 0; y < bd.Max.Y; y++ {
		for x := 0; x < bd.Max.X; x++ {
			if x == 0 || y == 0 || x == bd.Max.X || y == bd.Max.Y {
				imgBlurGradientAngle.Set(x, y, color.Gray{0})
			} else {
				imgBlurGradientAngle.Set(x, y, color.Gray{uint8(256.0 * imgBlured.GradientAngle(x, y) / 3.14)})
			}
		}
	}

	//fmt.Println(Gx[1][0])
	//fmt.Println(img.GrayCode(100, 100))
	//fmt.Println(int(img.GrayCode(100, 100)) * Gx[1][0])

	Show(infile.Name())

	outfileGray, err := os.Create("test03Gray.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfileGray.Close()
	png.Encode(outfileGray, imgGray) //Encode writes the Image m to w in PNG format.
	Show(outfileGray.Name())

	outfileBlur, err := os.Create("test03Blur.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfileBlur.Close()
	png.Encode(outfileBlur, imgGuassianBlur) //Encode writes the Image m to w in PNG format.
	Show(outfileBlur.Name())

	outfileGrad, err := os.Create("test03Gradient.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfileGrad.Close()
	png.Encode(outfileGrad, imgGradient) //Encode writes the Image m to w in PNG format.
	Show(outfileGrad.Name())

	outfileBlurGrad, err := os.Create("test03BlurGradient.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfileBlurGrad.Close()
	png.Encode(outfileBlurGrad, imgBlurGradient) //Encode writes the Image m to w in PNG format.
	Show(outfileBlurGrad.Name())

	outfileBlurGradAngle, err := os.Create("test03BlurGradientAngle.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfileBlurGradAngle.Close()
	png.Encode(outfileBlurGradAngle, imgBlurGradientAngle) //Encode writes the Image m to w in PNG format.
	Show(outfileBlurGradAngle.Name())

	fmt.Println(256.0 * imgBlured.GradientAngle(100, 110) / 3.14)

}

//Create Red Cross
func DrawRSq(m *image.RGBA, x, y, width, height int) (n int, err error) {
	//if x < m.Bounds().Min.X-12 && y < m.Bounds().Min.Y-12 && x > m.Bounds().Max.X+12 && y > m.Bounds().Max.Y+12 {
	//	return 0, err
	//}

	for i := 0; i < width; i++ {
		m.Set(x+i, y, colorscheme.Red)
		m.Set(x+i, y+height, colorscheme.Red)
	}
	for j := 0; j < height; j++ {
		m.Set(x, y+j, colorscheme.Red)
		m.Set(x+width, y+j, colorscheme.Red)
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

// Converted implements image.Image, so you can
// pretend that it is the converted image.
type Converted struct {
	Img image.Image
	Mod color.Model
}

// Module ColorModel
func (c *Converted) ColorModel() color.Model {
	return c.Mod
}

// ... but the original bounds
func (c *Converted) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image and
// then asks the color model to convert it.
func (c *Converted) At(x, y int) color.Color {
	return c.Mod.Convert(c.Img.At(x, y))
}

func (m *Converted) GrayCode(x, y int) uint8 {
	oldPixel := m.At(x, y)
	r, g, b, _ := oldPixel.RGBA()
	yGr := 0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256)
	return uint8(yGr)
	//newPixel := color.Gray{uint8(yGr)}
}

func (m *Converted) GradientX(x, y int) int {
	gradientX := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gradientX += Gx[j][i] * int(m.GrayCode(x+i-1, y+j-1))

		}
	}
	return gradientX
}

func (m *Converted) GradientY(x, y int) int {
	gradientY := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gradientY += Gy[j][i] * int(m.GrayCode(x+i-1, y+j-1))

		}
	}
	return gradientY
}

func (m *Converted) GradientTotal(x, y int) uint8 {
	return uint8(math.Sqrt(float64(m.GradientX(x, y))*float64(m.GradientX(x, y)) + float64(m.GradientY(x, y))*float64(m.GradientY(x, y))))
}

func (m *Converted) GradientAngle(x, y int) float64 {
	return math.Atan(float64(m.GradientY(x, y)) / float64(m.GradientX(x, y)))
}
