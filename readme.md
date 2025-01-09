# Command Line Pokedex

A terminal-based Pokedex application that allows users to catch and inspect Pokemon. This is a guided project from 
[Boot.dev](https://www.boot.dev) that focused on creating a CLI app using go while learning how to use JSON, make 
network request, and implement caching. 

## Features

* list Pokemon location areas
* Explore location area by name
* Capture pokemon 
  * Capture rate scales down as base experience of Pokemon increases
* Inspect Pokemon you've captured
* List all Pokemon discovered 
* Caches requests to the [Pokemon API](https://pokeapi.co/docs/v2)
* Basic help documentation

## Commands

- `help`: Displays a help message
- `map`: Displays list of location areas, each subsequent call will return the next page of location areas
- `mapb`: Displays list of location areas, each subsequent call will return the previous page of location areas
- `explore <location-area>`: Displays list of pokemon at given location
- `catch <pokemon>`: Attempts to catch designated pokemon
- `inspect <pokemon>`: Displays information of captured pokemon
- `pokedex`: Displays list of pokemon that have been captured
- `exit`: Exit the Pokedex


## Example Usage

```bash
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!

Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  -defense: 40
  -special-attack: 35
  -special-defense: 35
  -speed: 56
Types:
  - normal
  - flying
```

## Setup

### From Source

```bash
## Clone repo
git clone https://github.com/ZorahScope/pokedexcli.git
## Navigate into project directory
cd pokedexcli
## Build binary from project repo
go build .
## run application
./pokedexcli

## Alternative to the build and run step by running directly
go run .
```