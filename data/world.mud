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
        print source "As you beat upon the couch, a nickel falls out."
        publish "{source} beats upon the couch, and a shining nickel falls out from under a cushion."
        copy "Nickel" to room Room
    }
}

entity Nickel {
    has Identity {
        name is "Nickel"
        description is "A shining nickel, Thomas Jefferson smiles at you from his handsome side profile."
        aliases is ["nickel"]
        tags is ["item"]
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
            "Toilet",
            "Goblin"
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

entity Goblin {
   has Identity {
        name is "Goblin"
        description is "A funny goblin man no bigger than your fist smiles warmly."
        aliases is ["goblin", "man"]
        tags is ["npc"]
    } 

    when attack by #player {
        print source "You pummel the goblin as he whines helplessly"
        publish "{source} punches and kicks the helpless goblin in his head."
    }

    when kiss by #player {
        print source "You give the goblin a kiss upon his sweaty brow, and he hops into your pocket."
        publish "{source} gives the goblin a kiss, before the goblin jumps into {source}'s pocket."
        move source to source Inventory
    }
}

