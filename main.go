package main

import (
	"tag-game-v2/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	levelData    *ecs.Level
	ScreenWidth  int
	ScreenHeight int
}

func (g *Game) Update() error {
	for _, enemy := range g.levelData.Enemies {
		enemy.Update(g.levelData.Player, &g.levelData.Environment, g.ScreenHeight, g.ScreenWidth)
	}
	g.levelData.Player.Update(g.levelData.Environment, g.ScreenHeight, g.ScreenWidth)
	if g.levelData.Environment.AliveEnemyCount == 0 {
		g.levelData.NextLevel()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, enemy := range g.levelData.Enemies {
		enemy.Draw(screen)
		enemy.BoundingBox.Polygon.Draw(screen)
	}
	g.levelData.Player.Draw(screen)
	g.levelData.Player.BoundingBox.Polygon.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth + 30, g.ScreenHeight + 30
}

func main() {
	var err error
	game := &Game{
		ScreenWidth:  1000,
		ScreenHeight: 500,
		levelData: &ecs.Level{
			Id:     1,
			Player: ecs.Player{},
		},
	}
	err = game.levelData.Init()
	if err != nil {
		panic(err)
	}
	err = game.levelData.Player.Init()
	if err != nil {
		panic(err)
	}
	err = ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
