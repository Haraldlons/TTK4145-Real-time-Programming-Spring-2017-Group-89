package main

// Følger følende tutorial
// https://www.youtube.com/watch?v=CF9S4QZuV30&t=923s

// import (
// 	"fmt"
// 	"strings"
// )

/*
"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
*/

// func main() {
// 	// fmt.Println("Hello World! THis is another test")

// 	// /*Variables and printing*/
// 	// var age int = 40
// 	// var facNum float64 = 1.6143243
// 	// fmt.Println(age, ": ", facNum)

// 	/*Go can choose variable type for you*/
// 	// randNUm := 1
// 	// var numOne = 1.000
// 	// var num99 = 0.999
// 	// fmt.Println(randNUm)
// 	// fmt.Println(numOne - num99)
// 	// /* + . * kan brukes i Println*/

// 	// /*You can define constants without using them*/
// 	// const pi float64 = 3.1415
// 	// // var (
// 	// // 	varA = 2
// 	// // 	varB = 3
// 	// // )

// 	// var myName string = "Harald Lønsethagen"
// 	// fmt.Println(len(myName)) //Length of string
// 	// fmt.Println(myName + " is a robot")
// 	// fmt.Println("Newlines\n")

// 	// var isOver40 bool = true
// 	// fmt.Println("%.3f \n", pi)     //3 Decimals
// 	// fmt.Println("%T \n", pi)       //Outputs: float64
// 	// fmt.Println("%t \n", isOver40) //Ouputs: true
// 	// fmt.Println("%d \n", 100)      //Outputs: 100
// 	// // && - AND
// 	// // || - OR
// 	// // ! - NOT
// 	// // == - EQUALS
// 	// // != - NOT EQUALS

// 	/*For-loops*/
// 	// i := 1
// 	// for i <= 10 {
// 	// 	fmt.Println(i)
// 	// 	i++
// 	// }

// 	/*Prints from 0 to 4*/
// 	// for j := 0; j < 5; j++ {
// 	// fmt.Println(j)
// 	// }

// 	// yourAge := 17
// 	// // if yourAge >= 18 {
// 	// // 	fmt.Println("You can vote")
// 	// // } else if yourAge >= 16 {
// 	// // 	fmt.Println("You can drive")
// 	// // }

// 	// /*Switch*/
// 	// yourAge = 18 // not := cause allready assigned type above
// 	// switch yourAge {
// 	// case 16:
// 	// 	fmt.Println("Go drive")
// 	// case 18:
// 	// 	fmt.Println("Go vote")
// 	// default:
// 	// 	fmt.Println("Go have fun")
// 	// }

// 	/*Arrays*/
// 	// var favNum [5]float64
// 	// // All values in array get initiated to 0
// 	// favNum[0] = 163
// 	// favNum[1] = 1432
// 	// favNum[4] = 1.543

// 	// fmt.Println(favNum)

// 	// // Another array
// 	// favNum2 := [5]float64{1, 2, 3, 4, 5}

// 	// /*Range in every value out of our array*/
// 	// So i->0->4 and value->1->5
// 	// for i, value := range favNum2 {
// 	// 	fmt.Println(i, ": ", value)
// 	// }

// 	/*Slice is like an array but without the size*/
// 	// This is how to make an array of ints
// 	// numSlice := []int{5, 4, 3, 2, 1}
// 	// numSlice2 := numSlice[1:4] //From index 1 to 3. Not 4

// 	// fmt.Println("numSlice: =", numSlice)
// 	// fmt.Println("numSlice[1:4] =", numSlice2)
// 	// fmt.Println("numSlice[:2]:", numSlice[:2]) //Two first
// 	// fmt.Println("numSlice[2:]", numSlice[2:])  //From third and above

// 	// /*Default value of 0 for first 5. Max size of slice to 10. 10 total values*/
// 	// numSlice3 := make([]int, 5, 10)
// 	// fmt.Println("numSclie3: ", numSlice3)

// 	// copy(numSlice3, numSlice)
// 	// fmt.Println("numSlice3 after copy:", numSlice3)

// 	// // Appender 0 og -1 til listen
// 	// numSlice3 = append(numSlice3, 0, -1)
// 	// fmt.Println(numSlice3)

