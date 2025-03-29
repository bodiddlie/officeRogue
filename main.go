package main

import (
	"image/color"
	"log"
	"officeRogue/entities"
	"officeRogue/game"
	"officeRogue/gamemap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Entities []game.Drawable
	Player   *entities.Sprite
	GameMap  *gamemap.GameMap
}

func (g *Game) Update() error {
	// Update player position based on keyboard input
	if g.Player != nil {
		g.Player.Update(g.GameMap)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	g.GameMap.Draw(screen)

	for _, entity := range g.Entities {
		entity.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	screenWidth, screenHeight := 640, 480

	gameMap, err := gamemap.LoadMap("sample_map.json")
	if err != nil {
		log.Fatal(err)
	}

	spriteSheet, _, err := ebitenutil.NewImageFromFile("assets/sprites/spritesheet_characters.png")
	if err != nil {
		log.Fatal(err)
	}

	// Create player sprite and position it in the center of the screen
	sprite := &entities.Sprite{
		X:           screenWidth/2 - 16,  // Center horizontally (32x32 sprite, so -16 for center)
		Y:           screenHeight/2 - 16, // Center vertically
		SpriteSheet: spriteSheet,
		Speed:       3,
		Width:       32,
		Height:      32,
	}
	drawables := []game.Drawable{sprite}

	// Create game instance with entities, player reference, and game map
	game := &Game{
		Entities: drawables,
		Player:   sprite,
		GameMap:  gameMap,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Office Rogue")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
