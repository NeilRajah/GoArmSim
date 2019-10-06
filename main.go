//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	"image/color"

	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	// "math"
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
		updateModel(ctx)

		//clear the canvas
		ctx.SetColor(bgColor) //set the bg color
		ctx.Clear()           //empty the canvas

		drawRobot(ctx) //draw the robot to the screen

		//display the data to the screen
		displayData(ctx)
	})

} //end main

//create the arm
func createArm() {
	//set the values for the arm
	kP := 1.25
	kI := 0.00 //0.01 for PID
	kD := 0.07
	robotArm = NewArm(1.0, 40.0, 159.3, 2, kP, kI, kD, "cim", 0)
} //end createArm

//Update the arm for drawing purposes
func updateModel(ctx *canvas.Context) {
	//move with PID control until the target is reached
	// robotArm.movePID(ToRadians(135), robotArm.angle, ToRadians(1))
	// robotArm.movePIDFF(ToRadians(135), robotArm.angle, ToRadians(1))
	// robotArm.voltage = calcFFArm(robotArm)
	// robotArm.update()
	goal := Point{1, 1}
	rad := 100.0
	ctx.SetColor(colornames.White)
	ctx.DrawCircle(float64(width)-goal.x*pixelToMeters-rad, goal.y*pixelToMeters+rad, rad)
	ctx.Fill()
	robotArm.pointToGoal(goal, ToRadians(1))
} //end updateModel

//SIMULATOR

//draw the robot to the display
//ctx *canvas.Context - responsible for drawing
func drawRobot(ctx *canvas.Context) {
	ctx.Push() //save current state

	//draw the space the arm can be in (configuration space, or c-space)
	ctx.SetRGBA(1, 1, 1, 0.25) //switch to transparent green
	ctx.DrawCircle(robotArm.start.x, robotArm.start.y, robotArm.getLengthPxl()+armWidth/2)
	ctx.Fill() //fill the circle

	//draw the robot arm as lines between the joint points
	colors := robotArm.getColor()                  //arm color
	ctx.SetRGB255(colors[0], colors[1], colors[2]) //switch to the arm color
	ctx.SetLineWidth(armWidth)                     //change to the arm thickness
	ctx.DrawLine(robotArm.start.x, robotArm.start.y,
		robotArm.getEndPtPxl().x, robotArm.getEndPtPxl().y)
	ctx.Stroke() //draw the line

	ctx.Pop() //load last saved state
} //end drawRobot

//display the parameters of the robot onto the screen
//ctx *canvas.Context - responsible for drawing
func displayData(ctx *canvas.Context) {
	startX := 1200.0 //starting x-coordinate for the text
	ctx.SetColor(textColor)

	//draw corner frame
	y := 350.0
	ctx.Push()
	ctx.InvertY()
	ctx.SetLineWidth(boxThick)
	ctx.DrawLine(startX-20, 0, startX-20, y+fontSize+20)
	ctx.DrawLine(startX-20, y+fontSize+20, float64(width), y+fontSize+20)
	ctx.Stroke()
	ctx.Pop()

	//display the start and end coords
	displayPointCoords(ctx, robotArm.getStartPtM(), startX, fontSize)
	displayPointCoords(ctx, robotArm.getEndPtM(), startX, 70+fontSize)

	//display the state + voltage of the arm
	drawFloat(ctx, robotArm.getAngleDeg(), startX, 140+fontSize, "Angle Degrees")
	drawFloat(ctx, robotArm.vel, startX, 210+fontSize, "Velocity Rad/s")
	drawFloat(ctx, robotArm.acc, startX, 280+fontSize, "Acceleration Rad/s")
	drawFloat(ctx, robotArm.voltage, startX, 350+fontSize, "Voltage Volts")
} //end displayData
