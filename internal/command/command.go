package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Add struct {
	Title string
}

type Update struct {
	Id    int
	Title string
}

type Delete struct {
	Id int
}

type DeleteAll struct {
	Status string
}

type List struct {
	Status string
}

type Mark struct {
	Id     int
	Status string
}

func AddCmd(args []string) *Add {
	add := &Add{
		Title: "",
	}
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.Parse(os.Args[2:])
	add.Title = addCmd.Arg(0)

	return add
}

func UpdateCmd(args []string) *Update {
	update := &Update{}
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateCmd.Parse(args)

	update.Title = updateCmd.Arg(1)
	id, err := strconv.Atoi(updateCmd.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	update.Id = id

	return update
}

func DeleteCmd(args []string) *Delete {
	del := &Delete{}
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.Parse(args)

	id, err := strconv.Atoi(deleteCmd.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	del.Id = id

	return del
}

func DeleteAllCmd(args []string) *DeleteAll {
	delAll := &DeleteAll{}
	delAllCmd := flag.NewFlagSet("delete-all", flag.ExitOnError)
	delAllCmd.StringVar(&delAll.Status, "status", "", "status[todo, in-progress, done], status is optional")
	delAllCmd.Parse(args)

	return delAll
}

func MarkCmd(args []string) *Mark {
	mark := &Mark{}
	markCmd := flag.NewFlagSet("mark", flag.ExitOnError)
	markCmd.Parse(args)
	status := markCmd.Arg(1)
	id, err := strconv.Atoi(markCmd.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mark.Id = id
	mark.Status = status

	return mark
}

func ListCmd(args []string) *List {
	list := &List{}
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.Parse(args)
	list.Status = listCmd.Arg(0)

	return list
}
