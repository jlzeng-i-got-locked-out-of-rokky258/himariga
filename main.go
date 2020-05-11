package main

import (
	_ "bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "github.com/hajimehoshi/ebiten/examples/resources/images"

)
const (
	screenWidth  = 650
	screenHeight = 425
	maxAngle     = 256
)
var (
	ebitenImage *ebiten.Image
)

type Sprite struct {
	imageWidth  int
	imageHeight int
	x           int
	y           int
	vx          int
	vy          int
	angle       int
	life		int
}

func (s *Sprite) Update() {
	//reinitialize a new himariga
	xpos, ypos := ebiten.CursorPosition()
	if s.life == 0{
		s.x, s.y = xpos - 33,ypos - 33
		s.vx, s.vy = 10*rand.Intn(2)-1, 5*rand.Intn(2)-1
		s.angle = rand.Intn(maxAngle)
		s.life = 1

	}
	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if mx := screenWidth - s.imageWidth; mx <= s.x {
		s.x = 2*mx - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if my := screenHeight - s.imageHeight; my <= s.y {
		s.y = 2*my - s.y
		s.vy = -s.vy
	}
	s.angle++
	if s.angle == maxAngle {
		s.angle = 0
	}
}

type Sprites struct {
	sprites []*Sprite
	num     int
}

func (s *Sprites) Update() {
	for i := 0; i < s.num; i++ {
		s.sprites[i].Update()
	}
}

const (
	MinSprites = 0
	MaxSprites = 50000
)

var (
	sprites = &Sprites{make([]*Sprite, MaxSprites), 0}
	op      = &ebiten.DrawImageOptions{}
)

func init() {
	// Decode image from a byte slice instead of a file so that
	// this example works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	existingImageFile, err := os.Open("Assets/Sprites/himariga.png")
	if err != nil {
		// Handle error
	}
	defer existingImageFile.Close()
	img, _, err := image.Decode(existingImageFile)
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	w, h := origEbitenImage.Size()
	ebitenImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 1)
	ebitenImage.DrawImage(origEbitenImage, op)

	for i := range sprites.sprites {
		w, h := ebitenImage.Size()
		x, y := rand.Intn(screenWidth-w), rand.Intn(screenHeight-h)
		vx, vy := 2*rand.Intn(2)-1, 2*rand.Intn(2)-1
		a := rand.Intn(maxAngle)
		sprites.sprites[i] = &Sprite{
			imageWidth:  w,
			imageHeight: h,
			x:           x,
			y:           y,
			vx:          vx,
			vy:          vy,
			angle:       a,
			life:		0,
		}
	}
}

func update(screen *ebiten.Image) error {
	// Decrease the nubmer of the sprites by hovering over them.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight){

		DeleteSprites(sprites)
	}
	// Increase the number of the sprites.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		sprites.num += 1
		if MaxSprites < sprites.num {
			sprites.num = MaxSprites
		}
	}

	sprites.Update()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Draw each sprite.
	// DrawImage can be called many many times, but in the implementation,
	// the actual draw call to GPU is very few since these calls satisfy
	// some conditions e.g. all the rendering sources and targets are same.
	// For more detail, see:
	// https://pkg.go.dev/github.com/hajimehoshi/ebiten#Image.DrawImage
	w, h := ebitenImage.Size()
	for i := 0; i < sprites.num; i++ {
		s := sprites.sprites[i]
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / maxAngle)
		op.GeoM.Translate(float64(w)/2, float64(h)/2)
		op.GeoM.Translate(float64(s.x), float64(s.y))
		screen.DrawImage(ebitenImage, op)
	}
	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Num of sprites: %d
Press left mouse to Himariga
right mouse to anti-Himariga`, ebiten.CurrentTPS(), ebiten.CurrentFPS(), sprites.num)
	ebitenutil.DebugPrint(screen, msg)
	return nil
}

//The worst method ever created honestly
func DeleteSprites(s *Sprites) {
	xpos, ypos := ebiten.CursorPosition()
	for i := 0; i < sprites.num; i++ {
		if math.Abs(float64(xpos-55-s.sprites[i].x)) < 30 && math.Abs(float64(ypos-55-s.sprites[i].y)) < 30{
			s.sprites[i].life = 0
			temp := s.sprites[i]
			s.sprites[i] = s.sprites[sprites.num - 1]
			s.sprites[sprites.num - 1] = temp
			sprites.num--
			i--
			if sprites.num < MinSprites {
				sprites.num = MinSprites
			}
		}}
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Himariga"); err != nil {
		log.Fatal(err)
	}
}