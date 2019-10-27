//main
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//Main file that handles all of the graphics and the model of the arm

package main

import (
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
var canAdd bool   //whether a point can be added by clicking to the set or not
var canTrack bool //whether the arm can track its goal or not
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
	//scale the point from pixel to cartesian coordinates
	ghost = scalePoint(Point{ctx.Mouse.X - float64(width)/2, ctx.Mouse.Y}, 1.0/float64(pixelToMeters))
	//clamp the point to the configuration space of the arm
	ghost = ClampToCSpace(ghost, robotArm2.arm1.length, robotArm2.arm2.length)

	//add points with mouse click
	if canAdd { //if user can add points
		if ctx.IsMouseDragged { //mouse click
			pts = append(pts, ghost) //add it to the list of goals

			//canAdd and the Timer are used to prevent multiple points be added during one click
			canAdd = false
			t = time.NewTimer(time.Millisecond * time.Duration(250))
		} //if
	} else {
		//wait for timer's duration until the user can add points again
		go func() {
			<-t.C         //blocks until finished
			canAdd = true //let user add a point again
		}() //timer goroutine
	} //if

	//set goal
	if len(pts) != 0 { //if lists isn't empty
		if pts[pointIndex] != armloop.goal && armloop.state != goalTracking { //if last point in list isn't already the goal and the arm is finished
			time.Sleep(time.Millisecond * 250) //delay before setting goal
			armloop.setGoal(pts[pointIndex])   //set the next point as the arm's goal
		} //if
	} //if
} //end updateGoal

//Update the arm's state machine
func updateModel() {
	//update the state for the state machine
	if robotArm2.isStopped() { //if both joints are stopped
		armloop.setState(finished)   //set the state to finished
		if len(pts)-1 > pointIndex { //if there is another point to move to
			pointIndex++
		} //if
	} //if

	armloop.onLoop() //move the arm
} //end updateModel

//DRAWING

//Draw to the scree
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
