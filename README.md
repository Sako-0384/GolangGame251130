# Gopher the Channel Miner

A game project created with TinyGo and TIC-80.

## Requirements

- Go 1.21+
- TinyGo
- TIC-80

## Build Instructions

Run the following commands to build the game and embed it into the TIC-80 cart file.

```bash
tinygo build -o build/main.wasm -target=./target_tic80.json ./main.go
tic80 --fs . --cmd "load game.tic & import binary build/main.wasm & save & exit"
```

export to HTML (creates `build/game.zip`):

```bash
tic80 --fs . --cmd "load game.tic & export html build/game.zip & exit"
```
