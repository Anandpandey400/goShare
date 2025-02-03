# goShare

goShare is a Go-based command-line utility that allows you to:  - Upload files to [https://0x0.st](https://0x0.st)

## Setup and Installation

### Prerequisites

Before you can run the `goShare`, ensure that you have Go installed. If not, you can install Go from [here](https://golang.org/doc/install).

### Step 1: Clone the repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/goShare.git
cd goShare
```

### Step 2: Install Dependencies
This project uses the github.com/schollz/progressbar and github.com/fatih/color libraries. These should be automatically installed when you build the project, but you can manually install them using:

```bash
go get github.com/schollz/progressbar/v3
go get github.com/fatih/color
```
### Step 3: Build the Project
Run the following command to build the executable:

```bash
go build -o goShare main.go
```
This will create the goShare executable in the current directory.

### Usage
Once the program is built, you can use it from the command line.

### Upload a file
To upload a file, run the following command:

```bash
goShare /path/to/your/file.txt
```

### List active links
To view all active uploaded links, run

```bash
goShare --list
```
### Example
Upload a file:

```bash
goShare C:\path\to\file.txt
```

### Output:
 ```bash
File uploaded successfully! Access it at: https://0x0.st/abc123
Note: The file will expire in 12 hours.
```

#List active links:

```bash
goShare --list
```

### Output:

```bash
      ACTIVE UPLOADED LINKS
https://0x0.st/abc123

https://0x0.st/xyz456
```

## Contributing
If you'd like to contribute to this project, feel free to fork it, make changes, and create a pull request. Any contributions are welcome!

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements
https://0x0.st for the file hosting service.
Go Programming Language for providing an easy way to write cross-platform applications.
schollz/progressbar for the progress bar library.
