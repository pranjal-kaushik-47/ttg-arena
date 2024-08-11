package ecs

import (
	"encoding/json"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Position
	Sprite      Sprite
	Velocity    float64
	BoundingBox BoundingBox

	KillRadius float64
}

func (p *Player) Init() error {
	var playerConfig map[string]interface{}
	jsonData, err := os.ReadFile("resources/gamedata/player/player.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &playerConfig)
	if err != nil {
		return err
	}
	p.Position = Position{
		X: playerConfig["starting_x"].(float64),
		Y: playerConfig["starting_y"].(float64),
	}
	p.Velocity = playerConfig["velocity"].(float64)
	p.KillRadius = playerConfig["kill_radius"].(float64)
	p.Sprite = Sprite{
		ImageSource: playerConfig["sprite"].(string),
		Height:      playerConfig["height"].(float64),
		Width:       playerConfig["width"].(float64),
		IsActive:    true,
	}
	err = p.Sprite.Init()
	if err != nil {
		return err
	}
	p.BoundingBox.Init(p.Position.X, p.Position.Y, p.Sprite.Height, p.Sprite.Width)

	return nil
}

func (p *Player) Update(env Environment, screenHeight, screenWidth int) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.MoveTo(p.Position.X, p.Position.Y-p.Velocity, &p.BoundingBox, env, screenHeight, screenWidth)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.MoveTo(p.Position.X, p.Position.Y+p.Velocity, &p.BoundingBox, env, screenHeight, screenWidth)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.MoveTo(p.Position.X-p.Velocity, p.Position.Y, &p.BoundingBox, env, screenHeight, screenWidth)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.MoveTo(p.Position.X+p.Velocity, p.Position.Y, &p.BoundingBox, env, screenHeight, screenWidth)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Sprite.Draw(screen, p.Position.X, p.Position.Y)
}
