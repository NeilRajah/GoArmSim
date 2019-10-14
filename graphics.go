//graphics
//Author: Neil Balaskandarajah
//Created on: 09/25/2019
//Graphics functions used for the simulator

package main

import (
	// "github.com/faiface/pixel"
	"strconv"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

const armWidth = 30.0 //armWidth is the thickness of the arm lines in pixels
const boxThick = 10.0 //boxThick is the thickness of the box lines in pixels

//set up the canvas
//*canvas.Context ctx - responsible for drawing
func setUpCanvas(ctx *canvas.Context) {
	ctx.SetColor(colornames.Lightgray)                 //set the bg color
	ctx.Clear()                                        //empty the canvas
	ctx.SetColor(colornames.Black)                     //set the drawing color
	ctx.SetLineWidth(30)                               //set the line width
	ctx.LoadFontFace("../HelveticaNeue.ttf", fontSize) //set the font
} //end setUpCanvas

//display the end-point coordinate of the arm
//*canvas.Context ctx - responsible for drawing
//point p - endpoint in cartesian space
//float64 x - x coordinate for the string
//float64 y - y coordinate for the string
func displayPointCoords(ctx *canvas.Context, p Point, x, y float64) {
	//show coordinates in (x,y) form
	coordinates := "(" + strconv.FormatFloat(p.x, 'f', 2, 64) + " , " +
		strconv.FormatFloat(p.y, 'f', 2, 64) + ")"

	ctx.Push()

	ctx.InvertY() //flip y value to draw the string
	ctx.DrawString(coordinates, x, y)

	ctx.Pop()
} //end displayPointCoords

//draw a float at a specified location
//*canvas.Context ctx - responsible for drawing
//float64 f - float to be drawn
//string s - message associated with float
//float64 x - x coordinate for the string
//float64 y - y coordinate for the string
func drawFloat(ctx *canvas.Context, f, x, y float64, s string) {
	ctx.Push()

	ctx.InvertY() //flip y value to draw string
	ctx.DrawString(s+": "+strconv.FormatFloat(f, 'f', 2, 64), x, y)

	ctx.Pop()
} //end drawString

//draw a point to the canvas
//*canvas.Context ctx - responsible for drawing
//Point p - point to draw
//float64 r - radius of circle at point
func drawPoint(ctx *canvas.Context, p Point, r float64) {
	ctx.Push()

	ctx.DrawCircle(p.x*float64(pixelToMeters)+float64(width)/2, p.y*float64(pixelToMeters), r)
	ctx.Stroke()

	ctx.Pop()
} //end drawPoint

//convert the mouse coordinates to cartesian coordinates
//Point mouse - point for the mouse
func mouseToMeter(mouse Point) Point {
	return Point{(mouse.x - float64(width)/2) / pixelToMeters, mouse.y / pixelToMeters}
} //end mouseToMeter
