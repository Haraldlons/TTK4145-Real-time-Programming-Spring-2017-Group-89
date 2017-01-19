package main

import (
	"fmt"
)

func main() {
	// fmt.Println("Hello World! THis is another test")

	// /*Variables and printing*/
	// var age int = 40
	// var facNum float64 = 1.6143243
	// fmt.Println(age, ": ", facNum)

	/*Go can choose variable type for you*/
	// randNUm := 1
	// var numOne = 1.000
	// var num99 = 0.999
	// fmt.Println(randNUm)
	// fmt.Println(numOne - num99)
	// /* + . * kan brukes i Println*/

	// /*You can define constants without using them*/
	// const pi float64 = 3.1415
	// // var (
	// // 	varA = 2
	// // 	varB = 3
	// // )

	// var myName string = "Harald Lønsethagen"
	// fmt.Println(len(myName)) //Length of string
	// fmt.Println(myName + " is a robot")
	// fmt.Println("Newlines\n")

	// var isOver40 bool = true
	// fmt.Println("%.3f \n", pi)     //3 Decimals
	// fmt.Println("%T \n", pi)       //Outputs: float64
	// fmt.Println("%t \n", isOver40) //Ouputs: true
	// fmt.Println("%d \n", 100)      //Outputs: 100
	// // && - AND
	// // || - OR
	// // ! - NOT
	// // == - EQUALS
	// // != - NOT EQUALS

	// /*For-loops*/
	// i := 1
	// for i <= 10 {
	// 	fmt.Println(i)
	// 	i++
	// }

	/*Prints from 0 to 4*/
	// for j := 0; j < 5; j++ {
	// fmt.Println(j)
	// }

	// yourAge := 17
	// // if yourAge >= 18 {
	// // 	fmt.Println("You can vote")
	// // } else if yourAge >= 16 {
	// // 	fmt.Println("You can drive")
	// // }

	// /*Switch*/
	// yourAge = 18 // not := cause allready assigned type above
	// switch yourAge {
	// case 16:
	// 	fmt.Println("Go drive")
	// case 18:
	// 	fmt.Println("Go vote")
	// default:
	// 	fmt.Println("Go have fun")
	// }

	/*Arrays*/
	// var favNum [5]float64
	// // All values in array get initiated to 0
	// favNum[0] = 163
	// favNum[1] = 1432
	// favNum[4] = 1.543

	// fmt.Println(favNum)

	// // Another array
	// favNum2 := [5]float64{1, 2, 3, 4, 5}

	// /*Range in every value out of our array*/
	// So i->0->4 and value->1->5
	// for i, value := range favNum2 {
	// 	fmt.Println(i, ": ", value)
	// }

	/*Slice is like an array but without the size*/
	numSlice := []int{5, 4, 3, 2, 1}
	numSlice2 := numSlice[1:4] //From index 1 to 3. Not 4

	fmt.Println("numSlice: =", numSlice)
	fmt.Println("numSlice[1:4] =", numSlice2)
	fmt.Println("numSlice[:2]:", numSlice[:2]) //Two first
	fmt.Println("numSlice[2:]", numSlice[2:])  //From third and above

	/*Default value of 0 for first 5. Max size of slice to 10. 10 total values*/
	numSlice3 := make([]int, 5, 10)
	fmt.Println("numSclie3: ", numSlice3)

	copy(numSlice3, numSlice)
	fmt.Println("numSlice3 after copy:", numSlice3)

}
