package amesh

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var mux = newMux()

func Serve(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/clouds", handleClouds)

	return mux
}

func colorDistance(a color.Color, b color.Color) float64 {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return math.Pow(float64(ar-br), 2) + math.Pow(float64(ag-bg), 2) + math.Pow(float64(ab-bb), 2) + math.Pow(float64(aa-ba), 2)
}

var hexRegex = regexp.MustCompile("^#(?:[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$")

const hexFormatFullRGB = "#%02x%02x%02x"
const hexFormatFullRGBA = "#%02x%02x%02x%02x"

func hexToColor(s string) color.Color {
	s = strings.ToLower(s)

	if !hexRegex.MatchString(s) {
		return nil
	}

	var r, g, b uint8
	a := uint8(255)
	if len(s) == 4 {
		fmt.Sscanf(s, hexFormatFullRGB, &r, &g, &b)
	} else {
		fmt.Sscanf(s, hexFormatFullRGBA, &r, &g, &b, &a)
	}
	return color.RGBA{r, g, b, a}
}

func colorToHex(c color.Color) string {
	r, g, b, a := c.RGBA()
	if a == 255 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
}

var originalColorMap = []color.Color{
	hexToColor("#00000000"), // 晴れ
	hexToColor("#68349A"),   // 猛烈な雨
	hexToColor("#EA3323"),   // 非常に激しい雨
	hexToColor("#EA33F7"),   // やや強い雨
	hexToColor("#110AF5"),   // やや強い雨
	hexToColor("#FFFFFF"),   // より弱い雨
}

var newColorMap = []color.Color{
	hexToColor("#00000000"), // 晴れ
	hexToColor("#FF0000"),   // 猛烈な雨
	hexToColor("#FFFF00"),   // 非常に激しい雨
	hexToColor("#FF00FF"),   // やや強い雨
	hexToColor("#FF0000"),   // やや強い雨
	hexToColor("#FFFFFF"),   // より弱い雨
}

func handleClouds(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("/Users/keishi/amesh-proxy/202002271350.gif")
	if err != nil {
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	mrgba, ok := m.(*image.RGBA)
	if !ok {
		mrgba = image.NewRGBA(m.Bounds())
		draw.Draw(mrgba, m.Bounds(), m, image.Point{}, draw.Src)
	}

	bounds := m.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := m.At(x, y)
			st := 0
			md := math.Inf(1)
			for i, c := range originalColorMap {
				d := colorDistance(p, c)
				if d < md {
					st = i
					md = d
				}
			}
			mrgba.Set(x, y, newColorMap[st])
		}
	}

	of, err := os.Create("/tmp/image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(of, mrgba); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := of.Close(); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte("hello from one"))
}
