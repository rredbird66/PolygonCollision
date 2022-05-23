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
	//angle float32
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
var edgesArray []Edge

var tempPolygon Polygon

func processKeys() {
	if rl.IsKeyPressed(rl.KeyH) {
		settings.layerSwitch = !settings.layerSwitch
	}
	if rl.IsKeyPressed(rl.KeyE) {
		edgesArray = nil
		getEdges()
		transformEdges()
	}
	if rl.IsKeyPressed(rl.KeyP) {
		printEdges()
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
	if rl.IsMouseButtonPressed(rl.MouseRightButton) && tempPolygon.points != nil {
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
			tempColor := figure.color
			tempColor.A = 255
			rl.DrawLineEx(figure.points[0], figure.points[1], 4, tempColor)
		} else {
			for i := 0; (i + 1) < len(figure.points); i++ {
				rl.DrawTriangle(figure.points[i+1], figure.points[i], figure.points[0], figure.color)
			}
		}
	}
}

func drawPoints() {
	for _, figure := range polygonArray {
		for _, point := range figure.points {
			rl.DrawCircle(int32(point.X), int32(point.Y), 5, rl.DarkGreen)
		}
	}
}

func drawEdges() {
	for idx, edge := range edgesArray {
		rl.DrawLineEx(edge.Begin, edge.End, 1, rl.DarkBlue)
		rl.DrawText("edge #"+strconv.Itoa(idx)+" start", int32(edge.Begin.X), int32(edge.Begin.Y)+10, 10, rl.Red)
		rl.DrawText("edge #"+strconv.Itoa(idx)+" end", int32(edge.End.X), int32(edge.End.Y)-10, 10, rl.Red)
		//rl.DrawCircleV(edge.Begin, 10, rl.White)
		//rl.DrawCircleV(edge.End, 10, rl.Black)
	}
}

func pointArrayContains(pointsArray []rl.Vector2, other rl.Vector2) bool {
	for _, point := range pointsArray {
		if point.Y == other.Y {
			return true
		}
	}
	return false
}

func edgesArrayContains(arr []Edge, other Edge) bool {
	for _, edge := range arr {
		if (edge.Begin == other.Begin && edge.End == other.End) || (edge.Begin == other.End && edge.End == other.Begin) {
			return true
		}
	}
	return false
}

func getEdges() {
	for _, figure := range polygonArray {
		var tempEdge Edge
		for i := 0; (i + 1) < len(figure.points); i++ {
			tempEdge = Edge{figure.points[i], figure.points[i+1]}
			if !edgesArrayContains(edgesArray, tempEdge) {
				edgesArray = append(edgesArray, tempEdge)
			}
		}
		tempEdge = Edge{figure.points[len(figure.points)-1], figure.points[0]}
		if !edgesArrayContains(edgesArray, tempEdge) {
			edgesArray = append(edgesArray, tempEdge)
		}
	}
}

func transformEdges() {
	for idx, edge := range edgesArray {
		if edge.Begin.Y > edge.End.Y {
			edgesArray[idx].Begin, edgesArray[idx].End = edge.End, edge.Begin
		}
	}
	sort.SliceStable(edgesArray, func(i, j int) bool {
		return edgesArray[i].Begin.Y < edgesArray[j].Begin.Y
	})
	//printEdges()
}

func printEdges() {
	for idx, edge := range edgesArray {
		fmt.Printf("LOG:: Edge #%d | BeginX: %f BeginY: %f\n", idx, edge.Begin.X, edge.Begin.Y)
		fmt.Printf("LOG:: Edge #%d | EndX: %f EndY: %f\n", idx, edge.End.X, edge.End.Y)
	}
}

func sweepLine() {
	/*
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
	*/
}

func drawStats() {
	rl.DrawText("FPS:     "+strconv.Itoa(int(rl.GetFPS())), 20, 40, 20, rl.DarkGray)
	rl.DrawText("Figures: "+strconv.Itoa(len(polygonArray)), 20, 60, 20, rl.DarkGray)

	rl.DrawText("LMC - add vertex to polygon", 20, 80, 20, rl.DarkGray)
	rl.DrawText("RMC - create polygon from added vertexes", 20, 100, 20, rl.DarkGray)
	rl.DrawText("H - switch layer", 20, 120, 20, rl.DarkGray)
	rl.DrawText("S - output JSON-formatted figures in terminal", 20, 140, 20, rl.DarkGray)
	rl.DrawText("Z - exit", 20, 160, 20, rl.DarkGray)
}

func main() {

	settings.layerSwitch = false
	settings.screenWidth = int32(1200)
	settings.screenHeight = int32(1200)
	settings.spacer = settings.screenWidth / 30

	rl.InitWindow(settings.screenWidth, settings.screenHeight, "PolygonCollision")
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyZ)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		processKeys()
		drawCanvas()
		drawFigures()
		sweepLine()
		drawEdges()
		drawPoints()

		drawStats()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
