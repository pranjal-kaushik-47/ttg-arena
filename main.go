package main

import (
	"fmt"
	"image"
	"os"
	"tag-game-v2/common"
	"tag-game-v2/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player      *entity.Player
	Enemies     []*entity.Enemy
	Environment *entity.Environment
	MetaData    common.GameMetaData
}

func (g *Game) Update() error {
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

func (g *Game) Draw(screen *ebiten.Image) {
	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}
	g.Environment.Draw(screen, g.MetaData)
	g.Player.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.MetaData.ScreenWidth, g.MetaData.ScreenHeight
}

func (g *Game) NewLevel(metaData common.GameMetaData) error {
	g.MetaData = metaData

	// environment
	g.Environment = &entity.Environment{}
	g.Environment.ScreenHeight = metaData.ScreenHeight
	g.Environment.ScreenWidth = metaData.ScreenWidth
	g.Environment.BuildWalls(0, metaData.ScreenWidth, metaData.ScreenHeight)

	// enemies
	enemies := make([]*entity.Enemy, 0)
	for range metaData.TotalEnemies {
		enemy := &entity.Enemy{}
		enemy.Reset(g.Environment)
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

	game := &Game{}
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
