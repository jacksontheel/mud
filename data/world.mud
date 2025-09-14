entity LivingRoom {
    has Identity {
        name is "Living Room"
        description is "This is a perfectly nice living room."
        aliases is ["room"]
        tags is ["room"]
    }

    has Room {
        exits is {
            "north": "a"
        }

        children is [
            "Couch"
        ]
    }
}

entity Couch {
    has Identity {
        name is "Couch"
        description is "This is a comfy couch."
        aliases is ["couch"]
        tags is ["furniture"]
    }
}