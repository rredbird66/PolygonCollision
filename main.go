package main

////////////////////////////////////////////////////////////////////////////////
import (
	"encoding/json"
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"sort"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
//	Structures
////////////////////////////////////////////////////////////////////////////////
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
	Id     int      `json:"id"`
	Layer  int      `json:"layer"`
	Points []Point2 `json:"points"`
}

type Edge struct {
	Begin Point2  `json:"begin"`
	End   Point2  `json:"end"`
	Angle float64 `json:"angle"`
	Layer int     `json:"layer"`
	Open  bool    `json:"open"`
}

type Point2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

////////////////////////////////////////////////////////////////////////////////
// Global variables
////////////////////////////////////////////////////////////////////////////////
var g_settings GlobalSettings
var g_drawParams DrawParams

var g_polygonArray []Polygon
var g_edgesArray []Edge
var g_pointsArray []Point2
var g_processedEdgesArray []Edge

var g_tempPolygon Polygon

////////////////////////////////////////////////////////////////////////////////
// Functions
////////////////////////////////////////////////////////////////////////////////
func printEdges() {
	for idx, edge := range g_edgesArray {
		fmt.Printf("LOG:: Edge #%d | BeginX: %f BeginY: %f\n", idx, edge.Begin.X, edge.Begin.Y)
		fmt.Printf("LOG:: Edge #%d | EndX: %f EndY: %f\n", idx, edge.End.X, edge.End.Y)
	}
}

////////////////////////////////////////////////////////////////////////////////
// Converter functions
////////////////////////////////////////////////////////////////////////////////
func p2rlP(point Point2) rl.Vector2 {
	return rl.Vector2{point.X, point.Y}
}

func rlP2p(point rl.Vector2) Point2 {
	return Point2{point.X, point.Y}
}

////////////////////////////////////////////////////////////////////////////////
// Grafic related functions
////////////////////////////////////////////////////////////////////////////////
func drawCanvas() {
	counter := int32(0)
	for counter < g_settings.screenWidth {
		rl.DrawLine(counter, 0, counter, g_settings.screenHeight, rl.LightGray)
		rl.DrawLine(0, counter, g_settings.screenWidth, counter, rl.LightGray)
		counter += g_settings.spacer
	}
}

func drawFigures() {
	for _, point := range g_tempPolygon.Points {
		rl.DrawCircle(int32(point.X), int32(point.Y), 5, rl.Red)
	}
	var tempColor rl.Color
	for _, figure := range g_polygonArray {
		if figure.Layer == 0 {
			tempColor = rl.NewColor(0, 0, 255, 100)
		} else {
			tempColor = rl.NewColor(255, 255, 0, 100)
		}
		if len(figure.Points) == 2 {
			tempColor.A = 255
			rl.DrawLineEx(p2rlP(figure.Points[0]), p2rlP(figure.Points[1]), 4, tempColor)
		} else {
			for i := 0; (i + 1) < len(figure.Points); i++ {
				rl.DrawTriangle(p2rlP(figure.Points[i+1]), p2rlP(figure.Points[i]), p2rlP(figure.Points[0]), tempColor)
			}
		}
	}
}

func drawPoints() {
	for _, figure := range g_polygonArray {
		for _, point := range figure.Points {
			rl.DrawCircle(int32(point.X), int32(point.Y), 5, rl.DarkGreen)
		}
	}
	for _, point := range g_pointsArray {
		rl.DrawCircle(int32(point.X), int32(point.Y), 5, rl.Green)
	}
}

func drawEdges() {
	for idx, edge := range g_edgesArray {
		rl.DrawLineEx(p2rlP(edge.Begin), p2rlP(edge.End), 1, rl.DarkBlue)
		rl.DrawText("edge #"+strconv.Itoa(idx)+" start", int32(edge.Begin.X), int32(edge.Begin.Y)+10, 10, rl.Red)
		rl.DrawText("edge #"+strconv.Itoa(idx)+" end", int32(edge.End.X), int32(edge.End.Y)-10, 10, rl.Red)
		//rl.DrawCircleV(edge.Begin, 10, rl.White)
		//rl.DrawCircleV(edge.End, 10, rl.Black)
	}
}

func drawStats() {

	var heightCounter int32 = 20

	counterF := func(x *int32) int32 {
		*x = *x + 20
		return *x
	}

	rl.DrawText("FPS:     "+strconv.Itoa(int(rl.GetFPS())), 20, counterF(&heightCounter), 20, rl.DarkGray)
	rl.DrawText("Figures: "+strconv.Itoa(len(g_polygonArray)), 20, counterF(&heightCounter), 20, rl.DarkGray)

	rl.DrawText("LMC - add vertex to polygon", 20, counterF(&heightCounter), 20, rl.DarkGray)
	rl.DrawText("RMC - create polygon from added vertexes", 20, counterF(&heightCounter), 20, rl.DarkGray)
	rl.DrawText("H - switch layer", 20, counterF(&heightCounter), 20, rl.DarkGray)
	rl.DrawText("S - output JSON-formatted figures in terminal", 20, counterF(&heightCounter), 20, rl.DarkGray)
	rl.DrawText("Z - exit", 20, counterF(&heightCounter), 20, rl.DarkGray)
}

