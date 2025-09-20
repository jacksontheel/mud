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
        print source "You hit the egg upon the couch, gently, as to not disturb the egg."
        publish "{source} hits their egg upon the couch, smiling vacantly to themselves."
    }

    when attack {
        print source "The springs of the couch groan as you beat upon them."
        publish "{source} beats upon the couch, its old springs groaning."
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
        print source "You figure the lamp is roughly the dimensions of a person... you give it a kiss."
        publish "{source} grabs the lamp and plants a passionate kiss upon it."
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