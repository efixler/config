package config

import (
	"context"
	"fmt"
	"os"
)

func ExampleSetLoader() {
	f := func(ctx context.Context) Getter {
		os.Setenv("mySpecialConfig","happy")
		return Environment()	
	}
	SetLoader(f)
	config := Default()
	fmt.Println(config.Get("mySpecialConfig"))
	// Output:
	// happy
}