////////////////////////////////////////////////////////////////////////////////
//  Input related functions
////////////////////////////////////////////////////////////////////////////////
func processKeys() {
	if rl.IsKeyPressed(rl.KeyH) {
		g_settings.layerSwitch = !g_settings.layerSwitch
	}
	if rl.IsKeyPressed(rl.KeyE) {
		g_edgesArray = nil
		getEdges()
		transformEdges()
	}
	if rl.IsKeyPressed(rl.KeyL) {
		sweepLine()
	}
	if rl.IsKeyPressed(rl.KeyP) {
		data, err := json.MarshalIndent(g_edgesArray, "", " ")
		if err != nil {
			log.Fatalf("ERROR JSON: %s", err)
		}
		fmt.Printf("%s\n", data)
	}
	if rl.IsKeyPressed(rl.KeyS) {
		data, err := json.MarshalIndent(g_polygonArray, "", " ")
		if err != nil {
			log.Fatalf("ERROR JSON: %s", err)
		}
		fmt.Printf("%s\n", data)
	}
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g_tempPolygon.Points = append(g_tempPolygon.Points, nearestPoint(rl.GetMousePosition()))
	}
	if rl.IsMouseButtonPressed(rl.MouseRightButton) && g_tempPolygon.Points != nil {
		if g_settings.layerSwitch {
			g_tempPolygon.Layer = 1
		} else {
			g_tempPolygon.Layer = 0
		}
		g_tempPolygon.Id = len(g_polygonArray)
		g_polygonArray = append(g_polygonArray, g_tempPolygon)
		g_tempPolygon.Points = nil
	}
}

////////////////////////////////////////////////////////////////////////////////
//  Logic functions
////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
//	Point functions
////////////////////////////////////////////////////////////////////////////////
func nearestPoint(point rl.Vector2) Point2 {
	point.X = float32(math.Round(float64(point.X)/40.0)) * 40
	point.Y = float32(math.Round(float64(point.Y)/40.0)) * 40
	return rlP2p(point)
}

func pointArrayContains(pointsArray []rl.Vector2, other rl.Vector2) bool {
	for _, point := range pointsArray {
		if point.Y == other.Y {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// Edge functions
////////////////////////////////////////////////////////////////////////////////
func getAngle(p1 Point2, p2 Point2) float64 {
	return math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X))
}

// Closure
func getCollision(e1 Edge) func(e2 Edge) (bool, Point2) {
	return func(e2 Edge) (bool, Point2) {
		var crossingPoint rl.Vector2
		res := rl.CheckCollisionLines(p2rlP(e1.Begin), p2rlP(e1.End), p2rlP(e2.Begin), p2rlP(e2.End), &crossingPoint)
		return res, rlP2p(crossingPoint)
	}
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
	for _, figure := range g_polygonArray {
		var tempEdge Edge
		for i := 0; (i + 1) < len(figure.Points); i++ {
			if figure.Points[i].Y != figure.Points[i+1].Y {
				tempEdge = Edge{figure.Points[i], figure.Points[i+1], getAngle(figure.Points[i], figure.Points[i+1]), figure.Layer, true}
				fmt.Printf("Angle: %f\n", tempEdge.Angle*360/math.Pi)
				if !edgesArrayContains(g_edgesArray, tempEdge) {
					g_edgesArray = append(g_edgesArray, tempEdge)
				}
			}
		}
		tempEdge = Edge{figure.Points[len(figure.Points)-1], figure.Points[0], getAngle(figure.Points[len(figure.Points)-1], figure.Points[0]), figure.Layer, true}
		if !edgesArrayContains(g_edgesArray, tempEdge) {
			g_edgesArray = append(g_edgesArray, tempEdge)
		}
	}
}

func transformEdges() {
	for idx, edge := range g_edgesArray {
		if edge.Begin.Y > edge.End.Y {
			g_edgesArray[idx].Begin, g_edgesArray[idx].End = edge.End, edge.Begin
		}
	}
	sort.SliceStable(g_edgesArray, func(i, j int) bool {
		return g_edgesArray[i].Begin.Y < g_edgesArray[j].Begin.Y
	})
	//printEdges()
}

////////////////////////////////////////////////////////////////////////////////
// Sweep line algorithm
////////////////////////////////////////////////////////////////////////////////
func sweepLine() {

	g_pointsArray = nil
	for idx1, edge1 := range g_edgesArray {
		collider := getCollision(edge1)
		for idx2 := idx1 + 1; idx2 < len(g_edgesArray); idx2++ {
			res, point := collider(g_edgesArray[idx2])
			if res {
				g_pointsArray = append(g_pointsArray, point)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Main function
////////////////////////////////////////////////////////////////////////////////
func main() {

	g_settings.layerSwitch = false
	g_settings.screenWidth = int32(1000)
	g_settings.screenHeight = int32(1000)
	g_settings.spacer = g_settings.screenWidth / 25

	rl.InitWindow(g_settings.screenWidth, g_settings.screenHeight, "PolygonCollision")
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyZ)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		processKeys()
		drawCanvas()
		drawFigures()
		drawEdges()
		drawPoints()

		drawStats()
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

////////////////////////////////////////////////////////////////////////////////
