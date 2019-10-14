//arm
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a single-jointed arm in cartesian space

package main

import (
	// "github.com/faiface/pixel"
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
	start   Point   //the base point of the arm
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
//float64 angle - angle to start the arm at
func NewArm(length, mass, gearRatio, numMotors, kP, kI, kD float64, motorName string, angle float64) *Arm {
	//create the arm
	arm := new(Arm)

	//create and add all pre-determined values
	arm.start = Point{float64(width / 2), 0} //center of bottom edge of window
	arm.angle = angle                        //start at specified angle
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
	arm.kT = (numMotors * arm.motor.kStallTorque) / arm.motor.kStallCurrent //stall torque of whole arm (sum of all motor stall torques)

	//add and configure constants
	arm.maxVel = (arm.motor.kFreeSpeed / gearRatio) / 60 * 2 * math.Pi //radians per second
	arm.moi = 0.333333 * arm.mass * arm.length * arm.length            //moment of inertia

	return arm
} //end NewArm

//SETTERS AND GETTERS

//Get the end point of the arm in pixels
func (a Arm) getEndPtPxl() Point {
	endX := a.getLengthPxl()*math.Cos(a.angle) + a.start.x
	endY := a.getLengthPxl()*math.Sin(a.angle) + a.start.y

	return Point{endX, endY}
} //end getEndPt

//Get the end point of the arm in meters
func (a Arm) getEndPtM() Point {
	endX := a.length*math.Cos(a.angle) + a.getStartPtM().x
	endY := a.length*math.Sin(a.angle) + a.getStartPtM().y

	return Point{endX, endY}
} //end getEndPtM

//Get the end point of the arm if its a second joint
//float64 a1 - angle of the first joint of the arm in radians
func (a Arm) get2JEndPtPxl(a1 float64) Point {
	endX := a.getLengthPxl()*math.Cos(a1+a.angle) + a.start.x
	endY := a.getLengthPxl()*math.Sin(a1+a.angle) + a.start.y

	return Point{endX, endY}
} //end get2JEndPt

//Get the end point of the arm if its a second joint
//float64 a1 - angle of the first joint of the arm in radians
func (a Arm) get2JEndPtM(a1 float64) Point {
	endX := a.length*math.Cos(a1+a.angle) + a.getStartPtM().x
	endY := a.length*math.Sin(a1+a.angle) + a.getStartPtM().y

	return Point{endX, endY}
} //end get2JEndPt

//Get the start point of the arm in meters
func (a Arm) getStartPtM() Point {
	return Point{(a.start.x - float64(width)/2.0) / pixelToMeters, a.start.y / pixelToMeters}
} //end getStartPtM

//Set the start point of the arm as another point
//p point - point to set the arm to start at
func (a *Arm) setStartPt(p Point) {
	a.start = p
} //end setStartPt

//Get the angle of the arm in degrees
func (a Arm) getAngleDeg() float64 {
	return ToDegrees(a.angle)
} //end getAngleDeg

//Get the length of the arm in pixels
func (a Arm) getLengthPxl() float64 {
	return a.length * pixelToMeters
} //end getLengthPxl

//Set the angle of the arm in radians
//float64 newAngle - new angle for the arm
func (a *Arm) setAngle(newAngle float64) {
	a.angle = newAngle
} //end setAngle

//PHYSICS

//Calculate the torque caused on the arm by gravity
func (a Arm) calcGravTorque() float64 {
	return a.mass * g * a.length / 2 * math.Cos(a.angle) //mgrcosA
} //end calcGravTorque

//Calculate the current acceleration of the arm
//float64 output - output voltage to drive the arm
func (a *Arm) calcAccel(output float64) {
	voltConst := (a.gearRatio * a.kT) / (a.motor.kResistance * a.moi)                           //proportional to voltage
	velConst := (a.kT * a.gearRatio * a.gearRatio) / (a.motor.kV * a.motor.kResistance * a.moi) //proportional to velocity
	torqGrav := a.calcGravTorque()                                                              //mgrcosA

	a.acc = (output)*voltConst - a.vel*velConst - torqGrav/a.moi
} //end calcAccel

//MOTION

//set the voltage of the arm in a percentage of the max voltage
//float64 percent - percentage of the max voltage (between -1.0 and 1.0)
func (a *Arm) setOutput(percent float64) {
	a.voltage = MaxVoltage * percent
	a.update()
} //end setOutput

//drive the arm using PID control
func (a *Arm) movePID(setpoint, current, epsilon float64) {
	if robotArm.pid.atTarget && math.Abs(a.vel) < a.maxVel*0.1 { //if at target
		a.stopped = true
	} else { //if not
		a.voltage = MaxVoltage * OutputClamp(a.pid.calcPID(setpoint, current, epsilon), -1, 1)
	}
	a.update()
} //end movePID

//drive the arm using PIDFF control (PID + feedforward to hold arm)
//float64 setpoint - goal angle to move to
//float64 current - current angle of the arm
//float64 epsilon - tolerance for the angle in radians
func (a *Arm) movePIDFF(setpoint, current, epsilon float64) {
	if a.pid.atTarget && a.vel < a.maxVel*0.1 { //if at target
		a.stopped = false
	} else { //if not
		a.voltage = MaxVoltage*OutputClamp(a.pid.calcPID(setpoint, current, epsilon), -1, 1) + calcFFArm(a)
	}
	a.update()
} //end movePIDFF

//move the arm to the line formed by a goal point and origin (single-joint IK)
//Point goal - (x,y) point in meters
//float64 tolerance - tolerance for the angle in radians
func (a *Arm) pointToGoal(goal Point, tolerance float64) {
	angle := math.Atan2(goal.y, goal.x)
	a.movePIDFF(angle, a.angle, tolerance)
} //end pointToGoal

//UPDATE

//update the coordinates of the endpoint based on the angle
func (a *Arm) update() {
	a.voltage = OutputClamp(a.voltage, -12, 12)

	if a.stopped {
		// a.acc = 0
		// a.vel = 0
	} else {
		a.calcAccel(a.voltage)
		a.vel += a.acc * float64(1.0/float64(fps))
		a.moveArm(a.vel)
	}

	// fmt.Println(a.voltage)
} //end update

//move the arm with a given speed
//float64 angVel - the angular velocity of the arm in radians/second
func (a *Arm) moveArm(angVel float64) {
	dtheta := angVel * float64(1.0/float64(fps))
	a.angle += dtheta
} //end moveArm

//stop the arm by setting the velocity to zero
func (a *Arm) stop() {
	a.vel = 0
} //end stop

//GRAPHICS

//get the speed-proportional color for the arm
//int i - indicated whether R, G or B should be the color selected
func (a Arm) getColor(i int) [3]int {
	color := [3]int{0, 0, 0}
	color[i] = int(OutputClamp((math.Abs(a.vel)/a.maxVel)*127, 0, 127) + 127)
	if a.stopped {
		color = [3]int{0, 0, 255} //blue
	}
	return color
} //end getColor
