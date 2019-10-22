//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	// "fmt"
	// "github.com/faiface/pixel"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"image/color"
	// "math"
	"time"
)

//Constants

const width int = 1920                //WIDTH is the width of the window
const height int = 1080               //HEIGHT is the height of the window
const fps int = 50                    //FPS is the frame rate of the animation
const dt float64 = 1.0 / float64(fps) //timestamp duration
const fontSize float64 = 60           //FONT_SIZE is the font size for the canvas

//variables
var c canvas.Canvas                         //canvas instance
var robotArm *Arm                           //arm struct
var robotArm2 Arm2                          //2-jointed arm
var bgColor color.RGBA = colornames.Black   //background color
var textColor color.RGBA = colornames.White //text color
var a1, a2 float64                          //angles to move to
var p Point                                 //point to move to

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
	c.Setup(func(ctx *canvas.Context) { setUpCanvas(ctx) })

	//create the arm
	createArm2()

	//create a starting point to move to
	p = Point{0.75, 1.25}
	a1, a2 = InverseKinematics(p, robotArm2.arm1.angle, robotArm2.arm2.angle, robotArm2.arm1.length, robotArm2.arm2.length)

	t := time.NewTimer(time.Millisecond * 1500)

	go func() {
		<-t.C
		// fmt.Println("Point switching")
		p = Point{-1.0, 1.2}
		a1, a2 = InverseKinematics(p, robotArm2.arm1.angle, robotArm2.arm2.angle, robotArm2.arm1.length, robotArm2.arm2.length)
	}()

	//create the arm
	c.Draw(func(ctx *canvas.Context) {
		go updateModel(ctx)
		draw(ctx)
	}) //end Draw
} //end main

//MODEL

//create the arm
func createArm() {
	//set the values for the arm
	kP := 1.25
	kI := 0.00 //0.01 for PID
	kD := 0.07
	robotArm = NewArm(1.0, 40.0, 159.3, 2, kP, kI, kD, "cim", 0)
} //end createArm

//create the two-jointed arm
func createArm2() {
	//PID constants
	kP1 := 2.00
	kI1 := 0.0
	kD1 := 0.04

	kP2 := 1.75
	kI2 := 0.0
	kD2 := 0.02

	//joints 1 and 2
	joint1 := NewArm(1.0, 30.0, 159.3, 2, kP1, kI1, kD1, "cim", 0)
	joint2 := NewArm(0.8, 15.0, 159.3, 1, kP2, kI2, kD2, "cim", 0)
	robotArm2.arm1 = joint1
	robotArm2.arm2 = joint2

	//set start of second joint to beginning of first joint
	robotArm2.arm2.start = robotArm2.arm1.getEndPtPxl()
} //end createArm2

//Update the arm for drawing purposes
func updateModel(ctx *canvas.Context) {
	//move with PID Feedback control and Feedforward until at target
	robotArm2.arm1.movePIDFF(a1, robotArm2.arm1.angle, ToRadians(1))
	robotArm2.arm2.movePIDFF(a2, robotArm2.arm2.angle, ToRadians(1))

	// fmt.Println(a1, a2)

	robotArm2.update() //update the arm
} //end updateModel

//DRAWING

//Update the model and draw to the scree
//ctx *canvas.Context - responsible for drawing
func draw(ctx *canvas.Context) {
	//clear the canvas
	ctx.SetColor(bgColor) //set the bg color
	ctx.Clear()           //empty the canvas

	drawCSpace(ctx) //draw the configuration space of the arm

	drawArm2(ctx) //draw the 2-jointed arm to the screen

	//display the data to the screen
	// displayData(ctx)

	//draw the point
	ctx.SetColor(colornames.White)
	drawPoint(ctx, p, 5)
} //end draw

//draw the robot to the display
//ctx *canvas.Context - responsible for drawing
func drawArm(ctx *canvas.Context) {
	ctx.Push() //save current state

	//draw the space the arm can be in (configuration space, or c-space)
	ctx.SetRGBA(1, 1, 1, 0.25) //switch to transparent green
	ctx.DrawCircle(robotArm.start.x, robotArm.start.y, robotArm.getLengthPxl()+armWidth/2)
	ctx.Fill() //fill the circle

	//draw the robot arm as lines between the joint points
	colors := robotArm.getColor(0)                 //arm color
	ctx.SetRGB255(colors[0], colors[1], colors[2]) //switch to the arm color
	ctx.SetLineWidth(armWidth)                     //change to the arm thickness
	ctx.DrawLine(robotArm.start.x, robotArm.start.y,
		robotArm.getEndPtPxl().x, robotArm.getEndPtPxl().y)
	ctx.Stroke() //draw the line

	ctx.Pop() //load last saved state
} //end drawRobot

//draw the two-jointed arm to the display
//ctx *canvas.Context - responsible for drawing
func drawArm2(ctx *canvas.Context) {
	ctx.Push() //save current state

	//draw the robot arm as lines between the joint points
	colors := robotArm2.arm1.getColor(0)           //red
	ctx.SetRGB255(colors[0], colors[1], colors[2]) //switch to the arm color
	ctx.SetLineWidth(armWidth)                     //change to the arm thickness
	ctx.DrawLine(robotArm2.arm1.start.x, robotArm2.arm1.start.y,
		robotArm2.arm1.getEndPtPxl().x, robotArm2.arm1.getEndPtPxl().y)
	ctx.Stroke()

	colors = robotArm2.arm2.getColor(0)            //red
	ctx.SetRGB255(colors[0], colors[1], colors[2]) //switch to the arm color
	ctx.SetLineWidth(armWidth)                     //change to the arm thickness
	ctx.DrawLine(robotArm2.arm2.start.x, robotArm2.arm2.start.y,
		robotArm2.arm2.get2JEndPtPxl(robotArm2.arm1.angle).x, robotArm2.arm2.get2JEndPtPxl(robotArm2.arm1.angle).y)
	ctx.Stroke() //draw the line

	ctx.Pop() //load last saved state
} //end drawArm2

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
	displayPointCoords(ctx, robotArm2.arm1.getStartPtM(), startX, fontSize)
	displayPointCoords(ctx, robotArm2.arm1.getEndPtM(), startX, 70+fontSize)
	displayPointCoords(ctx, robotArm2.arm2.getStartPtM(), startX, 140+fontSize)
	displayPointCoords(ctx, robotArm2.arm2.get2JEndPtM(robotArm2.arm1.angle), startX, 210+fontSize)

	//display the state + voltage of the arm
	drawFloat(ctx, robotArm2.arm1.getAngleDeg(), startX, 280+fontSize, "Angle Degrees")
	drawFloat(ctx, robotArm2.arm2.getAngleDeg(), startX, 350+fontSize, "Angle Degrees")
	// drawFloat(ctx, robotArm.vel, startX, 280+fontSize, "Velocity Rad/s")
	// drawFloat(ctx, robotArm.acc, startX, 350+fontSize, "Acceleration Rad/s")
	// drawFloat(ctx, robotArm.voltage, startX, 420+fontSize, "Voltage Volts")
} //end displayData

//draw the configuration space of the arm
//ctx *canvas.Context - responsible for drawing
func drawCSpace(ctx *canvas.Context) {
	ctx.Push()

	ctx.SetRGBA(0, 1, 0, 0.33)
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length+robotArm2.arm2.length)*pixelToMeters)
	ctx.Fill()

	ctx.SetColor(colornames.Black)
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length-robotArm2.arm2.length)*pixelToMeters)
	ctx.Fill()

	ctx.Pop()
} //end drawCSpace