// 	// /*Maps. Key values*/
// 	// // Maps the age of presidents. Index is a string, while the value is age(int)
// 	// presAge := make(map[string]int)
// 	// presAge["ThieodoreRoosevelt"] = 42
// 	// fmt.Println(presAge["ThieodoreRoosevelt"]) //O: 42

// 	// // leng(presAge) //O: 1

// 	// presAge["John F. Kennedy"] = 43
// 	// delete(presAge, "John F. Kennedy") //Delete

// 	// /*Functions*/

// 	listNums := []float64{1, 2, 3, 4, 5}
// 	// fmt.Println(listNums)

// 	fmt.Println(addThemUp(listNums))

// }

/*
func function_name(parameter_name []array_type) return_type{

}


*/
// func addThemUp(numbers []float64) float64 {
// 	sum := 0.0
// 	for _, value := range numbers {
// 		sum += value
// 	}
// 	return sum
// 	// i needs to be used. if not used, i=_
// }

// func main() {
// 	num1, num2 := next2Values(5)
// 	fmt.Println("num1: ", num1, " and num2: ", num2)

// }

//  //A function can return two values
// func next2Values(number int) (int, int) {
// 	return number + 1, number + 2
// }

// func main() {
// 	fmt.Println(substractThem(1, 2, 3, 4, 5))
// }
// 	//Undefined how many input arguments
// func substractThem(arg ...int) int {
// 	finalValue := 0
// 	for _, value := range arg {
// 		finalValue -= value
// 	}
// 	return finalValue
// }

// func main() {
// 	num3 := 3

// 	// Define function inside main. Then function can 'reach' all variables in main
// 	doubleNum := func() int {
// 		num3 *= 2
// 		return num3
// 		//Can access num3 even tho its outside
// 	}
// 	fmt.Println(doubleNum())
// 	fmt.Println(doubleNum())

// }

// Recursion
// func main() {
// 	fmt.Println(factorial(3))
// }

// func factorial(start int) int {
// 	if start == 1 {
// 		return 1
// 	}
// 	return factorial(start-1) * start
// }

// Defer waits until going out of closing brackets
// In this cause this is main {}
// func main() {
// 	defer printTwo()
// 	printOne()
// }

// func printOne() { fmt.Println(1) }
// func printTwo() { fmt.Println(2) }

// func main() {
// 	fmt.Println(saveDiv(3, 0))
// 	fmt.Println(saveDiv(3, 2))
// }

// // If both parameters are of the same type
// func saveDiv(num1, num2 int) int {
// 	defer func() {
// 		fmt.Println(recover())
// 	}()
// 	// Recover will catch if error accour
// 	// Ex if num2 = 0 recover will take care of it
// 	solution := num1 / num2
// 	return solution
// }

// Notice: Program will still run even tho 3/0 causes error.

//Panic
// func main() {
// 	demPanic()
// }

// func demPanic() {
// 	defer func() {
// 		fmt.Println(recover())
// 	}()
// 	panic("PANIC")
// 	// I'm not sure, but I think panic fails the program with the error message "PANIC", which recover() recovers
// }

// Pointers (without here)
// func main() {
// 	x := 0
// 	changeXVal(x) // Send value of x
// 	fmt.Println("x =", x)
// }

// func changeXVal(x int) {
// 	x = 2
// }

// With pointers
// func main() {
// 	x := 0
// 	changeXVal(&x) // Send value of x
// 	fmt.Println("x =", x)
// }

// func changeXVal(x *int) {
// 	*x = 2
// }

// Structs
// func main() {
// 	// If I don't know the sequence of parameters in struct
// 	// rect1 := Rectangle{leftX: 0, topY: 50, width: 10, height: 20}
// 	// If I know the sequence of parameters
// 	rect1 := Rectangle{0, 50, 10, 20}
// 	fmt.Println("rect1: ", rect1)
// 	fmt.Println("Area: ", rect1.area())
// }

// type Rectangle struct {
// 	leftX  float64
// 	topY   float64
// 	height float64
// 	width  float64
// }

// func (rect *Rectangle) area() float64 {
// 	return rect.width * rect.height
// }

// func main() {
// 	rect := Rectangle{20, 50}
// 	circ := Circle{4}
// 	// fmt.Println(rect, circ)

// 	// Now comes the magic
// 	fmt.Println("Rectangle area: ", getArea(rect))
// 	fmt.Println("Circle area: ", getArea(circ))
// }

