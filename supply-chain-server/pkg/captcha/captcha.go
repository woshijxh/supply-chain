package captcha

import (
	"math/rand"
	"sync"
	"time"
)

type Captcha struct {
	ID      string
	Content string // 计算结果
	Expr    string // 显示的表达式
}

var (
	captchas = make(map[string]string)
	mu       sync.RWMutex
)

func Generate() *Captcha {
	id := randomString(16)
	expr, result := generateMathExpr()

	mu.Lock()
	captchas[id] = result
	mu.Unlock()

	return &Captcha{
		ID:      id,
		Content: result,
		Expr:    expr,
	}
}

func GetByID(id string) *Captcha {
	mu.RLock()
	content, ok := captchas[id]
	mu.RUnlock()

	if !ok {
		return nil
	}

	return &Captcha{
		ID:      id,
		Content: content,
	}
}

func Verify(id, code string) bool {
	mu.RLock()
	expected, ok := captchas[id]
	mu.RUnlock()

	if !ok {
		return false
	}

	mu.Lock()
	delete(captchas, id)
	mu.Unlock()

	return expected == code
}

// generateMathExpr 生成数学计算表达式
func generateMathExpr() (expr string, result string) {
	// 两个数字的范围
	a := rand.Intn(50) + 1  // 1-50
	b := rand.Intn(50) + 1  // 1-50

	op := rand.Intn(3) // 0: +, 1: -, 2: *

	var res int
	var opChar string

	switch op {
	case 0: // 加法
		res = a + b
		opChar = "+"
	case 1: // 减法，确保结果为正
		if a < b {
			a, b = b, a
		}
		res = a - b
		opChar = "-"
	case 2: // 乘法，数字小一点
		a = rand.Intn(9) + 1  // 1-9
		b = rand.Intn(9) + 1  // 1-9
		res = a * b
		opChar = "×"
	}

	return formatNum(a) + " " + opChar + " " + formatNum(b) + " = ?", intToStr(res)
}

func formatNum(n int) string {
	return intToStr(n)
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	var negative bool
	if n < 0 {
		negative = true
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}