package cmd

import (
    "fmt"
    "os"

    "aibuilder/internal/config"
    "aibuilder/internal/logger"
    "aibuilder/internal/task"

    "github.com/spf13/cobra"
)

var (
    toolDescription string
    debugMode       bool
)

var rootCmd = &cobra.Command{
    Use:   "aibuilder",
    Short: "A CLI tool that lets AI build tools in Golang",
    Run: func(cmd *cobra.Command, args []string) {
        log := logger.NewLogger(debugMode)

        cfg, err := config.LoadConfig()
        if err != nil {
            log.Fatalf("Failed to load config: %v", err)
        }

        taskManager := task.NewManager(cfg, log, debugMode)

        err = taskManager.Start(toolDescription)
        if err != nil {
            log.Fatalf("Task execution failed: %v", err)
        }
    },
}

func init() {
    rootCmd.Flags().StringVarP(&toolDescription, "description", "d", "", "Description of the tool to build")
    rootCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "Enable debug mode")
    rootCmd.MarkFlagRequired("description")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
