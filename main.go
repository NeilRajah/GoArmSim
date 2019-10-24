//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
	"fmt"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

//Constants

const width int = 1920                //WIDTH is the width of the window
const height int = 1080               //HEIGHT is the height of the window
const fps int = 50                    //FPS is the frame rate of the animation
const dt float64 = 1.0 / float64(fps) //timestamp duration
const fontSize float64 = 60           //FONT_SIZE is the font size for the canvas

//variables
var bgColor color.RGBA = colornames.Black            //background color
var textColor color.RGBA = colornames.White          //text color
var cspaceColor []float64 = []float64{0, 1, 0, 0.25} //configuration space color
var ghostColor []float64 = []float64{1, 0, 0, 0.25}  //ghost point color
var pointIndex int = 0                               //increasing index for drawing points
var goalIndex int = 0                                //increasing index for setting goal point
var c canvas.Canvas                                  //canvas instance
var ghost Point                                      //ghost point to draw

var robotArm *Arm   //arm struct
var robotArm2 Arm2  //2-jointed arm
var armloop ArmLoop //state machine for the arm

var pts []Point   //points to move to
var canAdd bool   //whether a point can be added to the set or not
var t *time.Timer //timer for the arm

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

	canAdd = true //can add mouse points

	//draw to the canvas
	c.Draw(func(ctx *canvas.Context) {
		updateGoal(ctx)
		updateModel()
		draw(ctx)
	})
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

	//state machine for the arm
	armloop = ArmLoop{arm2: robotArm2, state: waiting}
} //end createArm2

//add the mouse click coordinates as points for the arm
//*canvas.Context ctx - used for drawing
func updateGoal(ctx *canvas.Context) {
	ghost = scalePoint(Point{ctx.Mouse.X - float64(width)/2, ctx.Mouse.Y}, 1.0/float64(pixelToMeters))
	ghost = ClampToCSpace(ghost, robotArm2.arm1.length, robotArm2.arm2.length) //clamp it to c-space

	//add points with mouse click
	if canAdd { //if user can add points
		if ctx.IsMouseDragged { //mouse click
			pts = append(pts, ghost) //add it to the list of goals
			fmt.Println("goal:", ghost.x, ghost.y)

			//canAdd and the Timer are used to prevent multiple points be added during one click
			canAdd = false
			t = time.NewTimer(time.Millisecond * time.Duration(250))
		} //if
	} else {
		//wait for timer's duration until the user can add points again
		go func() {
			<-t.C         //blocks until finished
			canAdd = true //let user add a point again
		}()
	} //if

	//set goal
	if len(pts) != 0 { //if lists isn't empty
		if pts[pointIndex] != armloop.goal && armloop.state != goalTracking { //if last point in list isn't already the goal and the arm is finished
			armloop.setGoal(pts[pointIndex]) //set the next point as the arm's goal
		} //if
	} //if
} //end updateGoal

//Update the arm for drawing purposes
func updateModel() {
	//update the state for the state machine
	if robotArm2.isStopped() { //if both joints are stopped and there is another goal
		armloop.setState(finished)   //set the state to finished
		if len(pts)-1 > pointIndex { //if there is another point to move to
			pointIndex++
		} //if
	} else if pointIndex > len(pts) { //if there is a point in sequence
		armloop.setState(goalTracking) //start moving to the point
	} //if

	armloop.onLoop() //move the arm
} //end updateModel

//DRAWING

//Update the model and draw to the scree
//ctx *canvas.Context - responsible for drawing
func draw(ctx *canvas.Context) {
	//clear the canvas
	ctx.SetColor(bgColor) //set the bg color
	ctx.Clear()           //empty the canvas

	drawCSpace(ctx)  //draw the configuration space of the arm
	drawPoints(ctx)  //draw all the points the robot can move to
	drawGhost(ctx)   //draw a point based on mouse location to show potential goal
	displayData(ctx) //display the data to the screen
	drawArm2(ctx)    //draw the 2-jointed arm to the screen
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

//draw the configuration space of the arm
//ctx *canvas.Context - responsible for drawing
func drawCSpace(ctx *canvas.Context) {
	ctx.Push()

	ctx.SetRGBA(cspaceColor[0], cspaceColor[1], cspaceColor[2], cspaceColor[3])                      //transparent green
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length+robotArm2.arm2.length)*pixelToMeters) //outer limit
	ctx.Fill()

	ctx.SetColor(bgColor)
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length-robotArm2.arm2.length)*pixelToMeters) //inside limit
	ctx.Fill()

	ctx.Pop()
} //end drawCSpace

//Draw a point around where the mouse is to show where the potential goal would be
//*canvas.Context ctx - responsible for drawing
func drawGhost(ctx *canvas.Context) {
	ctx.Push()

	ctx.SetRGBA(ghostColor[0], ghostColor[1], ghostColor[2], ghostColor[3]) //transparent red

	// p := scalePoint(ghost, 1.0/pixelToMeters)
	// p = ClampToCSpace(p, robotArm2.arm1.length, robotArm2.arm2.length)
	// drawPoint(ctx, p, 30)
	drawPoint(ctx, ghost, 30)

	ctx.Pop()
} //end drawGhost

//Draw all the points in the list the robot can move to
//ctx *canvas.Context - responsible for drawing
func drawPoints(ctx *canvas.Context) {
	//iterate through list and draw points
	start := pointIndex                                    //don't draw the points the arm has moved to
	if len(pts)-1 == pointIndex && robotArm2.isStopped() { //if the arm is stopped and there is no goal
		start++ //don't draw the point its at as well
	} //if

	factor := 0.0
	if len(pts) == start {
		factor = 0.0
	} else {
		factor = 0.75 / float64((len(pts) - start)) //multiply factor by index
	}

	//loop through all the points to move to and draw them
	for i := start; i < len(pts); i++ {
		// val := float64(255 - factor*(i-start))
		val := 1.0 - (factor * float64((i - start)))
		ctx.SetRGBA(1, 0, 0, val)
		drawPoint(ctx, pts[i], 30)
	} //loop
} //end drawPoints

//display the parameters of the robot onto the screen
//ctx *canvas.Context - responsible for drawing
func displayData(ctx *canvas.Context) {
	ctx.Push()

	ctx.InvertY()
	//change text color based on state
	switch armloop.state {
	case waiting:
		ctx.SetColor(colornames.Yellow)
		break
	case goalTracking:
		ctx.SetColor(colornames.Green)
		break
	case finished:
		ctx.SetColor(colornames.Blue)
		break
	} //switch
	ctx.DrawString(armloop.state.String(), 1400, 200)
	ctx.InvertY()

	ctx.Pop()
} //end displayData
