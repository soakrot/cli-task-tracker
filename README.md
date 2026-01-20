# CLI Task Tracker

A simple and lightweight **command-line task tracker** written in **Go**.  
It allows you to manage tasks directly from your terminal with support
for adding, updating, deleting, listing, and marking task statuses.

---

## Features

- Add new tasks
- Update existing tasks
- Delete tasks
- List tasks (with optional status filtering)
- Mark tasks as:
  - `todo`
  - `in-progress`
  - `done`

---

## Command Structure

All commands follow this format:

```bash

cli-task-tracker <command> <arguments>

```

## Usage

### Add a Task

Adds a new task with the specified title.

```bash

tt add <task-title>

```

**Example**:

```bash

tt add "Write README file."

```

### Update a Task


