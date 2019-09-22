//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	// "github.com/faiface/pixel/pixelgl"
	// "fmt"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"strconv"
	"time"
)

//constants
const WIDTH int = 1920  //width of the window
const HEIGHT int = 1080 //height of the window

//variables
var c canvas.Canvas //canvas instance
var robotArm arm    //arm struct

//create the arm struct to be used and run the graphics
func main() {
	//create a new canvas instance
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     WIDTH,
		Height:    HEIGHT,
		FrameRate: 60, //same as regular monitor
		Title:     "Arm Simulator",
	})

	//set up the canvas
	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.Lightgray) //set the bg color
		ctx.Clear()                        //empty the canvas
		ctx.SetColor(colornames.Black)     //set the drawing color
		ctx.SetLineWidth(20)               //set the line width
		ctx.LoadFontFace("../HelveticaNeue.ttf", 40)
	})

	startPt := point{float64(WIDTH / 2), 0}
	length := 500.0
	angle := (45 * (2 * 3.1415)) / 360.0

	robotArm := arm{
		start:  startPt,
		length: length,
		angle:  angle}

	c.Draw(func(ctx *canvas.Context) {
		go func() {
			robotArm.angle += 0.02
		}()

		ctx.Push()
		ctx.DrawLine(robotArm.start.x, robotArm.start.y,
			robotArm.getEndPt().x, robotArm.getEndPt().y)
		displayPointCoords(ctx, robotArm.getEndPt(), 1400, 200)

		ctx.Push()
		ctx.InvertY()
		ctx.DrawString(strconv.FormatFloat(angle, 'f', -1, 64), 1400, 300)
		ctx.Pop()
		ctx.Pop()

		ctx.Stroke()
		time.Sleep(time.Millisecond * 200)

	})

} //end main

//display the end-point coordinate of the arm
//*canvas.Context ctx - coordinates of the endpoint
//point p - endpoint in cartesian space
//int x - x coordinate for the string
//int y - y coordinate for the string
func displayPointCoords(ctx *canvas.Context, p point, x, y float64) {
	coordinates := "(" + strconv.FormatFloat(p.x-float64(WIDTH/2), 'f', 2, 64) + " , " +
		strconv.FormatFloat(p.y, 'f', 2, 64) + ")"

	ctx.Push()

	ctx.InvertY()
	ctx.DrawString(coordinates, x, y)

	ctx.Pop()
} //end displayPointCoords
