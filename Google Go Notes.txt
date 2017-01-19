Go Tutorials 

https://www.youtube.com/watch?v=CF9S4QZuV30&t=923s
"Go Programming" - Youtube Tutorial
golang.org for install
"

package main
import "fmt" //format package
//Comments
/*  comment */
func main(){
	fmt.Println("Hello World")
	var age int = 40;
	var favNum float64 = 1.6180339
	fmt.Println(age, favNum)
	
	//Go can choose variable type for you
	randNum := 1; //Automatic create int. Must then remain integer
	var numOne = 1.000
	var num99 = 0.9999
	fmt.Println(numOne-num99)
	//
	// + . * / kan brukes i Println
	
	const pi float64 = 3.1415
	var (
		varA = 2
		varB = 3
	)
	
	var myName string = "Harald Lønsethagen"
	fmt.Println(len(myName)) //Length of string
	fmt.Println(myName + " is a robot")
	fmt.Println("Newlines\n")
	
	var isOver40 bool = true
	fmt.Printf("%.3f \n", pi) //3 decimals
	fmt.Printf("%T \n", pi) //outputs float64
	fmt.Printf("%t \n", isOver40) //outputs: true
	fmt.Printf("%d \n", 100) //outputs: 100
	// && and
	// || org
	// ! not
	// == 
	// != 
	
	
	//for-loops
	i := 1
	for i <= 10 {
		fmt.Println(i);
		i++ 
	}
	
	for j := 0; j < 5; j++{
		fmt.Println(j)
	}
	
	yourAge := 18
	if yourAge >= 16 {
		fmt.Println("You can drive")
	} else if yourAge >= 18 {
		fmt.Println("You can vote")
	} 
	
	//switch
	yourAge := 5
	switch yourAge {
		case 16: fmt.Println("Go drive")
		case 18: fmt.Println("Go vite");
		default: fmt.Println("Go have fun")
	}
	
	//Arrays
	var favNum2[5] float64
	favNum2[0] = 163
	favNum2[1] = 1432
	...
	favNum2[4] = 1.543
	
	fmt.Println(favNum2[3])
	
	//Another array
	favNum3 := [5]float64 {1,2,3,4,5}
	
	//range every value out of our array
	for i, value := range favNum3{
		fmt.Println(value, i)
	}
	//Printer ut indexen og verdien
	
	
	//Slice is like an array but without the size.
	numSlice := []int {5,4,3,2,1}
	numSlice2 := numSlice[3:5] //Tar 3. og 4. og IKKE 5. men 5. finnes ikke i numSlice
	fmt.Println("numSlice[0] =", numSlice2[0])
	fmt.Println("numSlice[:2] =", numSlice[:2]) //start from 0 index to 1, NOT 2
	fmt.Println("numSlice[2:] =", numSlice[2:]) //start from 2 index to end
	
	//Default value of 0 for first 5. Max size of slice to 10. 10 total values
	numSlice3 := make([]int, 5, 10)
	
	//kopierer numSlice til numSlice3
	copy(numSlice3, numSlice)
	
	//
	numSlice3 = append(numSlice3,0,-1);
	fmt.Println(numSlice3[6]) //-1 shows up
	
	//Maps. Key values
	//Maps the age of presidents. index is string, values is age.
	presAge := make(map[string] int)
	presAge["ThieodoreRoosevelt"] = 42
	fmt.Println(presAge["ThieodoreRoosevelt"])
	leng(presAge) //1
	presAge["John F. Kennedy"] = 43
	delete(presAge, "John F. Kennedy"); //Delete
	
	//functions
	
	listNums := []float64{1,2,3,4,5}

	
	
}

To run
go run herewego.go

variables cant change 


//functions
package main
import "fmt"
func main(){
	listNums := []float64{1,2,3,4,5}
	fmt.println("Sum: ", addThemUp(listNums))
}
//addThemUp is function name. numbers is paramater name of type array of float64. return value is float64 of function.
func addThemUp(numbers []float64) float64{
	//Local variables
	sum := 0.0
	//If dont care about index value in for-loop "_"
	for _, val :=range numbers {
		sum += val
	}
	return sum
}

