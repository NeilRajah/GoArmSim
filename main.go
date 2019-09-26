//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	// "fmt"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	// "math"
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
		colors := robotArm.getColor()
		ctx.SetRGB255(colors[0], colors[1], colors[2])

		drawRobot(ctx) //draw the robot to the screen

		//display the data to the screen
		ctx.SetColor(colornames.White)
		displayData(ctx)

		//restore canvas state
		ctx.Pop()
	})

} //end main

//create the arm
func createArm() {
	//set the values for the arm
	startPt := point{float64(WIDTH / 2), 0}
	length := 24.0
	angle := ToRadians(0)

	//create the arm
	robotArm = Arm{
		start:    startPt,
		length:   length,
		angle:    angle,
		topSpeed: 60} //60 degrees per second, 10RPM

	robotArm.setOutput(0.5) //50%
} //end createArm

//Update the arm for drawing purposes
func updateModel() {
	robotArm.update()
	if robotArm.getAngleDeg() > 45 {
		robotArm.stop()
	}
} //end updateModel

//SIMULATOR

//draw the robot to the display
//ctx *canvas.Context - responsible for drawing
func drawRobot(ctx *canvas.Context) {
	//draw the robot arm as lines between the joint points
	ctx.DrawLine(robotArm.start.x, robotArm.start.y,
		robotArm.getEndPtPxl().x, robotArm.getEndPtPxl().y)

	ctx.Stroke()
} //end drawRobot

//display the parameters of the robot onto the screen
//ctx *canvas.Context - responsible for drawing
func displayData(ctx *canvas.Context) {
	//display the start and end coords
	displayPointCoords(ctx, robotArm.getStartPtIn(), 1400, FONT_SIZE)
	displayPointCoords(ctx, robotArm.getEndPtIn(), 1400, 70+FONT_SIZE)

	//display the angle of the arm (combine with point and make into helper function)
	drawFloat(ctx, robotArm.angle, 1400, 140+FONT_SIZE)
	drawFloat(ctx, robotArm.getAngleDeg(), 1600, 140+FONT_SIZE)
} //end displayData
