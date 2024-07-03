# Hardware Monitoring Tool in Go

This project is a simple hardware monitoring tool for an Ubuntu server, written in Go. It monitors CPU usage, memory usage, and disk usage, and outputs this information at regular intervals. The program can be terminated by pressing the `q` key or by sending a SIGINT or SIGTERM signal (e.g., `Ctrl+C`).

## Installation

1. **Install Go:** If not already installed, you can install Go with:
   ```sh
   sudo apt-get install golang
   ```

2. **Clone the project:**
   ```sh
   git clone https://github.com/your-username/hardware-monitoring-tool.git
   cd hardware-monitoring-tool
   ```

3. **Install dependencies:**
   ```sh
   go get -u github.com/shirou/gopsutil/cpu
   go get -u github.com/shirou/gopsutil/mem
   go get -u github.com/shirou/gopsutil/disk
   ```

## Usage

1. **Build the program:**
   ```sh
   go build -o monitor
   ```

2. **Run the program:**
   ```sh
   ./monitor
   ```

   The program will start monitoring hardware parameters and output the results every 5 seconds.

3. **Stop the program:**
   - Press the `q` key
   - Or send a SIGINT or SIGTERM signal (e.g., `Ctrl+C`)

## Code Overview

The source code is in the `main.go` file. Here is a brief overview:

- The main logic runs in an infinite loop that queries CPU, memory, and disk usage every 5 seconds and outputs the results.
- A separate goroutine monitors user input and stops the program when the `q` key is pressed.
- Another goroutine ensures the program stops cleanly when a SIGINT or SIGTERM signal is received.



## Contributions

Contributions are welcome! Please fork this repository and create a pull request with your changes.

## Author

Floorey
