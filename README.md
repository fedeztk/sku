# sku
Simple TUI written in go to play sudoku in the terminal

[![GO](https://github.com/fedeztk/sku/actions/workflows/go.yaml/badge.svg)](https://github.com/fedeztk/sku/tree/master/.github/workflows/go.yml) [![GHCR](https://github.com/fedeztk/sku/actions/workflows/deploy.yaml/badge.svg)](https://github.com/fedeztk/sku/tree/release/.github/workflows/deploy.yml) [![AUR](https://img.shields.io/aur/version/sku-git?logo=archlinux)](https://aur.archlinux.org/packages/sku-git) [![Go Report Card](https://goreportcard.com/badge/github.com/fedeztk/sku)](https://goreportcard.com/report/github.com/fedeztk/sku)

## Table of Contents

[Usage](#orgfa2aa9c) -
[Features](#org26baa6c) -
[Testing](#org2744438)

sku is a simple TUI for playing sudoku inside the terminal. It uses the awesome [bubbletea](https://github.com/charmbracelet/bubbletea) TUI library for the UI and the [sudoku-go](https://github.com/forfuns/sudoku-go) library as the sudoku logic backend.

Screenshots [here](#org26baa6c)

> Disclaimer: there are probably many other sudoku TUIs around with all kind of features that sku does not have. PRs are welcomed but it is generally a better idea to just use those programs, since adding some features, like a pencil mode used to annotate the sudoku, would require too much effort.


<a id="orgab62fc1"></a>

# Usage

- Install `sku`: 

With the `go` tool:
```sh
go install github.com/fedeztk/sku/cmd/sku@latest
```
**Or** from source:
```sh
# clone the repo
git clone https://github.com/fedeztk/sku.git
# install manually 
make install
```
In both cases make sure that you have the go `bin` directory in your path:
```sh
export PATH="$HOME/go/bin:$PATH"
```
If you are an Arch user there is also an AUR package available:
```sh
paru -S sku-git
```
- Run it interactively:
```sh
sku            # use default mode (easy)
```
For more information about version and modes check the help (`sku -h`)
<a id="org26baa6c"></a>

# Features

- Minimal/clean interface, only displays the board, the help and the game state (timer and remaining cells). Cursor position is marked with green. Greyed cells are unmodifiable (the base of the sudoku); when the cursor is under those cells it will darken
![Screenshot from 2022-06-06 17-35-14](https://user-images.githubusercontent.com/58485208/172200830-677fb8f4-ea29-455c-989b-1c5ea774ae78.png)

- Game check:
	- when the last cell is filled, `sku` will perform a check of the sudoku:
		- if it is correct, an animation will let you know that you won the game
		
https://user-images.githubusercontent.com/58485208/172200661-78ce055f-b5b9-44aa-bf4d-a27e9f8fce85.mp4

    - otherwise it will color with red the errors 
   ![Screenshot from 2022-06-06 17-36-56](https://user-images.githubusercontent.com/58485208/172201574-e1ebe9ec-fc44-4d6c-a80a-287c8433133d.png)


- Simple keys to interact with the puzzle:
	- moving around: use the arrows or the vim motion keys, as preferred
	- setting a cell: just press the desired number
	- unsetting a cell: press spacebar or enter
	- toggle help: press the question mark
	- quit: you can quit anytime by pressing esc or q
	![Screenshot from 2022-06-06 17-36-19](https://user-images.githubusercontent.com/58485208/172201555-1dcc1851-6853-4760-ac7a-6d5285c2f0b6.png)


- Check the version with `sku -v`
- Get the help menu with `sku -h`
- Set a sudoku mode with `sku MODE`. Valid MODEs are: easy, medium, hard and expert (also displayed with `-h` flag)


<a id="org2744438"></a>

# Testing

Development is done through `docker`, build the container with:

    make docker-build

Check that the build went fine:

    docker images | grep sku

Test it with:

    make docker-run

Pre-built Docker image available [here](https://github.com/fedeztk/sku/pkgs/container/sku)
