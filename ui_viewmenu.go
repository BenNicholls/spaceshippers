package main

import "github.com/bennicholls/burl-E/burl"

type ViewMenu struct {
	burl.PagedContainer
}

func NewViewMenu() (vm *ViewMenu) {
	vm = new(ViewMenu)
	vm.PagedContainer = *burl.NewPagedContainer(40, 28, 39, 3, 10, true)

	vm.SetVisibility(false)

	return
}
