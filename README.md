# GoArmSim
GoArmSim is an interactive program that allows the user to give Cartesian coordinate points to a 2-jointed robotic arm and have it move to them using a simple feedback controller and a (somewhat realistic) physics model.

## User Interface
- canvas
- two libraries
- click
- ghost
- configuration space (explained later)
- state display (text, arm colors)

## Motor Model
Motor modelling formulae derived from this book: https://www.chiefdelphi.com/t/paper-practical-guide-to-state-space-control/166417/2
- VEX link for CIM motor model
- how constants were calculated

~ graph for CIM motor

## Arm Model
- 2D, 2DoF simple two-joint that can pass through itself
- base joint two motors with gear reduction, elbow joint two motors with gear reduction
- values for joint lengths, masses
- configuration space (Desmos image)
  - solving inequality to clamp point inside of c-space
     - shrinking/enlargening value by small factor to avoid rounding errors

## Inverse Kinematics of the Arm
Inverse Kinematics formulae derived from this video: https://robotacademy.net.au/lesson/inverse-kinematics-for-a-2-joint-robot-arm-using-geometry/

A geometric approach was chosen for solving the inverse kinematics of the arm due to its relative simplicity. Since the arm only has two joints, the algebraic approach involving complex matrix math would be overcomplicated, and a numerical approach where the solution is found using optimization less effective. The calculations involved in the geometric approach boil down to basic trigonometry after a few key triangles are formulated.

~ IK1
~ IK2

Because of the symmetry of the cosine function, there are two solutions to the joint angles. These configurations are referred to as "Elbow Up" and "Elbow Down".

~ IK3

To decide between the two solutions, the robot operates under the assumption that the end-effector must "face" the goal point. This is done by choosing the set with the negative elbow joint angle (Elbow Down) in quadrant one and positive elbow joint angle (Elbow Up) in quadrant two. This prevents the arm from moving below the y-axis and into the ground when moving between points. 

## Dynamics Model
- gravity
- torque provided by motor
- center of mass, moment of inertia
- discrete calculations
- "integration"

## Feedback Controller
- PID
- PIDFF
   - how F eliminates need for I (I gain is zero)

## State Machine
- waiting at beginning
- tracks goal
- finished
   - delay before given next point

## Potential Improvements
- better decision for pairs (ie. edges of inner c-space)
- coordinating joint movement to avoid ground (running one after another, slowing one down)
- joint space <--> task space
   - splines in joint space
   - thrust control
- coupled physics model instead of separate arms
- change motor configuration, gearbox ratio, feedback gains through file/ on screen
- timer to show duration of movement
- going into ERROR state when given point out of c-space instead of clamping inside
