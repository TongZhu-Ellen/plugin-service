package main 



import (
	"context"
	"fmt"

	_ "plugin_service/plugins"
	"plugin_service/core"
)






func main() {
	

	// enable
	core.Enable("normal", "v1.0")
	core.Enable("error", "v1.0")
	core.Enable("panic", "v1.0")
	core.Enable("timeout", "v1.0")

	ctx := context.Background()

	// run
	result := core.RunAll(ctx, map[string]any{
		"foo": "bar",
	})

	fmt.Println("==== RESULTS ====")

	 _, _, err1 := core.Status("normal", "v1.0")
	fmt.Println("normal |", result["normal@v1.0"], "|", err1)

	_, _, err2 := core.Status("error", "v1.0")
	fmt.Println("error |", result["error@v1.0"], "|", err2)

	_, _, err3 := core.Status("panic", "v1.0")
	fmt.Println("panic |", result["panic@v1.0"], "|", err3)

	_, _, err4 := core.Status("timeout", "v1.0")
	fmt.Println("timeout |", result["timeout@v1.0"], "|", err4)


	


}