package gamemap

import (
	"encoding/json"
	"image"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameMap represents the structure of the map in the JSON file.
type GameMap struct {
	Name       string `json:"name"`
	TileWidth  int    `json:"tileWidth"`
	TileHeight int    `json:"tileHeight"`
	Tileset    string `json:"tileset"`
	Layers     []struct {
		Name   string `json:"name"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Data   []int  `json:"data"`
	} `json:"layers"`
	TileImg *ebiten.Image
}

// LoadMap loads the map from a JSON file.
func LoadMap(filePath string) (*GameMap, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var gameMap GameMap
	if err := json.Unmarshal(bytes, &gameMap); err != nil {
		return nil, err
	}

	tileImg, _, err := ebitenutil.NewImageFromFile("assets/tiles/" + gameMap.Tileset)
	if err != nil {
		return nil, err
	}

	gameMap.TileImg = tileImg

	return &gameMap, nil
}

func (m *GameMap) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(m.TileImg, op)

	for _, layer := range m.Layers {
		for y := 0; y < layer.Height; y++ {
			for x := 0; x < layer.Width; x++ {
				idx := layer.Data[y*layer.Width+x]
				sx, sy := (idx%8)*32, (idx/8)*32
				subImage := m.TileImg.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image)
				op.GeoM.Reset()
				op.GeoM.Translate(float64(x*m.TileWidth), float64(y*m.TileHeight))
				screen.DrawImage(subImage, op)
			}
		}
	}
}

// IsWall checks if the tile at the given position is a wall (tile index 0)
func (m *GameMap) IsWall(x, y int) bool {
	// Convert pixel coordinates to tile coordinates
	tileX := x / m.TileWidth
	tileY := y / m.TileHeight

	// Check boundaries
	if tileX < 0 || tileY < 0 {
		return true // Out of bounds is considered a wall
	}

	// Check each layer (we're assuming the first layer contains walls)
	if len(m.Layers) > 0 {
		layer := m.Layers[0]
		if tileX >= layer.Width || tileY >= layer.Height {
			return true // Out of bounds is considered a wall
		}

		tileIndex := tileY*layer.Width + tileX
		if tileIndex >= 0 && tileIndex < len(layer.Data) {
			// Tile index 0 is a wall
			return layer.Data[tileIndex] == 0
		}
	}

	return false
}

// CanMoveToPosition checks if the entity can move to the given position
func (m *GameMap) CanMoveToPosition(x, y, width, height int) bool {
	// Check all four corners of the sprite
	return !m.IsWall(x, y) && // Top-left
		!m.IsWall(x+width-1, y) && // Top-right
		!m.IsWall(x, y+height-1) && // Bottom-left
		!m.IsWall(x+width-1, y+height-1) // Bottom-right
}
