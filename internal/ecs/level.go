package ecs

import (
	"encoding/json"
	"os"
	"strconv"
)

type Props struct {
	Sprite      Sprite
	Position    Position
	BoundingBox BoundingBox
}

type Environment struct {
	Objects         []Props
	AliveEnemyCount int
}

type Level struct {
	Id          int
	Environment Environment
	Player      Player
	Enemies     []*Enemy
}

type GameObjects interface {
	MoveTo()
	PolygonCollision()
}

func (p *Props) MoveTo(x, y float64) {}

func (p *Props) Update() {}

func (p *Props) Draw() {}

func (e *Environment) Colliding(Polygon Polygon) bool {

	c := false
	for _, obj := range e.Objects {
		if obj.BoundingBox.PolygonCollision(Polygon.Vertices) {
			c = true
		}
	}
	return c
}

func (l *Level) Init() error {
	enemyFile, err := os.ReadDir("resources/gamedata/enemies/level" + strconv.Itoa(l.Id) + "/")
	if err != nil {
		return err
	}

	enemies := make([]*Enemy, 0)

	for _, file := range enemyFile {
		var enemyConfig map[string]interface{}
		jsonData, err := os.ReadFile("resources/gamedata/enemies/level" + strconv.Itoa(l.Id) + "/" + file.Name())
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonData, &enemyConfig)
		if err != nil {
			return err
		}

		e := &Enemy{}

		e.Position = Position{
			X: enemyConfig["starting_x"].(float64),
			Y: enemyConfig["starting_y"].(float64),
		}
		e.Velocity = enemyConfig["velocity"].(float64)
		e.Sprite = Sprite{
			ImageSource: enemyConfig["sprite"].(string),
			Height:      enemyConfig["height"].(float64),
			Width:       enemyConfig["width"].(float64),
			IsActive:    true,
		}
		e.Eyesight = enemyConfig["eyesight"].(float64)
		err = e.Sprite.Init()
		if err != nil {
			return err
		}
		e.BoundingBox.Init(e.X, e.Y, e.Sprite.Height, e.Sprite.Width)

		enemies = append(enemies, e)

	}

	l.Enemies = enemies
	l.Environment.AliveEnemyCount = len(enemies)
	return nil
}

func (l *Level) NextLevel() {
	l.Id++
	l.Init()
}
