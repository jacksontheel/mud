entity LivingRoom {
    component Identity {
        name is "Living Room"
        description is "This is a perfectly nice living room."
        aliases is ["room"]
        tags is ["room"]
    }

    component Room {
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
    component Identity {
        name is "Couch"
        description is "This is a comfy couch."
        aliases is ["couch"]
        tags is ["furniture"]
    }

    trait Kissable

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
    component Identity {
        name is "Nickel"
        description is "A shining nickel, Thomas Jefferson smiles at you from his handsome side profile."
        aliases is ["nickel"]
        tags is ["item"]
    }
}

entity Lamp {
    component Identity {
        name is "Lamp"
        description is "A dimly lit lamp sits in the corner."
        aliases is ["lamp"]
        tags is ["furniture"]
    }

    trait Standard
}

entity BedRoom {
    component Identity {
        name is "Bed Room"
        description is "A fun little bedroom."
        aliases is ["room"]
        tags is ["room"]
    }

    component Room {
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
    component Identity {
        name is "Bed"
        description is "A bed is well-made and looks inviting."
        aliases is ["bed"]
        tags is ["furniture"]
    }

    trait Standard
}

entity Bathroom {
    component Identity {
        name is "Bathroom"
        description is "A bathroom, a perfect place to relax and excrete."
        aliases is ["room"]
        tags is ["room"]
    }

    component Room {
        exits is {
            "west": "LivingRoom",
            "east": "MedicineCabinet"
        }

        children is [
            "Toilet",
            "Goblin"
        ]
    }
}

entity Goblin {
   component Identity {
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
        move target to source Inventory
    }
}

entity Toilet {
    component Identity {
        name is "Toilet"
        description is "A toilet, you piss and poop in here."
        aliases is ["toilet"]
        tags is ["furniture"]
    }   

    trait Standard
}

trait Standard {
    trait Kissable
    trait Hittable
}

trait Kissable {
    when kiss {
        print source "You kiss the {target}"
        publish "{source} kisses the {target}."
    }
}

trait Hittable {
    when attack {
        print source "You hit the {target}"
        publish "{source} hits the {target}."
    }
}