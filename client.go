// package main

// import (
//     "bufio"
//     "fmt"
//     "net"
//     "os"
//     "strconv"
// 	"github.com/shirou/gopsutil/mem"
// )

// func periodicUpdate() {
// 	for {
// 		// Generate a new snapshot
// 		newSnapshot, err := getSnapshot()
// 		if err != nil {
// 			log.Printf("Error retrieving snapshot: %v", err)
// 			time.Sleep(1 * time.Second) // Wait before retrying
// 			continue
// 		}

// 		// Update the cache (write lock)
// 		mu.Lock()
// 		snapshot = newSnapshot
// 		mu.Unlock()

// 		// Log update for debugging
// 		log.Println("Snapshot updated in cache.")

// 		// Wait for the next update
// 		time.Sleep(5 * time.Second) // Update interval
// 	}
// }

// func getSnapshot() (DataSnapshot, error) {
// 	v, err := mem.VirtualMemory()
// 	if err != nil {
// 		return DataSnapshot{}, err
// 	}
// 	metrics := []Metric{}
// 	metrics = append(metrics, Metric{Name: "Total", Value: float64(v.Total / 1024 / 1024)})
// 	metrics = append(metrics, Metric{Name: "Used", Value: float64(v.Used / 1024 / 1024)})
// 	metrics = append(metrics, Metric{Name: "UsedPercent", Value: float64(v.UsedPercent)})

// 	// return DataSnapshot{Id: "1234", Timestamp: time.Now(), Metrics: metrics}, nil
// 	temp := DataSnapshot{Id: "1234", Timestamp: time.Now(), Metrics: metrics}
// 	time.Sleep(5 * time.Second)
// 	return temp, nil

// }

// // manyMetrics := []DataSnapshot{}
// 	// for i := 0 ; i < 1 ; i++ {

// 	// 	v, err := getSnapshot()
// 	// 	if err != nil {
// 	// 		http.Error(w, "Error retrieving memory stats", http.StatusInternalServerError)
// 	// 		return
// 	// 	}

// 	// 	// Convert the struct to JSON
// 	// 	manyMetrics = append(manyMetrics, v)
// 	// }

// const (
//     serverAddress = "localhost:" // Server address
//     numMessages   = 50               // Number of messages to send
//     messageLength = 60               // Length of each message
// )

// // GenerateMessage creates a message of exactly 60 characters
// func generateMessage(i int) string {
//     baseMessage := "This is around sixty chars " + strconv.Itoa(i)
//     padding := "iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii"

//     // Ensure the message is exactly 60 characters long
//     msg := baseMessage + padding
//     if len(msg) > messageLength {
//         msg = msg[:messageLength] // Trim to 60 characters if it exceeds
//     }
//     return msg
// }

// func initClient(port int) {
//     // Connect to the server
//     conn, err := net.Dial("tcp", serverAddress + strconv.Itoa(port))
//     if err != nil {
//         fmt.Println("Error connecting to server:", err)
//         return
//     }
//     defer conn.Close()
//     fmt.Println("Connected to server at", serverAddress + strconv.Itoa(port))

//     // Send multiple messages to the server
//     for i := 0; i < numMessages; i++ {
//         msg := generateMessage(i)

//         // Calculate the message length and format it with a 4-digit prefix
//         msgLength := len(msg)
//         prefixedMessage := fmt.Sprintf("%04d%s", msgLength, msg)

//         // Send the complete message (length prefix + message) to the server
//         _, err = conn.Write([]byte(prefixedMessage))
//         if err != nil {
//             fmt.Println("Error sending message:", err)
//             return
//         }
//         fmt.Printf("Sent message %d: %s\n", i+1, prefixedMessage)
//     }

//     // Wait for keypress to exit
//     fmt.Println("\nPress Enter to exit.")
//     bufio.NewReader(os.Stdin).ReadBytes('\n')
// }

package main

import (
	"fmt"
	"net"
)

func main() {
	for i:=0 ; i < 1 ; i ++{
		go createClient(i);
		
	}
	for {}
}

func createClient(increment int) {
	serverAddr := "172.20.10.11:8888" // Replace with the actual server address and port
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")


	directions := []string{"U", "R", "D", "L"}
	dirIndex := 0
	name := fmt.Sprintf("art%d", increment)

	// Send incrementing name
	_, err = conn.Write([]byte(name + ";"))
	if err != nil {
		fmt.Printf("Error sending name: %v\n", err)
		return
	}
	fmt.Printf("Sent name: %s\n", name)

	for {
		direction := directions[dirIndex]
		_, err = conn.Write([]byte(direction + ";"))
		if err != nil {
			fmt.Printf("Error sending direction: %v\n", err)
			return
		}
		fmt.Printf("Sent direction: %s\n", direction)

		dirIndex = (dirIndex + 1) % len(directions)

	}

}
