package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	CRUD        CRUD
	Tables      map[string]*DynamicModel
	GopherModel GopherModel
}

func NewApp(crud CRUD, gm GopherModel) *App {
	return &App{
		CRUD:        crud,
		GopherModel: gm,
	}
}

func main() {

	// Initaite Gorm
	db, err := gorm.Open(postgres.Open("postgresql://akinwumi:akin2024@localhost:5432/scriblz"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	gormDB := &GormCRUD{DB: db}
	gopherModel := NewGopherModel(textinput.New())

	app := NewApp(gormDB, *gopherModel)

	if err := app.LoadTables(); err != nil {
		fmt.Println("Error while loading table", err)
	}

	//tea program
	program := tea.NewProgram(app)
	_, err = program.Run()

	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
