package db

import (
	"log"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// InitPocketBase initializes an embedded PocketBase server.
// This allows Peregrine to use PocketBase as its core auth/data layer.
func InitPocketBase(dataDir string) *pocketbase.PocketBase {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: filepath.Join(dataDir, "pb_data"),
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		log.Println("PocketBase server started inside Peregrine!")
		return se.Next()
	})

	// Normally we would app.Start(), but since this is embedded in a TUI,
	// we will just return the app for now. 
	// In a real scenario, we'd start it in a background goroutine.
	return app
}
