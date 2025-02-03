package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// File to store the active links
const linksFile = "links.txt"

// uploadFile uploads a file to the specified server and returns the URL
func uploadFile(filePath string) (string, error) {
	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// Create a buffer and multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form field for the file
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create form file: %w", err)
	}

	// Get the file size for the progress bar
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("unable to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Initialize the progress bar
	bar := progressbar.NewOptions(int(fileSize),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Uploading"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
	)

	// Copy the file data to the form field with progress bar updates
	_, err = io.Copy(io.MultiWriter(part, bar), file)
	if err != nil {
		return "", fmt.Errorf("error copying file data: %w", err)
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("unable to close writer: %w", err)
	}

	// Prepare the HTTP request
	request, err := http.NewRequest("POST", "https://0x0.st", bytes.NewReader(body.Bytes()))
	if err != nil {
		return "", fmt.Errorf("unable to create request: %w", err)
	}

	// Set the appropriate Content-Type header
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	// Read the URL from the response body
	url, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Save the URL to the links file
	err = saveLink(url)
	if err != nil {
		return "", fmt.Errorf("error saving URL: %w", err)
	}

	return string(url), nil
}

// saveLink saves the uploaded URL to a file
func saveLink(url []byte) error {
	// Open the file in append mode to add the new URL
	file, err := os.OpenFile(linksFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open links file: %w", err)
	}
	defer file.Close()

	// Write the URL to the file
	_, err = file.WriteString(string(url) + "\n")
	if err != nil {
		return fmt.Errorf("unable to write URL to file: %w", err)
	}

	return nil
}

// listLinks displays all the active uploaded links
func listLinks() error {
	// Open the links file for reading
	file, err := os.Open(linksFile)
	if err != nil {
		return fmt.Errorf("unable to open links file: %w", err)
	}
	defer file.Close()

	// Read and display the URLs
	fmt.Println("Active links:")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading links file: %w", err)
	}

	return nil
}

func printHelp() {
	fmt.Println("File Upload CLI - Uploads a file to https://0x0.st and gives you a URL")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go <file-path>")
	fmt.Println("  go run main.go --list      List all active links")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  --list           Display all active links")
	fmt.Println()
}

func main() {
	// Check for help or list flags
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	// Handle listing active links
	if os.Args[1] == "--list" {
		err := listLinks()
		if err != nil {
			fmt.Printf("\nError listing links: %v\n", err)
		}
		return
	}

	// Get the file path from command-line argument
	filePath := os.Args[1]

	// Upload the file
	url, err := uploadFile(filePath)
	if err != nil {
		fmt.Printf("\nError uploading file: %v\n", err)
		return
	}

	// Display the URL where the file is uploaded
	fmt.Printf("\nFile uploaded successfully! Access it at: \033[1;34m%s\033[0m\n", url)

	// Notify user that the file will expire in 12 hours
	fmt.Println("\nNote: The file will expire in 12 hours.")
}
