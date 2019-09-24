//arm
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a single-jointed arm in cartesian space

package main

import (
	// "fmt"
	"math"
)

//A single-jointed arm in cartesian space
type Arm struct {
	start  point   //the base point of the arm
	length float64 //the length of the arm
	angle  float64 //the angle of the arm from the horizontal measured CCW in radians
	vel    float64 //angular velocity of the arm in radians/second
} //end struct

//SETTERS AND GETTERS

//set the angular velocity of the arm in degrees per second
//float64 newVel - the new angular velocity of the arm in degrees per second
func (a *Arm) setVelDPS(newVel float64) {
	a.vel = newVel
} //end setVel

//Get the end point of the arm
func (a Arm) getEndPt() point {
	endX := a.length*math.Cos(a.angle) + a.start.x
	endY := a.length*math.Sin(a.angle) + a.start.y

	return point{endX, endY}
} //end GetEndPt

//MOTION

//move the arm with a given speed
//float64 angVel - the angular velocity of the arm in degrees/second
func (a *Arm) moveArm(angVel float64) {
	dtheta := angVel * float64(1.0/float64(FPS))
	a.angle = a.angle + ToRadians(dtheta)
} //end moveArm

//UPDATE

//update the coordinates of the endpoint based on the angle
func (a *Arm) update() {
	a.moveArm(a.vel)
} //end update
