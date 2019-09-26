//pidcontroller
//Author: Neil Balaskandarajah
//Created on: 09/25/2019
//A simple PID controller that assumes regular loop intervals

package main

import ()

//pid controller struct with the three gains
type pidcontroller struct {
	kP        float64 //proportionality constant
	kI        float64 //integral constant
	kD        float64 //derivative constant
	errorSum  float64 //sum of all errors
	lastError float64 //last error for derivative calculation
} //end struct

//calculate the PID output based on the setpoint, current value and tolerance
//float64 setpoint - the desired goal value
//float64 currentValue - the current value
//float64 tolerance - the range you can be in to be considered "at goal"
func (pid pidcontroller) calcPID(setpoint, currentVal, tolerance float64) float64 {
	//get the error
	error := setpoint - currentVal

	//P value
	pOut := pid.kP * error //output proportional to error

	//I value
	pid.errorSum += error //add onto the error sum
	iOut := pid.kI * pid.errorSum

	//D value
	dError := (error - pid.lastError)
	dOut := pid.kD * dError

	return pOut + iOut + dOut
} //end calcPID
