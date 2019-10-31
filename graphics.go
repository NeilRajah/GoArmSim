//graphics
//Author: Neil Balaskandarajah
//Created on: 09/25/2019
//Graphics functions used for the simulator

package main

import (
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"strconv"
)

const armWidth = 30.0 //thickness of the arm lines in pixels

//set up the canvas
//*canvas.Context ctx - responsible for drawing
func setUpCanvas(ctx *canvas.Context) {
	ctx.SetColor(colornames.Lightgray)                           //set the bg color
	ctx.Clear()                                                  //empty the canvas
	ctx.SetColor(colornames.Black)                               //set the drawing color
	ctx.SetLineWidth(30)                                         //set the line width
	ctx.LoadFontFace("resources/HelveticaNeue.ttf", fontSize) //set the font
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
	ctx.Fill()

	ctx.Pop()
} //end drawPoint

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
	ctx.Push()                 //save current state
	ctx.SetLineWidth(armWidth) //change to the arm thickness

	//draw the robot arm as lines between the joint points
	colors := robotArm2.arm1.getColor(0) //switch to the joint color
	ctx.SetRGB255(colors[0], colors[1], colors[2])
	ctx.DrawLine(robotArm2.arm1.start.x, robotArm2.arm1.start.y,
		robotArm2.arm1.getEndPtPxl().x, robotArm2.arm1.getEndPtPxl().y)
	ctx.Stroke() //draw the line

	colors = robotArm2.arm2.getColor(0) //switch to the joint color
	ctx.SetRGB255(colors[0], colors[1], colors[2])
	ctx.DrawLine(robotArm2.arm2.start.x, robotArm2.arm2.start.y,
		robotArm2.arm2.get2JEndPtPxl(robotArm2.arm1.angle).x, robotArm2.arm2.get2JEndPtPxl(robotArm2.arm1.angle).y)
	ctx.Stroke() //draw the line

	ctx.Pop() //load last saved state
} //end drawArm2

//draw the configuration space of the arm
//ctx *canvas.Context - responsible for drawing
func drawCSpace(ctx *canvas.Context) {
	ctx.Push()

	ctx.SetRGBA(cspaceColor[0], cspaceColor[1], cspaceColor[2], cspaceColor[3])                      //c-space color
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length+robotArm2.arm2.length)*pixelToMeters) //outer limit
	ctx.Fill()

	ctx.SetColor(bgColor)                                                                            //background color
	ctx.DrawCircle(float64(width)/2, 0, (robotArm2.arm1.length-robotArm2.arm2.length)*pixelToMeters) //inside limit
	ctx.Fill()

	ctx.Pop()
} //end drawCSpace

//Draw a point around where the mouse is to show where the potential goal would be
//*canvas.Context ctx - responsible for drawing
func drawGhost(ctx *canvas.Context) {
	ctx.Push()

	ctx.SetRGBA(ghostColor[0], ghostColor[1], ghostColor[2], ghostColor[3]) //transparent red

	drawPoint(ctx, ghost, 30) //draw the ghost point

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
	} //if

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
