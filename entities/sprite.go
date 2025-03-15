package entities

import (
	"image"
	"officeRogue/game"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite represents a character sprite in the game.
type Sprite struct {
	X, Y        int
	SpriteSheet *ebiten.Image
	Speed       int
	Width       int
	Height      int
}

// Draw draws the sprite on the screen.
func (s *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	subImage := s.SpriteSheet.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

// Update handles the sprite movement based on keyboard input
func (s *Sprite) Update(collidable game.Collidable) {
	// Default speed if not set
	if s.Speed == 0 {
		s.Speed = 3
	}

	// Default size if not set
	if s.Width == 0 {
		s.Width = 32
	}
	if s.Height == 0 {
		s.Height = 32
	}

	// Calculate potential new positions
	newX, newY := s.X, s.Y

	// Handle arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		newY -= s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		newY += s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		newX -= s.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		newX += s.Speed
	}

	// Check if the new position is valid (no collision with walls)
	if collidable.CanMoveToPosition(newX, newY, s.Width, s.Height) {
		s.X, s.Y = newX, newY
	} else {
		// Try to slide along walls for smoother movement
		if newX != s.X && collidable.CanMoveToPosition(newX, s.Y, s.Width, s.Height) {
			s.X = newX // Can move horizontally
		} else if newY != s.Y && collidable.CanMoveToPosition(s.X, newY, s.Width, s.Height) {
			s.Y = newY // Can move vertically
		}
	}
}
