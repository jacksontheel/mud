-- data/rooms.lua
local Man = {
    name = "man",
    description = "A humble man, wary but kind.",
    aliases = { "man", "fellow" },
    tags = { "npc", "human" }
}

return {
    Man = Man,
    Room = {
        name = "Blue Room",
        description = "A calm blue chamber with echoing walls.",
        aliases = { "room", "blue room" },
        tags = { "safe", "indoors" },
        components = {
            room = {
                color = "blue",
                icon = "O",
                children = { Man },
                exits = { north = "Room2", south = "Room3" }
            }
        }
    }
}
