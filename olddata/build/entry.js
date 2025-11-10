"use strict";
(() => {
  // data/entities/maze.ts
  function generateMazeRooms(opts) {
    const width = Math.max(1, Math.min(10, opts?.width ?? 10));
    const height = Math.max(1, Math.min(10, opts?.height ?? 10));
    const idPrefix = opts?.idPrefix ?? "Maze";
    const roomIcon = opts?.roomIcon ?? "O";
    const roomColor = opts?.roomColor ?? "gray";
    const rnd = opts?.rng ?? Math.random;
    const dirs = ["north", "east", "south", "west"];
    const dRow = { north: -1, east: 0, south: 1, west: 0 };
    const dCol = { north: 0, east: 1, south: 0, west: -1 };
    const opposite = { north: "south", south: "north", east: "west", west: "east" };
    const idAt = (r, c) => `${idPrefix}_r${r}_c${c}`;
    const exits = Array.from(
      { length: height },
      () => Array.from({ length: width }, () => ({}))
    );
    const visited = Array.from(
      { length: height },
      () => Array.from({ length: width }, () => false)
    );
    function shuffle(arr) {
      for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(rnd() * (i + 1));
        [arr[i], arr[j]] = [arr[j], arr[i]];
      }
      return arr;
    }
    const stack = [[0, 0]];
    visited[0][0] = true;
    while (stack.length) {
      const [r, c] = stack[stack.length - 1];
      const order = shuffle([...dirs]);
      let progressed = false;
      for (const d of order) {
        const nr = r + dRow[d];
        const nc = c + dCol[d];
        if (nr < 0 || nr >= height || nc < 0 || nc >= width) continue;
        if (visited[nr][nc]) continue;
        exits[r][c][d] = idAt(nr, nc);
        exits[nr][nc][opposite[d]] = idAt(r, c);
        visited[nr][nc] = true;
        stack.push([nr, nc]);
        progressed = true;
        break;
      }
      if (!progressed) stack.pop();
    }
    const rooms = {};
    for (let r = 0; r < height; r++) {
      for (let c = 0; c < width; c++) {
        const id = idAt(r, c);
        rooms[id] = {
          id,
          name: `Maze Room (${r},${c})`,
          description: "Stone walls twist and turn through a tight, grid-cut maze.",
          aliases: [`maze ${r},${c}`, `cell ${r},${c}`],
          tags: ["room", "maze"],
          components: {
            room: {
              icon: roomIcon,
              color: roomColor,
              children: [],
              exits: exits[r][c]
            }
          }
        };
      }
    }
    return rooms;
  }

  // data/entry.ts
  var maze = generateMazeRooms({ width: 8, height: 8, idPrefix: "CoastMaze" });
  Orbis.load(maze);
})();
