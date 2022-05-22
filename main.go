package main

import (
	"encoding/json"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"sort"
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

type Polygon struct {
	Id     int
	Layer  int
	color  rl.Color
	points []rl.Vector2
}

type Edge struct {
	Begin rl.Vector2
	End   rl.Vector2
	angle float32
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
	point.X = float32(math.Round(float64(point.X)/40.0)) * 40
	point.Y = float32(math.Round(float64(point.Y)/40.0)) * 40
	return point
}

func collisions() int {
	return 0
}

var settings GlobalSettings
var drawParams DrawParams

var polygonArray []Polygon

var tempPolygon Polygon
var vertexCounter int

func processKeys() {
	if rl.IsKeyPressed(rl.KeyH) {
		settings.layerSwitch = !settings.layerSwitch
	}
	if rl.IsKeyPressed(rl.KeyS) {
		var layer int
		for id, figure := range polygonArray {
			var tempPointsArray []PointJSON
			for _, point := range figure.points {
				tempPointsArray = append(tempPointsArray, PointJSON{point.X, point.Y})
			}
			if figure.color == rl.NewColor(255, 255, 0, 100) {
				layer = 1
			} else if figure.color == rl.NewColor(0, 0, 255, 100) {
				layer = 0
			}
			tempPolyJS := PolygonJSON{id, layer, tempPointsArray}
			data, err := json.MarshalIndent(tempPolyJS, "", " ")
			if err != nil {
				log.Fatalf("ERROR JSON: %s", err)
			}
			fmt.Printf("%s\n", data)
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		tempPolygon.points = append(tempPolygon.points, nearestPoint(rl.GetMousePosition()))
	}
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if settings.layerSwitch {
			tempPolygon.color = rl.NewColor(255, 255, 0, 100)
		} else {
			tempPolygon.color = rl.NewColor(0, 0, 255, 100)
		}
		polygonArray = append(polygonArray, tempPolygon)
		tempPolygon.points = nil
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
	for _, point := range tempPolygon.points {
		rl.DrawCircle(int32(point.X), int32(point.Y), 5, rl.Red)
	}
	for _, figure := range polygonArray {
		if len(figure.points) == 2 {
			rl.DrawLineEx(figure.points[0], figure.points[1], 3, figure.color)
		} else {
			for i := 0; (i + 1) < len(figure.points); i++ {
				rl.DrawTriangle(figure.points[i+1], figure.points[i], figure.points[0], figure.color)
			}
		}
	}
}

func pointArrayContains(pointsArray []rl.Vector2, point rl.Vector2) bool {
	for _, a := range pointsArray {
		if a == point {
			return true
		}
	}
	return false
}

func sweepLine() {
	var tempPointsArray []rl.Vector2
	for _, figure := range polygonArray {
		for _, point := range figure.points {
			if !pointArrayContains(tempPointsArray, point) {
				tempPointsArray = append(tempPointsArray, point)
			}
		}
	}

	sort.SliceStable(tempPointsArray, func(i, j int) bool {
		return tempPointsArray[i].Y < tempPointsArray[j].Y
	})

	for idx, point := range tempPointsArray {
		rl.DrawText("line #"+strconv.Itoa(idx), 5, int32(point.Y), 20, rl.Red)
		rl.DrawLine(0, int32(point.Y), 1200, int32(point.Y), rl.Red)
	}
}

func drawStats() {
	rl.DrawText("Vertex  #"+strconv.Itoa(vertexCounter), 20, 20, 20, rl.DarkGray)
	rl.DrawText("FPS:     "+strconv.Itoa(int(rl.GetFPS())), 20, 40, 20, rl.DarkGray)
	rl.DrawText("Figures: "+strconv.Itoa(len(polygonArray)), 20, 60, 20, rl.DarkGray)
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
		sweepLine()

		drawStats()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
