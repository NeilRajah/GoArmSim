//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

//Constants

const width int = 1920      //WIDTH is the width of the window
const height int = 1080     //HEIGHT is the height of the window
const fps int = 50          //FPS is the frame rate of the animation
const fontSize float64 = 60 //FONT_SIZE is the font size for the canvas

//variables
var c canvas.Canvas                         //canvas instance
var robotArm *Arm                           //arm struct
var bgColor color.RGBA = colornames.Black   //background color
var textColor color.RGBA = colornames.White //text color

//create the arm struct to be used and run the graphics
func main() {
	//create a new canvas instance
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     width,
		Height:    height,
		FrameRate: fps,
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
		ctx.SetColor(bgColor) //set the bg color
		ctx.Clear()           //empty the canvas

		//save canvas state
		ctx.Push()

		drawRobot(ctx) //draw the robot to the screen

		//display the data to the screen
		displayData(ctx)

		//restore canvas state
		ctx.Pop()
	})

} //end main

//create the arm
func createArm() {
	//set the values for the arm
	kP := 1.0
	kI := 0.0
	kD := 0.01
	robotArm = NewArm(1.0, 10.0, 159.3, 2, kP, kI, kD, "cim", math.Pi/4)
} //end createArm

//Update the arm for drawing purposes
func updateModel() {
	//move with PID control until the target is reached
	// robotArm.movePID(ToRadians(160), robotArm.angle, 1)
	// robotArm.movePIDFF(ToRadians(135), robotArm.angle, ToRadians(1))
	robotArm.voltage = calcFFArm(robotArm)
	robotArm.update()
} //end updateModel

//SIMULATOR

//draw the robot to the display
//ctx *canvas.Context - responsible for drawing
func drawRobot(ctx *canvas.Context) {
	//switch to the arm color
	colors := robotArm.getColor()
	ctx.SetRGB255(colors[0], colors[1], colors[2])

	//draw the robot arm as lines between the joint points
	ctx.DrawLine(robotArm.start.x, robotArm.start.y,
		robotArm.getEndPtPxl().x, robotArm.getEndPtPxl().y)

	ctx.Stroke()
} //end drawRobot

//display the parameters of the robot onto the screen
//ctx *canvas.Context - responsible for drawing
func displayData(ctx *canvas.Context) {
	ctx.SetColor(textColor)

	//display the start and end coords
	displayPointCoords(ctx, robotArm.getStartPtM(), 1400, fontSize)
	displayPointCoords(ctx, robotArm.getEndPtM(), 1400, 70+fontSize)

	//display the angle of the arm (combine with point and make into helper function)
	// drawFloat(ctx, robotArm.angle, 1400, 140+fontSize, "Angle Radians")
	startX := 1200.0
	drawFloat(ctx, robotArm.getAngleDeg(), startX, 140+fontSize, "Angle Degrees")
	drawFloat(ctx, robotArm.vel, startX, 210+fontSize, "Velocity Rad/s")
	drawFloat(ctx, robotArm.acc, startX, 280+fontSize, "Acceleration Rad/s")
	drawFloat(ctx, robotArm.voltage, startX, 350+fontSize, "Voltage volts")
} //end displayData
