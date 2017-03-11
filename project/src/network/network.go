package network

import (
	"../definitions"
	// "../driver"
	// "../buttons"
	"../storage"
	// "../elevator"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
	// "math"
	//"bufio"
	// "log"
	//"os"
	//"strconv"
)

var bcAddress string = "129.241.187.255"

//var bcAddress string = "localhost"
var port string = ":55748"
var slaveSendPort string = ":55758"
var jsonSendPort string = ":55656"

var delay100ms = 100 * time.Millisecond

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func slave(udpListen *net.UDPConn) int {
	defer fmt.Println("Slave function ended") /*For debugging*/
	listenChan := make(chan int, 1)
	slaveCount := 0

	stopSlaveBroadcasting := make(chan int)

	// Run goroutine listening for sent values from master.
	go listen(listenChan, udpListen)
	go sendImAliveMessage(stopSlaveBroadcasting)

	for {
		select {
		case slaveCount = <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("Got listen message from master: ", slaveCount)
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			time.Sleep(delay100ms / 2) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			fmt.Println("Master is dead. Long live the the new king!")
			stopSlaveBroadcasting <- 1
			return slaveCount
		}
	}

}

func master(startCount int, udpBroadcast *net.UDPConn) {
	/* Launch new instance of "main".
	 * This creates the corresponding slave which will loop on listen until master dies
	 */
	newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err := newSlave.Run()
	check(err)
	msg := make([]byte, 8)

	count := startCount

	go listenAfterImAliveMessage()

	for {
		// Convert count from int to binary/byte and place in msg
		binary.BigEndian.PutUint64(msg, uint64(count))
		udpBroadcast.Write(msg)

		// fmt.Println(count)
		count++

		time.Sleep(10 * delay100ms) // Wait 1 cycle (100 ms)
	}
}

func listen(listenChan chan int, udpListen *net.UDPConn) {
	buf := make([]byte, 8)
	for {
		udpListen.ReadFromUDP(buf)

		// Convert byte from buf to int and send over channel.
		listenChan <- int(binary.BigEndian.Uint64(buf))
		time.Sleep(10 * delay100ms) // Wait 1 cycle (100 ms)
	}
}

func SetupNetwork() {

	sendJSON()

	udpAddr, err := net.ResolveUDPAddr("udp", port)
	check(err)

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)

	// Initialize slave
	// First run of program will return 0 and initialize master->slave topology
	fmt.Println("Run slave")
	count := slave(udpListen)

	udpListen.Close()

	udpAddr, err = net.ResolveUDPAddr("udp", bcAddress+port)
	check(err)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	count = 10
	fmt.Println("Run master")
	go recieveJSON()
	master(count, udpBroadcast)

	fmt.Println("Close broadcast")
	udpBroadcast.Close()
}

func sendImAliveMessage(stopSlaveBroadcasting chan int) {
	defer fmt.Println("Actually stopping sending sendImAliveMessage")

	fmt.Println("Sending order over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+slaveSendPort)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	binary.BigEndian.PutUint64(msg, uint64(66))
	// go func() {
	// 	select {
	// 	case <- stopSlaveBroadcasting:
	// 		fmt.Println("Ending sendImAliveMessage")
	// 		return
	// 	}
	// }()
	for {
		select {
		case <-stopSlaveBroadcasting:
			return
		default:
			udpBroadcast.Write(msg)
			// fmt.Println("Sending I'm Alive")
			time.Sleep(delay100ms * 20)
		}
	}

}

func listenAfterImAliveMessage() {
	fmt.Println("Listening after ImAlive messages")
	udpAddr, err := net.ResolveUDPAddr("udp", slaveSendPort)
	check(err)

	slaveMessagesRecieved := 0

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)

	listenChan := make(chan int, 1)
	slaveCount := 0

	go func() {
		buf := make([]byte, 8)
		for {
			udpListen.ReadFromUDP(buf)

			// Convert byte from buf to int and send over channel.
			listenChan <- int(binary.BigEndian.Uint64(buf))
			time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
		}
	}()

	for {
		select {
		case slaveCount = <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("ListenAfterSlaves: ", slaveCount, ", slaveMessagesRecieved: ", slaveMessagesRecieved)
			slaveMessagesRecieved++
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			// sendImAliveMessage()
			time.Sleep(delay100ms) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			fmt.Println("Have not recieved slave message for the last 3 seconds!")
			newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			err := newSlave.Run()
			check(err)
			// return
		}
	}

}

