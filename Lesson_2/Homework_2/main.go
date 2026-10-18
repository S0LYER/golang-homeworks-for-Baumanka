package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"hw02/task"
)

var (
	errEmptyName    = errors.New("пустое имя папки")
	errNoSuchFolder = errors.New("папка отсутствует")
	errSuchFolder   = errors.New("папка с таким именем уже существует")
	errNoSuchTask   = errors.New("задачи нет")
)

func findFolder(folders []string, name string) (int, error) {
	for i, f := range folders {
		if f == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%w: %s", errNoSuchFolder, name)
}

func openFolder(folders []string, n int) (string, error) {
	if n < 1 || n > len(folders) {
		return "", fmt.Errorf("%w: номер %d", errNoSuchFolder, n)
	}
	return folders[n-1], nil
}

func addFolder(folders []string, name string) ([]string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return folders, errEmptyName
	}
	if _, err := findFolder(folders, name); err == nil {
		return folders, fmt.Errorf("%w: %s", errSuchFolder, name)
	}
	return append(folders, name), nil
}

func deleteFolder(folders []string, name string) ([]string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return folders, errEmptyName
	}
	i, err := findFolder(folders, name)
	if err != nil {
		return folders, err
	}
	return append(folders[:i], folders[i+1:]...), nil
}

func renameFolder(folders []string, oldName, newName string) ([]string, error) {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return folders, errEmptyName
	}
	i, err := findFolder(folders, oldName)
	if err != nil {
		return folders, err
	}
	if oldName != newName {
		if _, err := findFolder(folders, newName); err == nil {
			return folders, fmt.Errorf("%w: %s", errSuchFolder, newName)
		}
	}
	folders[i] = newName
	return folders, nil
}

func findTask(list []*task.Task, n int) (*task.Task, error) {
	if n < 1 || n > len(list) {
		return nil, fmt.Errorf("%w: номер %d", errNoSuchTask, n)
	}
	return list[n-1], nil
}

func readLine(in *bufio.Scanner, prompt string) (string, bool) {
	fmt.Print(prompt)
	if !in.Scan() {
		return "", false
	}
	return strings.TrimSpace(in.Text()), true
}

func toChoice(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, errors.New("нужен номер пункта")
	}
	return n, nil
}

func main() {
	folders := []string{"Входящие", "Работа", "Личное"}
	tasks := make(map[string][]*task.Task)
	for _, f := range folders {
		tasks[f] = nil
	}

	current := "" 
	in := bufio.NewScanner(os.Stdin)

	for {
		if current == "" {
			fmt.Println()
			fmt.Println("=== папки ===")
			for i, f := range folders {
				fmt.Printf("%d. %s\n", i+1, f)
			}
			fmt.Println("0. Выход")

			line, ok := readLine(in, "> ")
			if !ok {
				fmt.Println("пока!")
				return
			}
			choice, err := toChoice(line)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}

			switch {
			case choice == 0:
				fmt.Println("пока!")
				return
			default:
				name, err := openFolder(folders, choice)
				if err != nil {
					fmt.Println("ошибка:", err)
					continue
				}
				current = name
			}
			continue
		}

		fmt.Println()
		fmt.Println("[" + current + "]")
		if len(tasks[current]) == 0 {
			fmt.Println("задач нет — создайте первую")
		} else {
			for i, t := range tasks[current] {
				fmt.Printf("%d. %s (создана %s)\n", i+1, t, t.CreatedAtString())
			}
		}
		fmt.Println()
		fmt.Println("1. Создать задачу")
		fmt.Println("2. Выполнить задачу")
		fmt.Println("3. Создать папку")
		fmt.Println("4. Удалить папку")
		fmt.Println("5. Переименовать папку")
		fmt.Println("6. Назад")
		fmt.Println("0. Выход")

		line, ok := readLine(in, "> ")
		if !ok {
			fmt.Println("пока!")
			return
		}
		choice, err := toChoice(line)
		if err != nil {
			fmt.Println("ошибка:", err)
			continue
		}

		switch choice {
		case 0:
			fmt.Println("пока!")
			return
		case 1:
			title, ok := readLine(in, "заголовок задачи: ")
			if !ok {
				fmt.Println("пока!")
				return
			}
			t, err := task.New(title) 
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			tasks[current] = append(tasks[current], t)
			fmt.Println("создана:", t)
		case 2:
			raw, ok := readLine(in, "номер задачи: ")
			if !ok {
				fmt.Println("пока!")
				return
			}
			n, err := strconv.Atoi(raw)
			if err != nil {
				fmt.Println("ошибка: нужен номер задачи")
				continue
			}
			t, err := findTask(tasks[current], n)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			if err := task.MarkDone(t); err != nil {
				fmt.Println("ошибка:", err) 
				continue
			}
			fmt.Println("выполнена:", t)
		case 3:
			name, ok := readLine(in, "имя новой папки: ")
			if !ok {
				fmt.Println("пока!")
				return
			}
			folders, err = addFolder(folders, name)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			tasks[strings.TrimSpace(name)] = nil
			fmt.Println("папка создана:", strings.TrimSpace(name))
		case 4:
			folders, err = deleteFolder(folders, current)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			delete(tasks, current) 
			fmt.Println("папка удалена:", current)
			current = "" 
		case 5:
			name, ok := readLine(in, "новое имя: ")
			if !ok {
				fmt.Println("пока!")
				return
			}
			newName := strings.TrimSpace(name)
			folders, err = renameFolder(folders, current, newName)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			tasks[newName] = tasks[current] 
			delete(tasks, current)
			fmt.Println("папка переименована:", current, "->", newName)
			current = newName
		case 6:
			current = ""
		default:
			fmt.Println("ошибка: нет такого пункта")
		}
	}
}
