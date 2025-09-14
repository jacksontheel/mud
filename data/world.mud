entity LivingRoom {
    has Identity {
        name is "Living Room"
        description is "This is a perfectly nice living room."
        aliases is ["room"]
        tags is ["room"]
    }

    has Room {
        exits is {
            "north": "BedRoom"
        }

        children is [
            "Couch",
            "Lamp"
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

entity Lamp {
    has Identity {
        name is "Lamp"
        description is "A dimly lit lamp sits in the corner."
        aliases is ["lamp"]
        tags is ["furniture"]
    }
}

entity BedRoom {
    has Identity {
        name is "Bed Room"
        description is "A fun little bedroom."
        aliases is ["room"]
        tags is ["room"]
    }

    has Room {
        exits is {
            "south": "LivingRoom"
        }

        children is [
            "Bed"
        ]
    }
}

entity Bed {
    has Identity {
        name is "Bed"
        description is "A bed is well-made and looks inviting."
        aliases is ["bed"]
        tags is ["furniture"]
    }
}