func setOrderOverNetwork(destinationFloor int) {
	fmt.Println("Sending order over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+port)
	check(err)
	msg := make([]byte, 8)

	// Create bcast Conn
	udpBroadcast, err := net.DialUDP("udp", nil, udpAddr)
	check(err)

	binary.BigEndian.PutUint64(msg, uint64(destinationFloor))
	udpBroadcast.Write(msg)
}

func sendJSON() {
	defer fmt.Println("Finished sending JSON")

	fmt.Println("Sending JSON over network")
	udpAddr, err := net.ResolveUDPAddr("udp", bcAddress+jsonSendPort)
	check(err)
	// msg := make([]byte, 128)
	// m := definitions.TestMessage{"Alice", "Hello", 1294706395881547000}
	// b, _ := json.Marshal(m)

	b := []byte(`[
  {
    "_id": "58bc3cbe7c880a1fb03cf518",
    "index": 0,
    "guid": "71963ee3-a7dc-45b0-bcf4-69823cf8ef1f",
    "isActive": false,
    "balance": "$2,463.19",
    "picture": "http://placehold.it/32x32",
    "age": 38,
    "eyeColor": "brown",
    "name": "Reyes Travis",
    "gender": "male",
    "company": "IPLAX",
    "email": "reyestravis@iplax.com",
    "phone": "+1 (837) 420-4000",
    "address": "505 Boynton Place, Accoville, Tennessee, 8967",
    "about": "Ad ad Lorem aute laborum eu aute consequat occaecat cupidatat veniam deserunt. Eu aliquip culpa nulla duis quis elit eu. Enim aute minim mollit id.\r\n",
    "registered": "2016-02-08T11:02:49 -01:00",
    "latitude": -84.569896,
    "longitude": 134.766522,
    "tags": [
      "laboris",
      "et",
      "veniam",
      "fugiat",
      "consequat",
      "minim",
      "minim"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Joanna Whitfield"
      },
      {
        "id": 1,
        "name": "Clayton Delaney"
      },
      {
        "id": 2,
        "name": "Francis Good"
      }
    ],
    "greeting": "Hello, Reyes Travis! You have 9 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "58bc3cbe8948d8853fa7d737",
    "index": 1,
    "guid": "0aa38531-3490-4986-b0cd-9649b215922c",
    "isActive": true,
    "balance": "$2,444.94",
    "picture": "http://placehold.it/32x32",
    "age": 26,
    "eyeColor": "green",
    "name": "Briana Levy",
    "gender": "female",
    "company": "NAMEBOX",
    "email": "brianalevy@namebox.com",
    "phone": "+1 (880) 480-2445",
    "address": "431 Pleasant Place, Churchill, Montana, 6745",
    "about": "Irure in in reprehenderit excepteur deserunt duis. Ullamco minim cupidatat eiusmod aliquip incididunt velit laborum tempor nisi do aliqua nulla deserunt sit. Cillum eu minim pariatur do. Sunt esse nostrud exercitation magna non ut qui aliqua sint in exercitation sit consectetur ea.\r\n",
    "registered": "2016-06-13T03:17:13 -02:00",
    "latitude": -50.55551,
    "longitude": -82.904125,
    "tags": [
      "laboris",
      "do",
      "cillum",
      "duis",
      "sint",
      "adipisicing",
      "dolor"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Virgie Frazier"
      },
      {
        "id": 1,
        "name": "Baird Cain"
      },
      {
        "id": 2,
        "name": "Bobbi Zamora"
      }
    ],
    "greeting": "Hello, Briana Levy! You have 2 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "58bc3cbee069ea0693bda1b4",
    "index": 2,
    "guid": "c448ace9-685f-43b8-80bb-03057b921b35",
    "isActive": false,
    "balance": "$3,708.77",
    "picture": "http://placehold.it/32x32",
    "age": 40,
    "eyeColor": "blue",
    "name": "Sophie Bender",
    "gender": "female",
    "company": "COGENTRY",
    "email": "sophiebender@cogentry.com",
    "phone": "+1 (804) 518-2791",
    "address": "353 Lancaster Avenue, Statenville, West Virginia, 9620",
    "about": "Nostrud ut Lorem magna ipsum excepteur culpa. Ex ad non duis eu sit aliquip do aliquip cupidatat dolore amet proident. Exercitation proident do officia incididunt dolore fugiat excepteur tempor esse fugiat ad consequat nulla. Enim aliqua consectetur eu in qui velit commodo amet ut laboris adipisicing dolore labore. Fugiat sit dolor cillum magna. Eu sint laboris eu qui elit.\r\n",
    "registered": "2014-12-10T11:32:29 -01:00",
    "latitude": -74.597476,
    "longitude": 132.724518,
    "tags": [
      "deserunt",
      "officia",
      "ut",
      "tempor",
      "elit",
      "velit",
      "anim"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Kimberley Mercer"
      },
      {
        "id": 1,
        "name": "Head Douglas"
      },
      {
        "id": 2,
        "name": "Terra Wilcox"
      }
    ],
    "greeting": "Hello, Sophie Bender! You have 3 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "58bc3cbe455d201287b2a86e",
    "index": 3,
    "guid": "216692f5-a149-4e9c-986a-45034b242a25",
    "isActive": false,
    "balance": "$1,344.18",
    "picture": "http://placehold.it/32x32",
    "age": 33,
    "eyeColor": "brown",
    "name": "Hampton Ochoa",
    "gender": "male",
    "company": "INTERLOO",
    "email": "hamptonochoa@interloo.com",
    "phone": "+1 (927) 415-3170",
    "address": "172 Dekoven Court, Sanborn, Northern Mariana Islands, 6722",
    "about": "Ex anim dolor eiusmod incididunt ipsum mollit reprehenderit sint id do tempor. Ea elit anim nulla veniam adipisicing ex elit mollit amet magna ut. Et labore est sit duis anim Lorem.\r\n",
    "registered": "2016-09-24T06:01:27 -02:00",
    "latitude": -17.876075,
    "longitude": -111.030459,
    "tags": [
      "dolor",
      "officia",
      "veniam",
      "do",
      "nostrud",
      "occaecat",
      "et"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Jasmine England"
      },
      {
        "id": 1,
        "name": "Gwendolyn Albert"
      },
      {
        "id": 2,
        "name": "Potter Morrow"
      }
    ],
    "greeting": "Hello, Hampton Ochoa! You have 8 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "58bc3cbecf74c06b5fc43c4f",
    "index": 4,
    "guid": "be2624dc-589a-4db6-b563-3f6b55255d29",
    "isActive": true,
    "balance": "$3,451.72",
    "picture": "http://placehold.it/32x32",
    "age": 40,
    "eyeColor": "green",
    "name": "Barr Lara",
    "gender": "male",
    "company": "CIPROMOX",
    "email": "barrlara@cipromox.com",
    "phone": "+1 (815) 598-3047",
    "address": "384 Cox Place, Highland, Nebraska, 6953",
    "about": "Minim cillum qui fugiat in Lorem aute. Consectetur non et aliquip nostrud consequat labore incididunt tempor. Aliqua incididunt nostrud duis amet proident consequat ex sunt sit laboris nisi. Veniam irure aliqua deserunt anim sint laboris pariatur deserunt ea est. Ad Lorem id officia magna quis mollit id aliquip proident nulla quis ex. Officia eu et aute ipsum proident qui nulla non in.\r\n",
    "registered": "2014-08-09T05:19:00 -02:00",
    "latitude": -78.464367,
    "longitude": 46.141521,
    "tags": [
      "adipisicing",
      "duis",
      "proident",
      "eiusmod",
      "sit",
      "in",
      "velit"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Norma Gay"
      },
      {
        "id": 1,
        "name": "Coffey Stark"
      },
      {
        "id": 2,
        "name": "Lea Mcneil"
      }
    ],
    "greeting": "Hello, Barr Lara! You have 10 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "58bc3cbee5c175af5b6676d7",
    "index": 5,
    "guid": "8b1463b3-b291-46d1-9bcd-e9a379e0e26e",
    "isActive": false,
    "balance": "$3,037.76",
    "picture": "http://placehold.it/32x32",
    "age": 20,
    "eyeColor": "blue",
    "name": "Gonzalez James",
    "gender": "male",
    "company": "EVEREST",
    "email": "gonzalezjames@everest.com",
    "phone": "+1 (826) 493-3049",
    "address": "545 Falmouth Street, Singer, Minnesota, 4644",
    "about": "Proident aliquip quis anim esse nulla. Labore Lorem ad veniam amet. Veniam magna deserunt enim adipisicing officia sint exercitation ea.\r\n",
    "registered": "2015-09-14T08:06:41 -02:00",
    "latitude": 63.001256,
    "longitude": 28.082829,
    "tags": [
      "ipsum",
      "qui",
      "et",
      "ipsum",
      "anim",
      "consequat",
      "eiusmod"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Danielle Norris"
      },
      {
        "id": 1,
        "name": "Houston Justice"
      },
      {
        "id": 2,
        "name": "English Wood"
      }
    ],
    "greeting": "Hello, Gonzalez James! You have 5 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "58bc3cbef097e05f09f0c7ae",
    "index": 6,
    "guid": "3b449f58-0b56-49b6-ba28-088d7fd9ff3b",
    "isActive": true,
    "balance": "$3,985.43",
    "picture": "http://placehold.it/32x32",
    "age": 37,
    "eyeColor": "green",
    "name": "Alfreda Sims",
    "gender": "female",
    "company": "UBERLUX",
    "email": "alfredasims@uberlux.com",
    "phone": "+1 (849) 498-2541",
    "address": "859 Lawrence Avenue, Bourg, Oklahoma, 186",
    "about": "Voluptate aliquip sunt sit fugiat labore do fugiat dolore in cillum nostrud. Occaecat tempor cillum occaecat enim mollit. Esse qui incididunt exercitation non nostrud adipisicing. Eiusmod exercitation nostrud pariatur nisi.\r\n",
    "registered": "2016-06-11T11:44:28 -02:00",
    "latitude": 25.431983,
    "longitude": -25.659437,
    "tags": [
      "elit",
      "non",
      "pariatur",
      "magna",
      "ullamco",
      "sunt",
      "occaecat"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Brady Sandoval"
      },
      {
        "id": 1,
        "name": "Maynard Holland"
      },
      {
        "id": 2,
        "name": "Jennings Drake"
      }
    ],
    "greeting": "Hello, Alfreda Sims! You have 8 unread messages.",
    "favoriteFruit": "banana"
  }
]`)

	fmt.Println("JSON in ByteArray:", b)
	jsonByteLength := len(b)
	firstByte := jsonByteLength / 255
	// fmt.Println("firstByte", firstByte)
	secondByte := jsonByteLength - firstByte*255
	// fmt.Println("secondByte:",secondByte)
	// fmt.Println("JSONByteArrayLength:",jsonByteLength)

	fmt.Println(byte(len(b)))

	b = append([]byte{byte(secondByte)}, b...)
	b = append([]byte{byte(firstByte)}, b...)
	// fmt.Println("WITH Length as first byte", b)

	// Create bcast Conn
	udpBroadcast, _ := net.DialUDP("udp", nil, udpAddr)
	// check(err)

	udpBroadcast.Write(b)
}

func recieveJSON() {
	fmt.Println("Listening after JSON Objectes")
	udpAddr, err := net.ResolveUDPAddr("udp", jsonSendPort)
	check(err)

	JSONobjectsRecieved := 0

	// Create listen Conn
	udpListen, err := net.ListenUDP("udp", udpAddr)
	check(err)
	var m interface{}

	listenChan := make(chan interface{}, 1)
	// slaveCount := 0

	go func() {
		buf := make([]byte, 65536) /*2^16 = max recovery size*/
		for {
			udpListen.ReadFromUDP(buf)

			// Two first bytes contains the size of the JSON byte array

			// fmt.Println("buffer after read from UDP: ", buf)
			jsonByteLength := int(buf[0])*255 + int(buf[1])
			// fmt.Println("jsonByteLength:",jsonByteLength)

			// Convert byte from buf to int and send over channel.
			err := json.Unmarshal(buf[2:jsonByteLength+2], &m)
			storage.SaveElevatorStateToFile(m) //This actually works
			// fmt.Println("Her kommer m som du skal se på: ", m)
			// fmt.Println("Ferdig med å vise m")
			check(err)
			listenChan <- m
			time.Sleep(delay100ms) // Wait 1 cycle (100 ms)
		}
	}()


	for {
		select {
		case JSONByteArray := <-listenChan:
			// fmt.Println("slaveCount: ", slaveCount)
			fmt.Println("got JSON object: ", JSONByteArray, ", json objects recieved: ", JSONobjectsRecieved)
			JSONobjectsRecieved++
			// if slaveCount < 4 && slaveCount > -1 {
			// 	fmt.Println("Going to floor from slave: ", slaveCount)
			// 	// go goToFloor(slaveCount, &elevatorState)
			// }
			// sendImAliveMessage()
			time.Sleep(delay100ms / 2) // wait 50 ms
			break
		case <-time.After(30 * delay100ms): // Wait 10 cycles (1 second). Master assumed dead
			// When master dies, slavecount is returned so that a new process of master -> slave
			// can continue from the last value sent over the network.
			// fmt.Println("Have not recieved any JSON message for the last 3 seconds!")
			// newSlave := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			// err := newSlave.Run()
			check(err)
			// return
		}
	}

}

// func decodeJSONRecievedAndStore(){

// }


func SendFromSlaveToMaster(externalButtonPressesChan chan, ordersChan chan, elevatorState definitions.ElevatorState, id string) { 
  MSG := definitions.MSG_to_master {
    Orders: Orders, 
    ElevatorState: state,
    ExternalButtonPresses: <-externalButtonPressesChan, 
    Id: id,//Id string 
  }

  sendJSON(MSG)
}

func SendFromMasterToSlave(elevators){
  MSG := definitions.MSG_to_slave {
    Elevators: elevators,
  }
  sendJSON(MSG)
}

func ListenToMaster() {

}

func ListenToSlave() {
  
}
>>>>>>> harald
