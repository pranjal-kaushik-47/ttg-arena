package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"tag-game-v2/common"
	"tag-game-v2/internal/entity"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TODO:
// 0. enemy AI (better logic ;P) : Done
// 1. map design system : Done

// #### add art and animation ####

// 2. hostile enemies
// 2.5 enemy from docile to hostile
// 3. room transition
// 4. fog of war
// 5. z axis movement
// 6. hiding

// ##### add music ####

// 7. ablities*
// 8. Boss fight

type Game struct {
	Player           *entity.Player
	Enemies          []*entity.Enemy
	Environment      *entity.Environment
	MetaData         common.GameMetaData
	HideSystemCursor bool

	// Temp Var for level editing
	BlockSize           int
	EnableBlockCollider bool
}

func (g *Game) Update() error {
	if g.HideSystemCursor {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
	x, y := ebiten.CursorPosition()
	// fmt.Println(x, y, inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft))
	CurrentEnemyCount := 0
	for _, enemy := range g.Enemies {
		if enemy.Sprite.IsActive {
			CurrentEnemyCount += 1
		}
	}
	if CurrentEnemyCount == 0 {
		fmt.Println("Level Completed ", g.MetaData.Level)
		g.MetaData.Level += 1
		g.MetaData.TotalEnemies = g.MetaData.TotalEnemies * 2
		g.MetaData.CurrentEnemyCount = g.MetaData.TotalEnemies
		g.NewLevel(g.MetaData)
	}
	g.Player.Update(&g.MetaData, g.Environment)
	for _, enemy := range g.Enemies {
		enemy.Update(&g.MetaData, g.Player, g.Environment)
	}

	// Edit Level
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.HideSystemCursor = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.HideSystemCursor = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.EnableBlockCollider = true
		fmt.Println("Enabled Collider")
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyQ) {
		g.EnableBlockCollider = false
		fmt.Println("Disabled Collider")
	}
	_, dy := ebiten.Wheel()
	g.BlockSize = g.BlockSize + int(dy)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Environment.BuildSquareWall(x, y, g.BlockSize, g.BlockSize, g.EnableBlockCollider)
	}

	if ebiten.IsKeyPressed(ebiten.KeyF1) {
		b, err := json.Marshal(g.Environment.Walls)
		if err != nil {
			panic(err)
		}
		os.WriteFile("internal/gamedata/levels/newlevel.json", b, 0644)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	clr := color.RGBA{135, 206, 235, 255} // RGB for light blue color
	screen.Fill(clr)
	if g.HideSystemCursor {
		x, y := ebiten.CursorPosition()
		rect := ebiten.NewImage(g.BlockSize, g.BlockSize)
		rect.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle

		// Draw the rectangle at the mouse cursor position
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x)-(float64(g.BlockSize)/2), float64(y)-(float64(g.BlockSize)/2))

		screen.DrawImage(rect, op)
	}
	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}
	g.Environment.Draw(screen, g.MetaData)
	g.Player.Draw(screen)

	// level details

	msg := fmt.Sprintf("Level: %d", g.MetaData.Level)
	op := &text.DrawOptions{}
	op.GeoM.Translate(24, 20)
	op.ColorScale.ScaleWithColor(color.White)

	s, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))

	text.Draw(screen, msg, &text.GoTextFace{
		Source: s,
		Size:   10,
	}, op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.MetaData.ScreenWidth, g.MetaData.ScreenHeight
}

func (g *Game) NewLevel(metaData common.GameMetaData) error {
	g.MetaData = metaData

	// environment
	g.Environment = &entity.Environment{}
	// g.Environment.BuildWalls(0, metaData.ScreenWidth, metaData.ScreenHeight)

	// enemies
	enemies := make([]*entity.Enemy, 0)
	for range metaData.TotalEnemies {
		enemy := &entity.Enemy{}
		enemy.Reset(g.Environment, &g.MetaData)
		enemies = append(enemies, enemy)
	}
	g.Enemies = enemies
	return nil
}

func main() {
	fmt.Println("Game Starting...")
	// sw, sh := ebiten.Monitor().Size()
	g := common.GameMetaData{}

	// screen
	g.ScreenHeight = 500
	g.ScreenWidth = 500
	g.TotalEnemies = 1
	g.CurrentEnemyCount = g.TotalEnemies
	g.BoundryEdgeBuffer = 15

	game := &Game{
		BlockSize:           10,
		EnableBlockCollider: true,
	}
	game.NewLevel(g)

	// player
	game.Player = &entity.Player{}
	err := game.Player.Reset()
	if err != nil {
		panic(err)
	}

	// window setting
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("TTG Arena")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	f, err := os.Open("resources\\images\\player.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowIcon([]image.Image{i})

	// start game
	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
