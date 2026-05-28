package main

import (
	"fmt"
	"os"

	"github.com/findmesh/findmesh/backend/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Fprintf(os.Stdout, "FindMesh admin CLI\nAPI base: %s\nUse FINDMESH_ADMIN_TOKEN for privileged API calls.\n", cfg.APIBaseURL)
}
