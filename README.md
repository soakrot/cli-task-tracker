# CLI Task Tracker

A simple and lightweight **command-line task tracker** written in **Go**.  
It allows you to manage tasks directly from your terminal with support
for adding, updating, deleting, listing, and marking task statuses.

---

## Features

- Add, update, list, delete, and mark tasks
- Simple command-based interface
- Supports task statuses: **todo**, **in-progress**, **done**
- Persistent task storage (local)
- Fast, minimal, and dependency-free

---

## Installation

Clone the repository and build the binary:

```sh
go build -o tt

```

## Command Structure

All commands follow this format:

```bash

tt <command> <arguments>

```

## Usage

### Add a Task

Adds a new task with the specified title.

```bash

tt add <task-title>

```

Example:

```bash
tt add "Write README file."

```

### Update a Task

```bash
tt update <id> <new-content>

```

Example:

```bash
tt update 1 "Write detailed README"
```

### Delete a Task

Delete a task by ID.
Returns the deleted task content.

```bash
tt delete <id>
```

Example:

```bash
tt delete 1
# Output: Write detailed README
```

### Delete All Tasks

Delete all tasks. if a status is provided, only tasks with that status are deleted.

```bash
tt delete-all [status]
```

Example:

```bash
tt delete-all
tt delete-all -status done
```

### Mark a Task

Mark a task with a new status.

```bash
tt mark <id> <status>
```

Example:

```bash
tt delete-all [status]
```

### List Tasks

List tasks. If no status is provided, all tasks are listed.

```bash
tt list [status]
```

Example:

```bash
tt list
tt list todo
```

### Example Workflow

```bash
tt add "Learn Go"
tt add "Build a CLI app"
tt mark 1 in-progress
tt list
tt delete 2
```

## License

MIT License
