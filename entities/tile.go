package entities

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tile represents a single tile in the game.
type Tile struct {
	Index   int
	X, Y    int
	TileMap *ebiten.Image
}

// Draw draws the tile on the screen.
func (t *Tile) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.X), float64(t.Y))
	sx, sy := (t.Index%8)*32, (t.Index/8)*32
	subImage := t.TileMap.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}
