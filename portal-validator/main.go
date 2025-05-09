package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus/v5"
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
	fmt.Println("Press Enter to move on to validate that cameras can be found...")
	_, _ = os.Stdin.Read(make([]byte, 1))

	findCameras()
}

func findCameras() {
	// Path where video devices are typically located
	devicePath := "/dev/video*"

	// Use filepath.Glob to find all matching video devices
	devices, err := filepath.Glob(devicePath)
	if err != nil {
		fmt.Println("Error searching for video devices:", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println("No cameras found.")
		return
	}

	fmt.Println("Cameras found:")
	for _, device := range devices {
		// Get device information
		info, err := os.Stat(device)
		if err != nil {
			fmt.Printf("Error getting info for %s: %v\n", device, err)
			continue
		}

		// Check if it's a character device (cameras are usually character devices)
		if info.Mode()&os.ModeCharDevice != 0 {
			fmt.Println(device)

			// Optionally, try to get more information about the device
			deviceInfo, err := getDeviceInfo(device)
			if err == nil {
				fmt.Printf("  Info: %s\n", deviceInfo)
			}
		}
	}
}

func getDeviceInfo(devicePath string) (string, error) {
	// Read the first few bytes of the device file
	// This might contain some identifying information
	file, err := os.Open(devicePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 256)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Convert to string and remove non-printable characters
	info := strings.Map(func(r rune) rune {
		if r >= 32 && r < 127 {
			return r
		}
		return -1
	}, string(buffer[:n]))

	return strings.TrimSpace(info), nil
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
