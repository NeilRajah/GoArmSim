//armloop
//Author: Neil Balaskandarajah
//Created by: 10/22/2019
//State machine for the arm to follow

package main

import (
// "fmt"
// "time"
)

//Colors
var yellow [3]int = [3]int{255, 255, 0} //yellow
var blue [3]int = [3]int{0, 0, 255}     //blue

//State represents the state the arm can be in
type State int

const (
	waiting      State = iota //WAITING state is for a point to move to
	goalTracking              //GOAL_TRACKING state is for moving towards a point
	finished                  //arm has successfully moved to a point
)

//ArmLoop is the loop that controls the arm
type ArmLoop struct {
	arm2  Arm2  //arm to control
	goal  Point //goal point for arm to move to
	state State //state the arm is in
} //end struct

//get a string representation of the state
func (s State) String() string {
	return [...]string{"waiting", "goalTracking", "finished"}[s]
} //end String

//Set the state
//State s - new state for the arm to be in
func (loop *ArmLoop) setState(s State) {
	loop.state = s
} //end setState

func (loop *ArmLoop) onLoop() {
	switch loop.state {
	case waiting:
		loop.arm2.setArmColors(yellow) //yellow for waiting
		break

	case goalTracking:
		//proportional green for tracking
		color1 := loop.arm2.arm1.CalcColor(1)
		loop.arm2.arm1.color = color1
		color2 := loop.arm2.arm2.CalcColor(1)
		loop.arm2.arm2.color = color2

		//calculate joint angles
		a1, a2 := InverseKinematics(loop.goal, loop.arm2.arm1.angle, loop.arm2.arm2.angle,
			loop.arm2.arm1.length, loop.arm2.arm2.length)

		//move to joint angles
		loop.arm2.arm1.movePIDFF(a1, loop.arm2.arm1.angle, ToRadians(1))
		loop.arm2.arm2.movePIDFF(a2, loop.arm2.arm2.angle, ToRadians(1))

		//graphically update the arm
		loop.arm2.update()
		break

	case finished:
		loop.arm2.setArmColors(blue) //blue for finished

		// p := loop.goal
		// t := time.NewTimer(time.Millisecond * 1000)
		// go func() {
		// 	<-t.C
		// 	if loop.goal != p { //if goal is same after the delay
		// 		loop.state = finished
		// 	} else {
		// 		loop.state = goalTracking
		// 	} //if
		// 	return
		// }()
		break
	} //switch
} //end onLoop

//Set the goal point for the state machine
//Point p - new point to be the goal for the state machine
func (loop *ArmLoop) setGoal(p Point) {
	loop.goal = p
	loop.state = goalTracking
	loop.arm2.arm1.stopped = false
	loop.arm2.arm2.stopped = false
} //end setGoal
