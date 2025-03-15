package game

import "github.com/hajimehoshi/ebiten/v2"

// Drawable is an interface for anything that can be drawn on screen
type Drawable interface {
	Draw(screen *ebiten.Image)
}

// Collidable is an interface for map objects that can be checked for collision
type Collidable interface {
	CanMoveToPosition(x, y, width, height int) bool
}
