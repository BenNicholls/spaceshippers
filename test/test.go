package burltest

// import (
// 	"fmt"
// 	"github.com/bennicholls/burl/util"
// 	"math"
// )

// func main() {
// 	curveSim()
// }

// //Simulation that tests the 1-D solution to see if burn times would be accurate for sub-max fuel courses.
// func fuelSim() {
// 	D := 1.4876516e11
// 	v_i := 0.0
// 	v_f := 500000.0
// 	T := 10.0

// 	target_fuel := 2*math.Sqrt((v_f*v_f+v_i*v_i)/2+T*D)/T - (v_f+v_i)/T
// 	fmt.Println("Fuel to use:", target_fuel)

// 	t_c := 0.0
// 	for F := 0.0; t_c >= 0; F += 3000.0 {
// 		t_c = (2*D + (v_f*v_f+v_i*v_i)/(2*T) - v_f*v_i/T - F*(v_f+v_i) - T*F*F/2) / (v_f + v_i + T*F)
// 		t_d := math.Sqrt((v_f*v_f+v_i*v_i)/2+T*T*t_c*t_c/4+T*D)/T - v_f/T - t_c/2
// 		t_a := math.Sqrt((v_f*v_f+v_i*v_i)/2+T*T*t_c*t_c/4+T*D)/T - v_i/T - t_c/2
// 		fmt.Println("Fuel:", F)
// 		fmt.Println("Coat Time:", t_c)
// 		fmt.Println("Accel Time:", t_a)
// 		fmt.Println("Decel Time:", t_d)
// 	}
// }

// //Newtonian Turning Simulation. Ended up not working, saved for posterity??
// func curveSim() {

// 	pos := util.Vec2{-1.49e11, 0}
// 	v := util.Vec2Polar{1000000.0, 0.5}

// 	target := util.Vec2{0, 0}

// 	targetCourse := target.Sub(pos).ToPolar()
// 	a := 10.0
// 	turnCircleCenter := pos.ToPolar().Add(util.Vec2Polar{v.R * v.R / a, v.Phi - math.Pi/2}).ToRect()

// 	t := 0
// 	for ; math.Abs(v.AngularDistance(targetCourse)) > math.Abs(math.Atan2(a, v.R)); t++ {

// 		v = v.Add(util.Vec2Polar{a, v.Phi - math.Pi/2})

// 		pos = pos.Add(v.ToRect())
// 		targetCourse = target.Sub(pos).ToPolar()
// 	}

// 	fmt.Println("Correction took: ", t, "s")
// 	fmt.Println("New angle: ", v.Phi)
// 	fmt.Println("Angular Diff", math.Abs(v.AngularDistance(targetCourse)))
// 	fmt.Println("New pos", pos.X, pos.Y)
// 	fmt.Println("Distance after correction:", targetCourse.R)

// 	fmt.Println("Circle Center", turnCircleCenter.X, turnCircleCenter.Y)
// 	fmt.Println("final speed", v.R)
// 	fmt.Println("NewPos On Circle?", turnCircleCenter.Sub(pos).Mag()-v.R*v.R/a)

// 	fmt.Println("")
// }
