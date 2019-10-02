//motor
//Author: Neil Balaskandarajah
//Created on: 10/01/2019
//Model of a simple DC electric motor for simulations

package main

import (
	"math"
	"strings"
)

//Constants

//MaxVoltage is the maximum voltage of the robot in Volts
const MaxVoltage = 12.0

//Motor struct is an instance of a simple DC motor
type Motor struct {
	//configured
	kStallTorque  float64 //stall torque in Nm
	kStallCurrent float64 //stall current in Amps
	kFreeSpeed    float64 //free speed in RPM
	kFreeCurrent  float64 //free current in Amps

	//calculated
	kResistance float64 //resistance in Ohms
	kV          float64 //velocity constant of the motor
} //end struct

//NewMotor returns a newly configured motor pointer based on its name
//string motorName - the name of the motor
func NewMotor(motorName string) *Motor {
	m := new(Motor)

	if strings.EqualFold(motorName, "cim") { //CIM Motor
		m.kStallTorque = 2.42
		m.kStallCurrent = 133
		m.kFreeSpeed = 5330
		m.kFreeCurrent = 2.7

		m.kResistance = MaxVoltage / m.kStallCurrent
		m.kV = (m.kFreeSpeed / 60 * 2 * math.Pi) / (MaxVoltage - m.kResistance*m.kFreeCurrent)
	} //if

	return m
} //end new motor

//MakeMotor returns a newly configured motor based on its name
//string motorName - name of the motor
func MakeMotor(motorName string) Motor {
	if strings.EqualFold(motorName, "cim") {
		resistance := MaxVoltage / 133
		kV := (5330 / 60 * 2 * math.Pi) / (MaxVoltage - (MaxVoltage/133)*2.7)
		return Motor{2.42, 133, 5330, 2.7, resistance, kV}
	} //if
	return Motor{}
} //end MakeMotor
