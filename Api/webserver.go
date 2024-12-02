

package main

import (
    "net"
    "fmt"
    "strconv"
)

func initWebServer(port int) {
    // Listen for incoming connections
    listener, err := net.Listen("tcp", "0.0.0.0:" + strconv.Itoa(port))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on port "+ strconv.Itoa(port))

    for {
        // Accept incoming connections
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Handle client connection in a goroutine
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("Welcome to the server! Type something and press enter.\n"))
    fmt.Println("Connection from " + conn.RemoteAddr().String())

    // Buffer to hold the 4-byte length prefix
    lengthBuffer := make([]byte, 4)

    for {
        // Read the first 4 bytes to determine the length of the incoming message
        _, err := conn.Read(lengthBuffer)
        if err != nil {
            fmt.Println("Error reading length from client:", err)
            break
        }

        // Convert lengthBuffer to an integer
        length, err := strconv.Atoi(string(lengthBuffer))
        if err != nil {
            fmt.Println("Invalid length prefix:", string(lengthBuffer))
            break
        }

        // Read the message based on the specified length
        messageBuffer := make([]byte, length)
        _, err = conn.Read(messageBuffer)
        if err != nil {
            break
        }

        // Process the received message
        clientInput := string(messageBuffer)
        fmt.Println("Received from client:", clientInput)

        // Send a response back to the client without a length prefix
        response := "Message received: " + clientInput + "\n"
        conn.Write([]byte(response))
    }
}