//arm
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a single-jointed arm in cartesian space

package main

import (
	// "fmt"
	"math"
)

//Arm is a single-jointed arm in cartesian space
type Arm struct {
	//Configured attributes
	length    float64 //the length of the arm in meters
	mass      float64 //mass of the arm in kg
	gearRatio float64 //gear ratio of the gearbox powering the arm

	//Calculated attributes
	start   point   //the base point of the arm
	angle   float64 //the angle of the arm from the horizontal measured CCW in radians
	vel     float64 //angular velocity of the arm in radians/second
	maxVel  float64 //maximum possible velocity of the arm
	acc     float64 //angular acceleration of the arm in radians/second^2
	moi     float64 //moment of inertia of the arm
	voltage float64 //current voltage being output

	numMotors float64 //number of motors powering the arm
	kT        float64 //torque constant of the arm

	pid   pidcontroller //PID controller for the arm
	motor Motor         //motor controlling the arm

	stopped bool //whether the arm is stopped or not
} //end struct

//NewArm creates a new arm given configurable parameters
//float64 length - length of the arm in meters
//float64 mass - mass of the arm in kg
//float64 numMotors - number of motors powering the arm
//float64 kP - proportionality constant
//float64 kI - integral constant
//float64 kD - derivative constant
//pidcontroller pid - calculates PID outputs
//string motorName - name of the motor
func NewArm(length, mass, gearRatio, numMotors, kP, kI, kD float64, motorName string) *Arm {
	//create the arm
	arm := new(Arm)

	//create and add all pre-determined values
	arm.start = point{float64(width / 2), 0} //center of bottom edge of window
	arm.angle = 0                            //start at horizontal
	arm.vel = 0                              //start at rest
	arm.acc = 0                              //start with no acceleration
	arm.voltage = 0

	//add all passed values
	arm.length = length
	arm.mass = mass
	arm.gearRatio = gearRatio
	arm.numMotors = numMotors

	//create and configure PID controller
	arm.pid = pidcontroller{kP: kP, kI: kI, kD: kD}

	//create and configure motor
	arm.motor = MakeMotor(motorName)
	arm.kT = (numMotors * arm.motor.kStallTorque) / arm.motor.kStallCurrent

	//add and configure constants
	arm.maxVel = (arm.motor.kFreeSpeed / gearRatio) / 60 * 2 * math.Pi //radians per second
	arm.moi = 0.333333 * arm.mass * arm.length * arm.length

	return arm
} //end NewArm

//SETTERS AND GETTERS

//Get the end point of the arm in pixels
func (a Arm) getEndPtPxl() point {
	endX := a.getLengthPxl()*math.Cos(a.angle) + a.start.x
	endY := a.getLengthPxl()*math.Sin(a.angle) + a.start.y

	return point{endX, endY}
} //end getEndPt

//Get the end point of the arm in meters
func (a Arm) getEndPtM() point {
	endX := a.length * math.Cos(a.angle)
	endY := a.length * math.Sin(a.angle)

	return point{endX, endY}
} //end getEndPtM

//Get the start point of the arm in meters
func (a Arm) getStartPtM() point {
	return point{0, 0}
} //end getStartPtM

//Get the angle of the arm in degrees
func (a Arm) getAngleDeg() float64 {
	return ToDegrees(a.angle)
} //end getAngleDeg

//Get the length of the arm in pixels
func (a Arm) getLengthPxl() float64 {
	return a.length * pixelToMeters
} //end getLengthPxl

//PHYSICS

//Calculate the torque caused on the arm by gravity
func (a Arm) calcGravTorque() float64 {
	return a.mass * g * a.length / 2 * math.Cos(a.angle) //mgrcosA
} //end calcGravTorque

//Calculate the current acceleration of the arm
func (a *Arm) calcAccel(output float64) {
	voltConst := (a.gearRatio * a.kT) / (a.motor.kResistance * a.moi)
	velConst := (a.kT * a.gearRatio * a.gearRatio) / (a.motor.kV * a.motor.kResistance * a.moi)

	a.acc = (output)*voltConst - a.vel*velConst - (a.mass * g * a.length / 2 * math.Cos(a.angle))
} //end calcAccel

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
	a.voltage = MaxVoltage * percent
	a.update()
} //end setOutput

//drive the arm using PID control
func (a *Arm) movePID(setpoint, current, epsilon float64) {
	if !robotArm.pid.atTarget { //if not at target
		a.voltage = MaxVoltage * OutputClamp(a.pid.calcPID(setpoint, current, epsilon), -a.maxVel, a.maxVel)
	} else { //if at target
		a.stopped = true
	}
	a.update()
} //end movePID

//UPDATE

//update the coordinates of the endpoint based on the angle
func (a *Arm) update() {
	if a.stopped {
		a.acc = 0
		a.vel = 0
	} else {
		a.calcAccel(a.voltage)
		a.vel += a.acc * float64(1.0/float64(fps))
		a.moveArm(a.vel)

		if a.angle < 0 {
			a.angle = 0
		} else if a.angle > math.Pi {
			a.angle = math.Pi
		}
	}
} //end update

//move the arm with a given speed
//float64 angVel - the angular velocity of the arm in radians/second
func (a *Arm) moveArm(angVel float64) {
	dtheta := angVel * float64(1.0/float64(fps))
	a.angle = a.angle + dtheta
} //end moveArm

//stop the arm by setting the velocity to zero
func (a *Arm) stop() {
	a.vel = 0
} //end stop

//GRAPHICS

//get the speed-proportional color for the arm
func (a Arm) getColor() [3]int {
	c := [3]int{0, int((a.vel/a.maxVel)*164 + 90), 0}
	return c
} //end getColor
