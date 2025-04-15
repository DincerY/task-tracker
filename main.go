package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// JSON olayını burada yapabilirim belki
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := os.ReadFile("tasks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tasks := []Task{}
	if len(data) != 0 {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage -> task-tracker <command>")
		return
	}

	switch os.Args[1] {
	case "add":
		fmt.Println("Before Task Lists : ")
		list()
		fmt.Println("------------------")

		if len(os.Args) <= 2 {
			fmt.Println("Please provide a task to add")
			return
		}
		var newTask Task
		if len(tasks) == 0 {
			newTask.ID = 1
		} else {
			newTask.ID = tasks[len(tasks)-1].ID + 1
		}
		newTask.Description = os.Args[2]
		newTask.Status = "todo"
		newTask.CreatedAt = time.Now()
		add(newTask)

		fmt.Println("------------------")
		fmt.Println("After Task Lists : ")
		list()
	case "update":
		fmt.Println("Before Task Lists : ")
		list()
		fmt.Println("------------------")

		if len(os.Args) <= 3 {
			fmt.Println("Please provide a task to update")
			//örnek yazabiliriz
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID is not int:", os.Args[2])
			return
		}
		description := os.Args[3]

		update(id, description)

		fmt.Println("------------------")
		fmt.Println("After Task Lists : ")
		list()
	case "delete":
		fmt.Println("Before Task Lists : ")
		list()
		fmt.Println("------------------")

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID is not int:", os.Args[2])
			return
		}

		delete(id)

		fmt.Println("------------------")
		fmt.Println("After Task Lists : ")
		list()
	case "list":
		if len(os.Args) == 2 {
			list()
			return
		}
		status := os.Args[2]
		switch status {
		case "todo":
			listWithStatus(status)
		case "in-progress":
			listWithStatus(status)
		case "done":
			listWithStatus(status)
		default:
			fmt.Println("Unknown status")
		}
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID is not int:", os.Args[2])
			return
		}
		markInProgress(id)

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task ID")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("ID is not int:", os.Args[2])
			return
		}
		markDone(id)
	default:
		fmt.Println("Unknown command")
	}
}

func add(newTask Task) {
	tasks := readFileConvertToTasks()
	tasks = append(tasks, newTask)
	writeFile(tasks)
	fmt.Println("Task added successfully ID:", newTask.ID)
}

func delete(id int) {
	tasks := readFileConvertToTasks()
	if id < 1 || id > len(tasks) {
		fmt.Println("Invalid index")
		return
	}
	index := -1
	for i, t := range tasks {
		if t.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Task id not found")
		return
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	writeFile(tasks)
}

func update(id int, description string) {
	tasks := readFileConvertToTasks()
	if id < 1 || id > len(tasks) {
		fmt.Println("Invalid index")
		return
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			writeFile(tasks)
			fmt.Println("Task updated successfully")
			break
		}
	}
}

func markInProgress(id int) {
	tasks := readFileConvertToTasks()
	if id < 1 || id > len(tasks) {
		fmt.Println("Invalid index")
		return
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = "in-progress"
			tasks[i].UpdatedAt = time.Now()
			writeFile(tasks)
			fmt.Println("Task marked as in-progress")
			return
		}
	}
}

func markDone(id int) {
	tasks := readFileConvertToTasks()
	if id < 1 || id > len(tasks) {
		fmt.Println("Invalid index")
		return
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = "done"
			tasks[i].UpdatedAt = time.Now()
			writeFile(tasks)
			fmt.Println("Task marked as done")
			return
		}
	}
}

func list() {
	tasks := readFileConvertToTasks()
	if tasks == nil {
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	for _, t := range tasks {
		fmt.Printf("Id: %d, Description: %s, Status: %s\n", t.ID, t.Description, t.Status)
	}
}

func listWithStatus(status string) {
	tasks := readFileConvertToTasks()
	if tasks == nil {
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	for _, t := range tasks {
		if t.Status == status {
			fmt.Printf("Id: %d, Description: %s, Status: %s\n", t.ID, t.Description, t.Status)
		}
	}
}

func readFileConvertToTasks() []Task {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	tasks := []Task{}
	if len(data) != 0 {
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return nil
		}
	}
	return tasks
}

func writeFile(tasks []Task) {
	data, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile("tasks.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("File written successfully")
}
