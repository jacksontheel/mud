export const man: Entity = {
    id: "Man",
    name: "man",
    description: "man",
    aliases: ["man"],
    reactions: {
        attack: function(ev: Event) {
            ev.print("source", "the man says, ouch!")
        }
    }
}

export const room: Entity = {
    id: "Room",
    name: "Room",
    description: "Room",
    aliases: ["room"],
    components: {
        room: {
            color: "blue",
            icon: "O",
            children: [man],
            exits: {}
        }
    }
}

