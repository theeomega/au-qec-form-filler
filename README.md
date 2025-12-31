# AU QEC Form Filler

![Banner](screenshots/interface.png)

## Overview

AU QEC Form Filler is an automation tool designed to streamline the process of filling out various forms on the AU (Air University) QEC (Quality Enhancement Cell) portal. It simplifies the task of completing Course QEC forms, Teacher Evaluation forms, and Online Learning Feedback forms.

Note: Version 4.0 (Go Version) is the currently recommended release for performance and stability.

## Features

- **Automated Login:** Securely logs into the AU QEC portal.
- **Course QEC Forms:** Automatically fills out forms for all available subjects.
- **Teacher Evaluations:** Completes evaluation forms for all listed teachers.
- **Interactive Grading (New):** Assign specific grades (A, B, C, D) to individual teachers via an interactive terminal table.
- **Online Learning Feedback:** Submits feedback forms for all available courses.
- **Cross-Platform:** Runs on Windows, Linux, and macOS.

## Quick Start: Download Binaries (No Setup)

Binaries are available in the Releases section. You can run the tool directly without installing Go, Python, or Docker.

1. Navigate to the [Releases Section](https://github.com/Aw4iskh4n/au-qec-form-filler/releases) on GitHub.
2. Download the binary matching your operating system:
   - **Windows:** `qec-windows.exe`
   - **Linux:** `qec-linux`
   - **macOS:** `qec-mac`
3. Run the file:
   - **Windows:** Double-click `qec-windows.exe`.
   - **Linux:** Open a terminal, grant permission (`chmod +x qec-linux`), and run (`./qec-linux`).
   - **macOS:** Open a terminal, grant permission (`chmod +x qec-mac`), and run (`./qec-mac`).
4. Enter your login credentials when prompted.

## Installation and Usage (For Developers)

If you prefer to run the code from source, you can use Docker, Go, or Python.

### Option 1: Using Docker (Recommended)

This method requires no dependency installation other than Docker.

1. Build and start the container:
    ```bash
    docker-compose up --build -d
    ```
2. Attach to the interactive terminal:
    ```bash
    docker attach au_qec_bot
    ```

### Option 2: Running from Source (Go)

Recommended for speed and the new interactive grading features.

**Requirements:**
- Go 1.21 or higher

**Steps:**

1. Navigate to the `go` directory:
    ```bash
    cd go
    ```

2. Initialize the module and install dependencies:
    ```bash
    go mod init au_portal_bot
    go get [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
    go get [github.com/fatih/color](https://github.com/fatih/color)
    go get [github.com/jedib0t/go-pretty/v6/table](https://github.com/jedib0t/go-pretty/v6/table)
    go get [github.com/jedib0t/go-pretty/v6/text](https://github.com/jedib0t/go-pretty/v6/text)
    go get golang.org/x/term
    ```

3. Run the application:
    ```bash
    go run main.go
    ```

### Option 3: Running from Source (Python)

The legacy version of the script.

**Requirements:**
- Python 3.6+
- Libraries: `requests`, `beautifulsoup4`, `rich`

**Steps:**
1. Install dependencies:
    ```bash
    pip install -r requirements.txt
    ```
2. Run the script:
    ```bash
    python o2mation.py
    ```

![Filled](screenshots/filled.png)

## License

This project is [licensed](https://github.com/Aw4iskh4n/au-qec-form-filler/blob/main/LICIENSE).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

**Disclaimer:** This tool is intended for educational purposes only. Use it responsibly and at your own risk.
