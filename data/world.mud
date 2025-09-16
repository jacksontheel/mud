entity LivingRoom {
    has Identity {
        name is "Living Room"
        description is "This is a perfectly nice living room."
        aliases is ["room"]
        tags is ["room"]
    }

    has Room {
        exits is {
            "north": "BedRoom",
            "east": "Bathroom"
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

    when attack by #player with #egg {
        say "You hit the egg upon the couch, gently, as to not disturb the egg."
    }

    when attack {
        say "As you beat the couch, a nickel falls out from under the cushion."
    }
}


entity Lamp {
    has Identity {
        name is "Lamp"
        description is "A dimly lit lamp sits in the corner."
        aliases is ["lamp"]
        tags is ["furniture"]
    }

    when kiss by #player {
        say "You figure the lamp is roughly the dimensions of a person... you give it a kiss."
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
            "Bed",
            "Lamp"
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

entity Bathroom {
    has Identity {
        name is "Bathroom"
        description is "A bathroom, a perfect place to relax and excrete."
        aliases is ["room"]
        tags is ["room"]
    }

    has Room {
        exits is {
            "west": "LivingRoom"
        }

        children is [
            "Toilet"
        ]
    }
}

entity Toilet {
 has Identity {
        name is "Toilet"
        description is "A toilet, you piss and poop in here."
        aliases is ["toilet"]
        tags is ["furniture"]
    }   
}