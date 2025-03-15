package gamemap

import (
	"encoding/json"
	"image"
	"io"
	"os"

	"officeRogue/game"

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

	return &gameMap, nil
}

// CreateTiles creates a list of tiles from the map data.
func (m *GameMap) CreateTiles() ([]game.Drawable, error) {
	tileImg, _, err := ebitenutil.NewImageFromFile("assets/tiles/" + m.Tileset)
	if err != nil {
		return nil, err
	}

	var tiles []game.Drawable

	for _, layer := range m.Layers {
		for y := 0; y < layer.Height; y++ {
			for x := 0; x < layer.Width; x++ {
				idx := layer.Data[y*layer.Width+x]
				// Import tile from the entities package when needed
				tile := &Tile{
					Index:   idx,
					X:       x * m.TileWidth,
					Y:       y * m.TileHeight,
					TileMap: tileImg,
				}
				tiles = append(tiles, tile)
			}
		}
	}

	return tiles, nil
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
	return !m.IsWall(x, y) &&                     // Top-left
	       !m.IsWall(x+width-1, y) &&             // Top-right
	       !m.IsWall(x, y+height-1) &&            // Bottom-left
	       !m.IsWall(x+width-1, y+height-1)       // Bottom-right
}

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
