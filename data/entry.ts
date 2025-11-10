import { generateMazeRooms } from "./entities/maze";

const maze = generateMazeRooms({ width: 8, height: 8, idPrefix: "CoastMaze" });

Orbis.load(maze)