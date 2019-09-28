//arm
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a single-jointed arm in cartesian space

package main

import (
	"fmt"
	"math"
)

//Constants
const PixelToInches = 30 //the number of pixels per 'inch' in the GUI

//A single-jointed arm in cartesian space
type Arm struct {
	//Configured attributes
	start    point   //the base point of the arm
	length   float64 //the length of the arm
	angle    float64 //the angle of the arm from the horizontal measured CCW in radians
	vel      float64 //angular velocity of the arm in radians/second
	topSpeed float64 //top speed of the arm in deg/s

	//Calculated attributes
	pid pidcontroller //PID controller for the arm
} //end struct

//SETTERS AND GETTERS

//set the PID controller for the arm
//pidcontroller pid - PID controller with configurable parameters
func (a *Arm) setPIDController(pid pidcontroller) {
	a.pid = pid
} //end setPIDController

//Get the end point of the arm in pixels
func (a Arm) getEndPtPxl() point {
	endX := a.getLengthPxl()*math.Cos(a.angle) + a.start.x
	endY := a.getLengthPxl()*math.Sin(a.angle) + a.start.y

	return point{endX, endY}
} //end getEndPt

//Get the end point of the arm in inches
func (a Arm) getEndPtIn() point {
	endX := a.length * math.Cos(a.angle)
	endY := a.length * math.Sin(a.angle)

	return point{endX, endY}
} //end getEndPt

//get the start point of the arm in inches
//Get the end point of the arm in inches
func (a Arm) getStartPtIn() point {
	return point{0, 0}
} //end getEndPt

//Get the angle of the arm in degrees
func (a Arm) getAngleDeg() float64 {
	return ToDegrees(a.angle)
} //end getAngleDeg

//Get the length of the arm in pixels
func (a Arm) getLengthPxl() float64 {
	return a.length * PixelToInches
} //end getLengthPxl

//MOTION

//set the angular velocity of the arm in degrees per second
//float64 newVel - the new angular velocity of the arm in degrees per second
func (a *Arm) setVelDPS(newVel float64) {
	a.vel = newVel
	a.update()
} //end setVel

//set the vel of the arm in a percentage of the top speed
//float64 percent - percentage of the top speed for velocity (between -1.0 and 1.0)
func (a *Arm) setOutput(percent float64) {
	a.vel = percent * a.topSpeed
	a.update()
} //end setOutput

//drive the arm using PID control
func (a *Arm) movePID(setpoint, current, epsilon float64) {
	if !robotArm.pid.atTarget {
		a.vel = OutputClamp(a.pid.calcPID(setpoint, current, epsilon), -a.topSpeed, a.topSpeed)
		if a.vel > a.topSpeed {
			fmt.Println(a.vel)
		}
	} else {
		a.vel = 0
	}
	a.update()
} //end movePID

//UPDATE

//update the coordinates of the endpoint based on the angle
func (a *Arm) update() {
	a.moveArm(a.vel)
} //end update

//move the arm with a given speed
//float64 angVel - the angular velocity of the arm in degrees/second
func (a *Arm) moveArm(angVel float64) {
	dtheta := angVel * float64(1.0/float64(fps))
	a.angle = a.angle + ToRadians(dtheta)
} //end moveArm

//stop the arm by setting the velocity to zero
func (a *Arm) stop() {
	a.vel = 0
} //end stop

//GRAPHICS

//get the speed-proportional color for the arm
func (a Arm) getColor() [3]int {
	if a.vel == 0 {
		return [3]int{255, 0, 0}
	} else {
		c := [3]int{0, int((a.vel/a.topSpeed)*127 + 127), 0}
		return c
	} //if
} //end getColor
