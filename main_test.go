//main_test
//Author: Neil Balaskandarajah
//Created on: 10/12/2019
//Testing the inverse kinematics of the two jointed arm

package main

import (
	"math"
	"testing"
)

//if the forward kinematics produced with the inverse kinematics angles is not within a tolerance, fail the test
func TestIK(t *testing.T) {
	target := Point{0.375, 1.0}
	a1, a2 := InverseKinematics(target, 1.0, 0.8)
	p2 := forwardKinematics(1.0, 0.8, a1, a2)

	t.Log("Forward kinematics produced: (a1, a2, goal point)", ToDegrees(a1), ToDegrees(a2), p2.x, p2.y)

	if !withinBounds(target, p2, 1.0) {
		t.Error("Angles do not produce point")
	}
} //end testIK

//Calculate the coordinates of the two arm joints given their angles and their lengths
//float64 a1 - length of the first arm
//float64 a2 - length of the second arm
//float64 q1 - angle of the first joint
//float64 q2 - angle of the second joint
func forwardKinematics(a1, a2, q1, q2 float64) Point {
	p1 := Point{a1 * math.Cos(q1), a1 * math.Sin(q1)}
	p2 := Point{p1.x + a2*math.Cos(q1+q2), p1.y + a2*math.Sin(q1+q2)}

	return p2
} //end forward kinematics
