package main

import (
	"fmt"

	"gio.mleku.dev/widget/material"
)

func main() {
	t := material.NewTheme()
	fmt.Printf("Background: %+v\n", t.Background())
	fmt.Printf("OnSurface: %+v\n", t.OnSurface())
	fmt.Printf("Primary: %+v\n", t.Primary())
	fmt.Printf("Surface: %+v\n", t.Surface())
}
