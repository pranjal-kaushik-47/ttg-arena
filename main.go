package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
	"tag-game-v2/common"
	"tag-game-v2/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Bugs:
// wall collider getting disabled when moving toward and back

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
	Player      *entity.Player
	Enemies     []*entity.Enemy
	Environment *entity.Environment
	ScreenText  *entity.ScreenText
	MetaData    common.GameMetaData

	// Temp Var for level editing
	BlockSize int

	// Game Engine Controls
	DrawSprite          bool
	DrawCollider        bool
	ChangeCursor        bool
	EnableBlockCollider bool
}

func (g *Game) EditLevel() error {
	// Edit Level

	if g.ChangeCursor {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
	x, y := ebiten.CursorPosition()

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.ChangeCursor = !g.ChangeCursor
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.DrawSprite = !g.DrawSprite
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyB) {
			g.EnableBlockCollider = !g.EnableBlockCollider
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.DrawCollider = !g.DrawCollider
			entity.RenderSpriteBoundingBox = !entity.RenderSpriteBoundingBox
		}
	}

	_, dy := ebiten.Wheel()
	g.BlockSize = g.BlockSize + int(dy)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.DrawSprite {
			g.Environment.BuildSquareWall(x, y, g.BlockSize, g.BlockSize, g.EnableBlockCollider)
		}
		if g.DrawCollider {
			g.Environment.DrawCollider(x, y)
		}
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

func (g *Game) Update() error {
	g.EditLevel()
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
	return nil
}

func (g *Game) UpdateDisplayText() {

	// msgs := []string{fmt.Sprintf("Level: %d", g.MetaData.Level), fmt.Sprintf("Sprite(D): %v", g.DrawSprite), fmt.Sprintf("Collider(Q): %v", g.DrawCollider), fmt.Sprintf("Cursor(C): %v", g.ChangeCursor), fmt.Sprintf("BlockCollider(B): %v", g.EnableBlockCollider)}
	displayLevel := entity.TextMap["level"]
	displayLevel.Variables = []any{g.MetaData.Level}

	// displayEnemy := entity.TextMap["enemy"]
	// displayEnemy.Variables = []any{len(g.Enemies)}
}

func (g *Game) Draw(screen *ebiten.Image) {
	clr := color.RGBA{40, 44, 52, 255}
	screen.Fill(clr)
	g.UpdateDisplayText()
	if g.ChangeCursor {
		x, y := ebiten.CursorPosition()
		rect := ebiten.NewImage(g.BlockSize, g.BlockSize)
		rect.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle

		// Draw the rectangle at the mouse cursor position
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x)-(float64(g.BlockSize)/2), float64(y)-(float64(g.BlockSize)/2))

		screen.DrawImage(rect, op)
	}
	for _, enemy := range g.Enemies {
		if enemy.Sprite.IsActive {
			enemy.Draw(screen)
		}
	}
	g.Environment.Draw(screen, g.MetaData, entity.Point{X: g.Player.Sprite.PosX, Y: g.Player.Sprite.PosY})
	g.Player.Draw(screen)
	entity.DrawAllText(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.MetaData.ScreenWidth, g.MetaData.ScreenHeight
}

func (g *Game) NewLevel(metaData common.GameMetaData) error {
	g.MetaData = metaData

	// environment
	g.Environment = &entity.Environment{}

	// load saved level by uncommenting the line below
	g.Environment.BuildWalls(0, metaData.ScreenWidth, metaData.ScreenHeight)

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
	g := common.GameMetaData{}

	// screen
	g.ScreenHeight = 500
	g.ScreenWidth = 500
	g.TotalEnemies = 2
	g.CurrentEnemyCount = g.TotalEnemies
	g.BoundryEdgeBuffer = 15

	game := &Game{
		BlockSize:  10,
		ScreenText: &entity.ScreenText{},
	}
	entity.TextMap["level"] = &entity.TextMessage{
		Message:   "Level %v",
		Variables: []any{1},
		Position:  &entity.Point{X: 20, Y: 10},
		Size:      10,
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
