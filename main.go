package main

import (
	"encoding/json"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"strconv"
)

type GlobalSettings struct {
	layerSwitch  bool
	screenWidth  int32
	screenHeight int32
	spacer       int32
}

type DrawParams struct {
	triColor rl.Color
}

type Triangle struct {
	points [3]rl.Vector2
	color  rl.Color
}

type PolygonJSON struct {
	Id     int `json:"id"`
	Layer  int `json:"layer"`
	Points []PointJSON
}

type PointJSON struct {
	PointX float32 `json:"x"`
	PointY float32 `json:"y"`
}

func nearestPoint(point rl.Vector2) rl.Vector2 {
	fmt.Printf("x:%f y:%f\n", point.X, point.Y)

	point.X = float32(math.Round(float64(point.X)/40.0)) * 40
	point.Y = float32(math.Round(float64(point.Y)/40.0)) * 40

	fmt.Printf("x:%f y:%f\n", point.X, point.Y)
	return point
}

func collisions(a, b Triangle) []rl.Vector2 {
	var collisionPoints []rl.Vector2
	var tempPoint rl.Vector2

	for i := 0; i < len(a.points); i++ {
		for j := 0; j < len(a.points); j++ {
			for k := 0; k < len(b.points); k++ {
				for l := 0; l < len(b.points); l++ {
					if rl.CheckCollisionLines(a.points[i], a.points[j], b.points[k], b.points[l], &tempPoint) {
						collisionPoints = append(collisionPoints, tempPoint)
					}
				}
			}
		}
	}

	keys := make(map[rl.Vector2]bool)
	uniqueCollisionPoints := []rl.Vector2{}
	for _, entry := range collisionPoints {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueCollisionPoints = append(uniqueCollisionPoints, entry)
		}
	}
	return uniqueCollisionPoints
}

var settings GlobalSettings
var drawParams DrawParams

var triangleArray []Triangle

var tempTriangle Triangle
var vertexCounter int

func processKeys() {
	if rl.IsKeyPressed(rl.KeyH) {
		settings.layerSwitch = !settings.layerSwitch
	}
	if rl.IsKeyPressed(rl.KeyS) {
		var layer int
		for id, figure := range triangleArray {
			PA := PointJSON{figure.points[0].X, figure.points[0].Y}
			PB := PointJSON{figure.points[1].X, figure.points[1].Y}
			PC := PointJSON{figure.points[2].X, figure.points[2].Y}
			if figure.color == rl.NewColor(255, 255, 0, 100) {
				layer = 1
			} else if figure.color == rl.NewColor(0, 0, 255, 100) {
				layer = 0
			}
			temTriJS := PolygonJSON{id, layer, []PointJSON{PA, PB, PC}}
			data, err := json.MarshalIndent(temTriJS, "", " ")
			if err != nil {
				log.Fatalf("ERROR JSON: %s", err)
			}
			fmt.Printf("%s\n", data)
		}
		fmt.Printf("%d\n", len(collisions(triangleArray[0], triangleArray[1])))
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		switch vertexCounter {
		case 0:
			tempTriangle.points[2] = nearestPoint(rl.GetMousePosition())
		case 1:
			tempTriangle.points[1] = nearestPoint(rl.GetMousePosition())
		case 2:
			tempTriangle.points[0] = nearestPoint(rl.GetMousePosition())
		}
		vertexCounter++
	}
	if vertexCounter == 3 {
		if settings.layerSwitch {
			tempTriangle.color = rl.NewColor(255, 255, 0, 100)
		} else {
			tempTriangle.color = rl.NewColor(0, 0, 255, 100)
		}
		triangleArray = append(triangleArray, tempTriangle)
		vertexCounter = 0
	}
}

func drawCanvas() {
	counter := int32(0)
	for counter < settings.screenWidth {
		rl.DrawLine(counter, 0, counter, settings.screenHeight, rl.LightGray)
		rl.DrawLine(0, counter, settings.screenWidth, counter, rl.LightGray)
		counter += settings.spacer
	}
}

func drawFigures() {
	for _, figure := range triangleArray {
		rl.DrawTriangle(figure.points[0], figure.points[1], figure.points[2], figure.color)
	}

	if len(triangleArray) > 1 {
		for i := 0; i < len(triangleArray); i++ {
			for j := i + 1; j < len(triangleArray); j++ {
				for _, point := range collisions(triangleArray[i], triangleArray[j]) {
					rl.DrawCircleV(point, float32(4), rl.Red)
				}
			}
		}
	}

}

func drawStats() {
	rl.DrawText("Vertex  #"+strconv.Itoa(vertexCounter), 20, 20, 20, rl.DarkGray)
	rl.DrawText("FPS:     "+strconv.Itoa(int(rl.GetFPS())), 20, 40, 20, rl.DarkGray)
	rl.DrawText("Figures: "+strconv.Itoa(len(triangleArray)), 20, 60, 20, rl.DarkGray)
}

func main() {

	settings.layerSwitch = false
	settings.screenWidth = int32(1200)
	settings.screenHeight = int32(1200)
	settings.spacer = settings.screenWidth / 30

	vertexCounter = 0

	rl.InitWindow(settings.screenWidth, settings.screenHeight, "Some random app")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		processKeys()
		drawCanvas()
		drawFigures()

		drawStats()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