// type Shape interface {
// 	area() float64
// }

// type Rectangle struct {
// 	height float64
// 	width  float64
// }

// type Circle struct {
// 	radius float64
// }

// func (r Rectangle) area() float64 {
// 	return r.height * r.width
// }

// func (c Circle) area() float64 {
// 	return math.Pi * math.Pow(c.radius, 2) //Need to include "math"
// }

// func getArea(shape Shape) float64 {
// 	return shape.area()
// }

// String cuntion, IO, input user , webserver, go-routines, channels
// func main() {

// 	sampString := "Hello World"
// 	fmt.Println(strings.Contains(sampString, "lo")) /*O: true*/
// 	fmt.Println(strings.Index(sampString, "lo")) /*O: 3*/
// 	fmt.Println(strings.Count(sampString, "l")) /*O: 3*/
// 	fmt.Println(strings.Replace(sampString, "l", "x", 3)) /*Replace all l's with x's for the first 3*/
// 	// fmt.Println(strings.Contains(samplString, "lo")

// }

// import (
// 	"fmt"
// 	"sort"
// 	"strings"
// )

// func main() {
// 	csvString := "1,2,3,4,5,6"
// 	fmt.Println(strings.Split(csvString, ","))

// 	listOfLetters := []string{"c", "a", "b"}
// 	sort.Strings(listOfLetters)

// 	fmt.Println(listOfLetters)

// 	listOfNums := strings.Join([]string{"3", "2", "1"}, ",")
// 	fmt.Println(listOfNums)

// }

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os" /*Package for */
// )

// func main() {

// 	file, err := os.Create("samp.txt")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	file.WriteString("This is some random text\n")
// 	file.WriteString("This is even more text\n")

// 	file.Close()

// 	stream, err := ioutil.ReadFile("samp.txt")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	readString := string(stream)
// 	fmt.Println(readString)
// }

// import (
// 	"fmt"
// 	"strconv"
// )

// func main() {
// 	randInt := 5
// 	randFloat := 10.5
// 	randString := "100"
// 	randString2 := "250.5"

// 	/*Cast from int to float64*/
// 	fmt.Println(float64(randInt))
// 	fmt.Println(int(randFloat))

// 	newInt, _ := strconv.ParseInt(randString, 0, 64)
// 	fmt.Println(newInt)

// 	newFloat, _ := strconv.ParseFloat(randString2, 64)
// 	fmt.Println(newFloat)

// }

//Network-Server
// import (
// 	"fmt"
// 	"net/http"
// )

// func main() {

// 	http.HandleFunc("/", handler)

// 	http.HandleFunc("/earth", handler2)

// 	http.ListenAndServe(":8000", nil)

// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello World\n")
// }
// func handler2(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello Earth\n")
// }

//Routines

// import (
// 	"fmt"
// 	"time"
// )

// func count(id int) {
// 	for i := 0; i < 10; i++ {
// 		fmt.Println(id, ": ", i)
// 		time.Sleep(time.Millisecond * 1000)
// 	}
// }

// func main() {
// 	for i := 0; i < 10; i++ {
// 		go count(i)
// 		time.Sleep(time.Millisecond * 10)
// 	}

// 	time.Sleep(time.Millisecond * 11000)
// }

import (
	"fmt"
	"strconv"
	"time"
)

var pizzaNum = 0
var pizzaName = ""

//GO routines
func makeDough(stringChan chan string) {
	pizzaNum++
	pizzaName = "Pizza #" + strconv.Itoa(pizzaNum)
	fmt.Println("Make Dough and Send for Sauce")

	stringChan <- pizzaName

	time.Sleep(time.Millisecond * 10)
}

func addSauce(stringChan chan string) {
	pizza := <-stringChan

	fmt.Println("Add Toppings to", pizza, "and ship")

	stringChan <- pizzaName

	time.Sleep(time.Millisecond * 10)
}

func addToppings(stringChan chan string) {
	pizza := <-stringChan

	fmt.Println("Add Sauce and Send", pizza, " for toppings")

	stringChan <- pizzaName

	time.Sleep(time.Millisecond * 10)
}

func main() {
	stringChan := make(chan string)

	for i := 0; i < 3; i++ {
		go makeDough(stringChan)
		go addSauce(stringChan)
		go addToppings(stringChan)

		time.Sleep(time.Millisecond * 5000)
	}

}
