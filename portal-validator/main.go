package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/godbus/dbus/v5"
	// "gocv.io/x/gocv"
)

const (
	portalBusName        = "org.freedesktop.portal.Desktop"
	portalObjectPath     = "/org/freedesktop/portal/desktop"
	fileChooserInterface = "org.freedesktop.portal.FileChooser"
	openFileMethod       = "OpenFile"
)

func main() {
	// Connect to the session bus
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	// Get the FileChooser object
	obj := conn.Object(portalBusName, dbus.ObjectPath(portalObjectPath))

	// Prepare method call arguments
	parentWindow := ""
	title := "Test File Chooser"
	options := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(true),
	}

	// Call the OpenFile method
	var response dbus.ObjectPath
	err = obj.CallWithContext(context.Background(),
		fileChooserInterface+"."+openFileMethod,
		0,
		parentWindow,
		title,
		options).Store(&response)

	if err != nil {
		log.Fatalf("Failed to call OpenFile method: %v", err)
	}

	fmt.Printf("File chooser dialog opened successfully. Response object path: %s\n", response)

	// Wait for user input to keep the program running
	fmt.Println("Press Enter to exit...")
	_, _ = os.Stdin.Read(make([]byte, 1))

	// accessCamera()
}

// func accessCamera() {

// 	// Open the default camera
// 	webcam, err := gocv.VideoCaptureFile("")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer webcam.Close()

// 	// Check if the camera is opened
// 	if ok := webcam.IsOpened(); !ok {
// 		fmt.Println("Failed to open the camera")
// 		return
// 	}

// 	// Get the camera's width and height
// 	// width := int(webcam.Get(gocv.VideoCaptureFrameWidth))
// 	// height := int(webcam.Get(gocv.VideoCaptureFrameHeight))

// 	// Create a window to display the video feed
// 	window := gocv.NewWindow("Camera Feed")

// 	// Process frames from the camera
// 	for {
// 		// Capture a new frame
// 		frame := gocv.NewMat()
// 		defer frame.Close()
// 		if ok := webcam.Read(&frame); !ok {
// 			fmt.Println("Failed to capture frame")
// 			break
// 		}

// 		// Display the frame in the window
// 		window.IMShow(frame)

// 		// Break the loop if 'q' is pressed
// 		if window.WaitKey(1) == 'q' {
// 			break
// 		}
// 	}

// 	// Release resources
// 	// frame.Close()
// 	window.Close()
// }
