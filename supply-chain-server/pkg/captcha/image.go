package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
)

// 简化的数字和符号字体 (6x10) - 更大更清晰
var font6x10 = map[rune][]bool{
	'0': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
	},
	'1': {
		false, false, true, true, false, false,
		false, true, true, true, false, false,
		true, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		true, true, true, true, true, true,
	},
	'2': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		false, false, false, false, true, true,
		false, false, false, true, true, false,
		false, false, true, true, false, false,
		false, true, true, false, false, false,
		true, true, false, false, false, false,
		true, false, false, false, false, false,
		true, true, false, false, true, true,
		true, true, true, true, true, true,
	},
	'3': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		false, false, false, false, true, true,
		false, false, false, true, true, false,
		false, true, true, true, true, false,
		false, false, false, true, true, false,
		false, false, false, false, true, true,
		false, false, false, false, true, true,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
	},
	'4': {
		false, false, false, true, true, false,
		false, false, true, true, true, false,
		false, true, false, true, true, false,
		true, false, false, true, true, false,
		true, true, true, true, true, true,
		false, false, false, true, true, false,
		false, false, false, true, true, false,
		false, false, false, true, true, false,
		false, false, false, true, true, false,
		false, false, false, true, true, false,
	},
	'5': {
		true, true, true, true, true, true,
		true, true, false, false, false, false,
		true, false, false, false, false, false,
		true, true, true, true, true, false,
		false, false, false, false, true, true,
		false, false, false, false, false, true,
		false, false, false, false, false, true,
		false, false, false, false, true, true,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
	},
	'6': {
		false, false, true, true, true, true,
		false, true, true, false, false, false,
		true, true, false, false, false, false,
		true, false, false, false, false, false,
		true, true, true, true, true, false,
		true, true, false, false, true, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
	},
	'7': {
		true, true, true, true, true, true,
		true, true, true, true, true, true,
		false, false, false, true, true, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
	},
	'8': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
	},
	'9': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		true, false, false, false, false, true,
		true, false, false, false, false, true,
		true, true, false, false, true, true,
		false, true, true, true, true, true,
		false, false, false, false, true, true,
		false, false, false, false, true, true,
		false, false, false, true, true, false,
		true, true, true, true, false, false,
	},
	'+': {
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		true, true, true, true, true, true,
		true, true, true, true, true, true,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
		false, false, true, true, false, false,
	},
	'-': {
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		true, true, true, true, true, true,
		true, true, true, true, true, true,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
	},
	'×': {
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		true, true, false, false, true, true,
		false, true, true, true, true, false,
		false, false, true, true, false, false,
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
	},
	'=': {
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		true, true, true, true, true, true,
		true, true, true, true, true, true,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		true, true, true, true, true, true,
		true, true, true, true, true, true,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
	},
	'?': {
		false, true, true, true, true, false,
		true, true, false, false, true, true,
		false, false, false, false, true, true,
		false, false, false, true, true, false,
		false, false, true, true, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
		false, false, false, false, false, false,
		false, true, true, false, false, false,
		false, true, true, false, false, false,
	},
	' ': {
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, false, false, false, false,
	},
}

// GenerateImage 生成验证码图片，expr 是表达式（如 "12 + 5 = ?"）
func GenerateImage(expr string) ([]byte, error) {
	width := 180
	height := 60
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 白色背景
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// 绘制表达式 - 使用深色，确保清晰
	textColor := color.RGBA{20, 20, 20, 255} // 深黑色

	startX := 8
	for _, c := range expr {
		drawChar(img, startX, 8, c, textColor)
		// 根据字符宽度调整间距
		if c == ' ' {
			startX += 10
		} else if c == '+' || c == '-' || c == '×' || c == '=' {
			startX += 22
		} else {
			startX += 20
		}
	}

	// 添加少量干扰线（颜色浅，不影响阅读）
	for i := 0; i < 3; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		lineColor := color.RGBA{200, 200, 200, 255} // 浅灰色干扰线
		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func drawChar(img *image.RGBA, x, y int, char rune, col color.RGBA) {
	pixels, ok := font6x10[char]
	if !ok {
		return
	}

	if len(pixels) != 60 { // 6x10 = 60
		return
	}

	// 放大到 4x4 像素点
	scale := 4
	for i, on := range pixels {
		if on {
			px := x + (i%6)*scale
			py := y + (i/6)*scale
			for dx := 0; dx < scale; dx++ {
				for dy := 0; dy < scale; dy++ {
					img.Set(px+dx, py+dy, col)
				}
			}
		}
	}
}

// drawLine 绘制干扰线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}

	if steps == 0 {
		img.Set(x1, y1, col)
		return
	}

	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)

	x := float64(x1)
	y := float64(y1)
	for i := 0; i <= steps; i++ {
		img.Set(int(x), int(y), col)
		x += xInc
		y += yInc
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}