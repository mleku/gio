// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	var (
		port    = flag.Int("port", 8080, "Port to serve on")
		distDir = flag.String("dist", "dist", "Distribution directory")
		open    = flag.Bool("open", true, "Open browser automatically")
		help    = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No main package specified\n")
		showHelp()
		os.Exit(1)
	}

	mainPkg := args[0]
	if !strings.HasPrefix(mainPkg, "./") && !strings.HasPrefix(mainPkg, "/") {
		mainPkg = "./" + mainPkg
	}

	fmt.Printf("Compiling %s to WASM...\n", mainPkg)

	// Create dist directory if it doesn't exist
	if err := os.MkdirAll(*distDir, 0755); err != nil {
		log.Fatalf("Failed to create dist directory: %v", err)
	}

	// Compile to WASM
	if err := compileToWASM(mainPkg, *distDir); err != nil {
		log.Fatalf("Failed to compile to WASM: %v", err)
	}

	fmt.Printf("WASM compilation complete. Serving from %s on port %d\n", *distDir, *port)

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.FileServer(http.Dir(*distDir)),
	}

	// Open browser if requested
	if *open {
		go func() {
			time.Sleep(1 * time.Second) // Give server time to start
			openBrowser(fmt.Sprintf("http://localhost:%d", *port))
		}()
	}

	// Start server
	fmt.Printf("Server starting at http://localhost:%d\n", *port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}

func compileToWASM(mainPkg, distDir string) error {
	// Set GOOS and GOARCH for WASM
	env := append(os.Environ(), "GOOS=js", "GOARCH=wasm")

	// Compile the main package to WASM
	cmd := exec.Command("go", "build", "-o", filepath.Join(distDir, "main.wasm"), mainPkg)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %v", err)
	}

	// Copy wasm_exec.js from Go installation
	if err := copyWasmExec(distDir); err != nil {
		return fmt.Errorf("failed to copy wasm_exec.js: %v", err)
	}

	// Create a basic index.html if it doesn't exist
	if err := createIndexHTML(distDir); err != nil {
		return fmt.Errorf("failed to create index.html: %v", err)
	}

	return nil
}

func copyWasmExec(distDir string) error {
	// Find Go installation directory
	goRoot := runtime.GOROOT()
	if goRoot == "" {
		return fmt.Errorf("GOROOT not set")
	}

	src := filepath.Join(goRoot, "lib", "wasm", "wasm_exec.js")
	dst := filepath.Join(distDir, "wasm_exec.js")

	// Copy the file
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read wasm_exec.js: %v", err)
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return fmt.Errorf("failed to write wasm_exec.js: %v", err)
	}

	return nil
}

func createIndexHTML(distDir string) error {
	indexPath := filepath.Join(distDir, "index.html")

	// Check if index.html already exists
	if _, err := os.Stat(indexPath); err == nil {
		return nil // File exists, don't overwrite
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Gio WASM App</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background: #f0f0f0;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }
        #loading {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            text-align: center;
        }
        #giowindow {
            width: 100%;
            height: 100vh;
            position: relative;
        }
        canvas {
            position: fixed;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
        }
    </style>
</head>
<body>
    <div id="loading">Loading...</div>
    <div id="giowindow"></div>
    <script src="wasm_exec.js"></script>
    <script>
        console.log("Starting WASM load...");
        console.log("Document ready state:", document.readyState);
        console.log("giowindow element:", document.getElementById("giowindow"));
        
        const go = new Go();
        console.log("Go object created:", go);
        
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            console.log("WASM loaded successfully, running...");
            console.log("WASM instance:", result.instance);
            
            go.run(result.instance);
            console.log("WASM run completed");
            
            // Check if canvas was created
            const canvas = document.querySelector("canvas");
            console.log("Canvas element:", canvas);
            if (canvas) {
                console.log("Canvas dimensions:", canvas.width, "x", canvas.height);
                console.log("Canvas style:", canvas.style.cssText);
            }
            
            console.log("WASM running, hiding loading...");
            document.getElementById("loading").style.display = "none";
        }).catch((err) => {
            console.error("Failed to load WASM:", err);
            document.getElementById("loading").innerHTML = "Failed to load application: " + err.message;
        });
        
        // Additional debugging
        setTimeout(() => {
            console.log("5 seconds later - checking DOM:");
            console.log("giowindow children:", document.getElementById("giowindow").children.length);
            const canvas = document.querySelector("canvas");
            if (canvas) {
                console.log("Canvas found after timeout:", canvas);
            } else {
                console.log("No canvas found after timeout");
            }
        }, 5000);
    </script>
</body>
</html>`

	return os.WriteFile(indexPath, []byte(html), 0644)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		fmt.Printf("Please open %s in your browser\n", url)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please open %s in your browser\n", url)
	}
}

func showHelp() {
	fmt.Printf(`webgio - Compile Gio apps to WASM and serve them

Usage:
    webgio [flags] <main-package>

Flags:
    -port int
        Port to serve on (default 8080)
    -dist string
        Distribution directory (default "dist")
    -open
        Open browser automatically (default true)
    -help
        Show this help message

Examples:
    webgio ./examples/hello
    webgio -port 3000 ./examples/hello
    webgio -dist ./build ./examples/hello

The command will:
1. Compile the specified main package to WASM
2. Copy necessary files to the dist directory
3. Start an HTTP server serving the dist directory
4. Open a web browser to view the application
`)
}
