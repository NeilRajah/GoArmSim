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
	// "math"
	"strconv"
	// "time"
)

//constants
const WIDTH int = 1920       //width of the window
const HEIGHT int = 1080      //height of the window
const FPS int = 60           //frame rate of the animation
const FONT_SIZE float64 = 60 //font size

//variables
var c canvas.Canvas //canvas instance
var robotArm Arm    //arm struct

//create the arm struct to be used and run the graphics
func main() {
	//create a new canvas instance
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     WIDTH,
		Height:    HEIGHT,
		FrameRate: FPS, //same as regular monitor
		Title:     "Arm Simulator",
	})

	//set up the canvas
	c.Setup(func(ctx *canvas.Context) {
		setUpCanvas(ctx)
	})

	//create the arm
	createArm()

	//create the arm
	c.Draw(func(ctx *canvas.Context) {
		//update the arm
		updateModel()

		//clear the canvas
		ctx.SetColor(colornames.Black) //set the bg color
		ctx.Clear()                    //empty the canvas

		//save canvas state
		ctx.Push()

		//switch to the arm color
		ctx.SetColor(colornames.White)

		//draw the robot arm (make helper func)
		ctx.DrawLine(robotArm.start.x, robotArm.start.y,
			robotArm.getEndPt().x, robotArm.getEndPt().y)
		// fmt.Println("x0:", robotArm.start.x, "y0:", robotArm.start.y, "x1:", robotArm.getEndPt().x, "y1:", robotArm.getEndPt().y)

		//display the start and aend coords
		displayPointCoords(ctx, robotArm.start, 1400, FONT_SIZE)
		displayPointCoords(ctx, robotArm.getEndPt(), 1400, 70+FONT_SIZE)

		//display the angle of the arm (combine with point and make into helper function)
		drawFloat(ctx, robotArm.angle, 1400, 140+FONT_SIZE)

		//draw
		ctx.Stroke()

		//restore canvas state
		ctx.Pop()
	})

} //end main

//create the arm
func createArm() {
	//set the values for the arm
	startPt := point{float64(WIDTH / 2), 0}
	length := 700.0
	angle := ToRadians(0)

	//create the arm
	robotArm = Arm{
		start:  startPt,
		length: length,
		angle:  angle,
		vel:    0}
} //end createArm

//Update the arm for drawing purposes
func updateModel() {
	robotArm.setVelDPS(10)
	robotArm.update()
} //end updateModel

//DRAWING

//set up the canvas
//*canvas.Context ctx - responsible for drawing
func setUpCanvas(ctx *canvas.Context) {
	ctx.SetColor(colornames.Lightgray)                  //set the bg color
	ctx.Clear()                                         //empty the canvas
	ctx.SetColor(colornames.Black)                      //set the drawing color
	ctx.SetLineWidth(30)                                //set the line width
	ctx.LoadFontFace("../HelveticaNeue.ttf", FONT_SIZE) //set the font
} //end setUpCanvas

//display the end-point coordinate of the arm
//*canvas.Context ctx - responsible for drawing
//point p - endpoint in cartesian space
//float64 x - x coordinate for the string
//float64 y - y coordinate for the string
func displayPointCoords(ctx *canvas.Context, p point, x, y float64) {
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
//float64 x - x coordinate for the string
//float64 y - y coordinate for the string
func drawFloat(ctx *canvas.Context, f, x, y float64) {
	ctx.Push()

	ctx.InvertY() //flip y value to draw string
	ctx.DrawString(strconv.FormatFloat(f, 'f', 2, 64), x, y)

	ctx.Pop()
} //end drawString
