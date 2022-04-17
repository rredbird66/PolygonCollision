package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"fmt"
	"log"
	"strconv"
	"encoding/json"
)

type GlobalSettings struct {
	layerSwitch bool
	screenWidth int32
	screenHeight int32
	spacer int32
}

type DrawParams struct {
	triColor rl.Color
}

type Triangle struct {
	vertexA rl.Vector2
	vertexB rl.Vector2
	vertexC rl.Vector2
	color rl.Color 
}

type PolygonJSON struct {
	Id int `json:"id"`
	Layer int `json:"layer"`
	Points []PointJSON
}

type PointJSON struct {
	PointX float32 `json:"x"`
	PointY float32 `json:"y"`
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
		for _, figure := range triangleArray {
			PA := PointJSON{figure.vertexA.X,figure.vertexA.Y}
			PB := PointJSON{figure.vertexB.X,figure.vertexB.Y}
			PC := PointJSON{figure.vertexC.X,figure.vertexC.Y}
			if figure.color == rl.NewColor(0, 82, 172, 100) {
				layer = 1
			} else if figure.color == rl.NewColor(230, 41, 55, 100) {
				layer = 0
			}
			temTriJS := PolygonJSON{0, layer, []PointJSON{PA,PB,PC}}
			data, err := json.MarshalIndent(temTriJS, "", " ")
			if err != nil {
				log.Fatalf("ERROR JSON: %s", err)
			}
			fmt.Printf("%s\n", data)
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		switch vertexCounter {
		case 0:
			tempTriangle.vertexC = rl.GetMousePosition()
		case 1:
			tempTriangle.vertexB = rl.GetMousePosition()
		case 2:
			tempTriangle.vertexA = rl.GetMousePosition()
		}
		vertexCounter++
	}
	if vertexCounter == 3 {
		if settings.layerSwitch {
			DarkBlueTrans := rl.NewColor(0, 82, 172, 100)
			tempTriangle.color = DarkBlueTrans
		} else {
			RedTrans := rl.NewColor(230, 41, 55, 100)
			tempTriangle.color = RedTrans
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
		rl.DrawTriangle(figure.vertexA, figure.vertexB, figure.vertexC, figure.color)
	}
}

func drawStats() {
	rl.DrawText("Vertex  #" + strconv.Itoa(vertexCounter), 20, 20, 20, rl.DarkGray)
	rl.DrawText("FPS:     " + strconv.Itoa(int(rl.GetFPS())), 20, 40, 20, rl.DarkGray)
	rl.DrawText("Figures: " + strconv.Itoa(len(triangleArray)), 20, 60, 20, rl.DarkGray)
}

func main() {

	settings.layerSwitch = false
	settings.screenWidth = int32(1200)
	settings.screenHeight = int32(1200)
	settings.spacer = settings.screenWidth / 40

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