func main(){
	num1, num2 := next2Values(5)
	fmt.Println(num1,num2)
}

//Return two values from a function
func next2Values(number int)(int, int){
	return number+1, number+2
}

//Dont know how many input arguments we need to function
func main(){
	fmt.Println(substractThem(1,2,3,4,5))
}

func substractThem(arg ...int) int {
	finalValue := 0
	for _, value:=range args {
		finalValue -= value
	}
	return finalValue
}

//New
func main{
	num3 := 3
	doubleNum := func() int{
		num3 *= 2
		return num3
		//Can access num3 even tho its outside
	}
	fmt.Println(doubleNum())
	fmt.Println(doubleNum())
}

//Recursion
func main(){
	fmt.Println(factorial(3))
}
//factorial(3) = 3 * factorial(2) osv..
func factorial(num int) int {
	if num == 0 {
		return 1
		}
	return num * factorial(num-1)
}

// Defer
func main(){
	defer printTwo()
	printOne()
}
//printer 1 før 2. defer gjør at venter til etter enclosing avslutter(main i dette tilfellet)
func printOne(){fmt.Println(1)}
func printTwo(){fmt.Println(2)}

//
func main(){
	fmt.Println(saveDiv(3,0))
	fmt.Println(saveDiv(3,2))
	//Program continues even tho 3/0 undefined cause of recover
}
//Both paramters are integers
func safeDiv(num1, num2 int) int {
	defer func(){
		fmt.Println(recover()) //Continue execution even tho error accour
	}()
		//Recover will catch if error accour.
		//Ex if num2 = 0 recover will take care of it
	solution := num1 / num2
	return solution
}

//panic
func main(){
	demPanic()
}
func demPanic(){
	defer func8){
		fmt.Println(recover())
	}()
	panic("PANIC")
}

//Pointers (First ex is without)
func main(){
	x := 0
	changeXVal(x)//Send value of x
	fmt.Println("x =", x)
}
func chagneXVal(x int) {
	x = 2//Only effect inside function
}
//With pointers
func main(){
	x:=0
	chagneXValNow(&x)//Pass the referance to x
	fmt.Println("x: ",x)
	fmt.Println("Memory address for x =", &x)
}
func changeXValNow(x *int){ //* referance. Can then change the value at the memory address referanced by the pointer
	*x=2 //Store 2 in the place where x is stored in memory
}
//Slictly change
func main(){
	yPtr := new(int)
	changeYValNow(yPtr)
	fmt.Println("y:",*yPtr)
	
}
func changeYValNow(yPtr *int){
	*yPtr = 100
}

//Structs
func main(){
	//If I dont know the sequence of paramaters
	rect1 := Rectangle{leftX: 0, topY: 50, width: 10, height: 20}
	//If I KNOW the sequence
	rect1 := Rectangle{0,50,10,10}
	fmt.Println("Rectangle is", rect1.width, "wide")
	//Use function with struct
	fmt.Println("Area: ",rect1.area())
}
type Rectangle struct{
	leftX float64
	topY float64
	height float64
	width float64
} //Function name is area() . Methods to a struct
//Think this function "attackes" to the Rectangle struct
//area() not recieves any attributes, so no paramaters
func (rect *Rectangle) area() float64{
	return rect.width * rect.height
}

// New example
func main(){
	rect := Rectangle{20,50}
	circ := Circle{4}
	
		//This is the magic
	fmt.Prinln("Rectangle area: ", getArea(rect))
	fmt.Println("Circle Area: ", getArea(circ))
}

type Shape interface {
	are() float64
}
type Rectangle struct {
	height float64
	width float64
}
type Circle struct {
	radius float64
}
func (r Rectangle) area() float64{
	return r.height * r.width
}
func (c Circle) area() float64{
	return math.Pi * math.Pow(c.radius,2)
}
//All the area functions are tied together because of the interface
func getArea(shape Shape) float64{
	return shape.area()
}

//String functions
package main
import "fmt"
import "math"
func main(){
	sampString := 
}


