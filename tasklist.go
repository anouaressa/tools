package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Content     string    `json:"content"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

var tasks []Task
var nextID int
var filename = "tasks.json"

// Load tasks from the file
func loadTasks() error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, no tasks to load
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return err
	}

	// Set the next ID based on the last task's ID
	if len(tasks) > 0 {
		nextID = tasks[len(tasks)-1].ID + 1
	}
	return nil
}

// Save tasks to the file
func saveTasks() error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(tasks)
}

// Function to list tasks in a table format
func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks to show.")
		return
	}

	// Print header
	fmt.Println("\nID | Task Content           | State     | Created At          | Completed At")
	fmt.Println("-------------------------------------------------------------------------")

	// Print each task in a row
	for _, task := range tasks {
		completedAt := "N/A"
		if task.CompletedAt != nil {
			completedAt = task.CompletedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%-3d | %-23s | %-10s | %-19s | %-19s\n", task.ID, task.Content, task.State, task.CreatedAt.Format("2006-01-02 15:04:05"), completedAt)
	}
}

// Function to add a task
func addTask(content string) {
	task := Task{
		ID:        nextID,
		Content:   content,
		State:     "Pending",
		CreatedAt: time.Now(),
	}
	tasks = append(tasks, task)
	nextID++
	saveTasks() // Save the tasks after adding
	fmt.Printf("Added task: %s\n", content)
}

// Function to remove a task
func removeTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks() // Save the tasks after removal
			fmt.Printf("Removed task: %d\n", id)
			return
		}
	}
	fmt.Println("Task not found.")
}

// Function to update the state of a task
func updateTaskState(id int, state string) {
	for i, task := range tasks {
		if task.ID == id {
			if state == "Completed" && task.State != "Completed" {
				completedAt := time.Now()
				tasks[i].State = "Completed"
				tasks[i].CompletedAt = &completedAt
			} else if state == "Pending" && task.State != "Pending" {
				tasks[i].State = "Pending"
				tasks[i].CompletedAt = nil
			}
			saveTasks() // Save the tasks after updating state
			fmt.Printf("Task %d marked as %s\n", id, state)
			return
		}
	}
	fmt.Println("Task not found.")
}

// Function to display the CLI menu
func showMenu() {
	fmt.Println("\nTask List CLI")
	fmt.Println("1. List tasks")
	fmt.Println("2. Add task")
	fmt.Println("3. Remove task")
	fmt.Println("4. Update task state (P = Pending, C = Completed)")
	fmt.Println("5. Exit")
}

func main() {
	// Load tasks from file when the program starts
	err := loadTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()

		// Read user choice
		fmt.Print("Enter your choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			listTasks()
		case "2":
			// Read task content (including spaces)
			fmt.Print("Enter task content: ")
			scanner.Scan()
			content := scanner.Text()
			addTask(content)
		case "3":
			// Read task ID to remove
			fmt.Print("Enter task ID to remove: ")
			scanner.Scan()
			idStr := scanner.Text()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID.")
				break
			}
			removeTask(id)
		case "4":
			// Read task ID to update state
			fmt.Print("Enter task ID to update state: ")
			scanner.Scan()
			idStr := scanner.Text()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID.")
				break
			}
			fmt.Print("Enter new state (P for Pending, C for Completed): ")
			scanner.Scan()
			state := scanner.Text()

			// Normalize the state input
			if state == "P" || state == "p" || state == "Pending" {
				updateTaskState(id, "Pending")
			} else if state == "C"|| state == "c" || state == "Completed" {
				updateTaskState(id, "Completed")
			} else {
				fmt.Println("Invalid state. Use 'P' for Pending or 'C' for Completed.")
			}
		case "5":
			fmt.Println("Exiting Task List CLI.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
