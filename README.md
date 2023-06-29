# Tetris
This project is my first Go project, which is a Tetris game implemented using the Ebiten library. It was developed as a learning project to explore the Go programming language and gain hands-on experience with its concepts and features.

## Installation
1. __Clone the repository__:
      ```shell
   git clone https://github.com/alexandengstrom/tetris.git
      ```
2. __Navigate to the project directory__:
      ```shell
   cd tetris
      ```    
3. __Download the project dependencies using the go.mod file__:
      ```shell
   go mod download
      ```

3. __Build the project__:
     ```shell
   ./build.sh
      ```
## Usage
1. Run the compiled binary to start the Tetris game:
    ```shell
   ./Tetris
      ```
2. Game controls:
   * __Left Arrow__: Move the current block to the left.
   * __Right Arrow__: Move the current block to the right.
   * __Up Arrow__: Rotate the current block.
   * __Down Arrow__: Drop the current block.
3. Objective:
     * Clear lines by filling them completely with blocks to score points.
     * Prevent the blocks from reaching the top of the playfield.
     * The game ends when the blocks stack up to the top.
  
  ##  Pictures

![Screenshot from 2023-06-29 10-03-23](https://github.com/alexandengstrom/tetris/assets/123507241/2f2605b0-1488-4fea-a904-6ab4c5292f27)

