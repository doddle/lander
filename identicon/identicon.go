package identicon

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/llgcode/draw2d/draw2dimg"
	// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
	"github.com/icza/gox/imagex/colorx"
)

type point struct {
	x, y int
}

type drawingPoint struct {
	topLeft     point
	bottomRight point
}

type gridPoint struct {
	value byte
	index int
}

type Identicon struct {
	Name       string
	hash       [16]byte
	color      [3]byte
	grid       []byte
	gridPoints []gridPoint
	pixelMap   []drawingPoint
	width      int
	height     int
}

// WriteTo writes the identicon image to the given writer
func (i Identicon) WriteImage(w io.Writer) error {
	var img = image.NewRGBA(image.Rect(0, 0, i.width, i.height))
	col := color.RGBA{R: i.color[0], G: i.color[1], B: i.color[2], A: 255}

	for _, pixel := range i.pixelMap {
		rect(img, col, float64(pixel.topLeft.x), float64(pixel.topLeft.y), float64(pixel.bottomRight.x), float64(pixel.bottomRight.y))
	}

	return png.Encode(w, img)
}

type applyFunc func(Identicon) Identicon

// Generate is the main entrypoint to creating the identicon
func Generate(input []byte, hex string, width int, height int) Identicon {

	// orange #e88726
	// color, err := colorx.ParseHexColor("#e88726")
	color, err := colorx.ParseHexColor(hex)
	// blue 26a4e8
	// color, err := colorx.ParseHexColor("#26a4e8")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(color)

	pickedColour := [3]byte{color.R, color.G, color.B}

	identiconPipe := []applyFunc{
		buildGrid, filterOddSquares, buildPixelMap,
	}
	identicon := hashInput(input)

	// over-ride the settings
	identicon.color = pickedColour
	identicon.width = width
	identicon.height = height

	for _, applyFunc := range identiconPipe {
		identicon = applyFunc(identicon)
	}
	return identicon
}

func hashInput(input []byte) Identicon {
	checkSum := md5.Sum(input)
	return Identicon{
		Name: string(input),
		hash: checkSum,
	}
}

func buildGrid(identicon Identicon) Identicon {
	var grid []byte
	for i := 0; i < len(identicon.hash) && i+3 <= len(identicon.hash)-1; i += 3 {
		chunk := make([]byte, 5)
		copy(chunk, identicon.hash[i:i+3])
		chunk[3] = chunk[1]
		chunk[4] = chunk[0]
		grid = append(grid, chunk...)

	}
	identicon.grid = grid
	return identicon
}

func filterOddSquares(identicon Identicon) Identicon {
	var grid []gridPoint
	for i, code := range identicon.grid {
		if code%2 == 0 {
			point := gridPoint{
				value: code,
				index: i,
			}
			grid = append(grid, point)
		}
	}
	identicon.gridPoints = grid
	return identicon
}

func rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 float64) {
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(col)
	gc.MoveTo(x1, y1)
	gc.LineTo(x1, y1)
	gc.LineTo(x1, y2)
	gc.MoveTo(x2, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.SetLineWidth(0)
	gc.FillStroke()
}

func buildPixelMap(identicon Identicon) Identicon {
	var drawingPoints []drawingPoint

	pixelFunc := func(p gridPoint) drawingPoint {
		horizontal := (p.index % 5) * 50
		vertical := (p.index / 5) * 50
		topLeft := point{x: horizontal, y: vertical}
		bottomRight := point{x: horizontal + 50, y: vertical + 50}

		return drawingPoint{
			topLeft,
			bottomRight,
		}
	}

	for _, gridPoint := range identicon.gridPoints {
		drawingPoints = append(drawingPoints, pixelFunc(gridPoint))
	}
	identicon.pixelMap = drawingPoints
	return identicon
}
