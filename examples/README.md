# Gio Examples

This directory contains example applications demonstrating how to use Gio.

## Hello World

The `hello` example shows a basic Gio application that displays a simple "Hello, Gio!" message.

### Running the Example

#### Desktop (Linux/X11)
```bash
go run ./examples/hello
```

#### Web (WASM)
```bash
webgio ./examples/hello
```

The `webgio` command will:
1. Compile the example to WebAssembly
2. Create a `dist` directory with the necessary files
3. Start a local HTTP server
4. Open your web browser to view the application

### Using webgio

The `webgio` command is a tool for compiling Gio applications to WebAssembly and serving them locally.

```bash
# Basic usage
webgio ./examples/hello

# Custom port
webgio -port 3000 ./examples/hello

# Custom output directory
webgio -dist ./build ./examples/hello

# Don't open browser automatically
webgio -open=false ./examples/hello

# Show help
webgio -help
```
