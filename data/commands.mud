command Attack {
    aliases is ["attack", "hit", "beat"]

    pattern {
        syntax is "attack {target}"
        noMatch is "You don't want to attack that."
    }

    pattern {
        syntax is "attack {target} with {instrument}"
        noMatch is "You don't want to attack that with that."
    }
}

command Kiss {
    aliases is ["smooch"]

    pattern {
        syntax is "kiss {target}"
        noMatch is "you don't want to kiss that."
    }
}

command Take {
    aliases is ["take", "grab", "pickup"]
    
    pattern {
        syntax is "take {target}"
        noMatch is "you can't pick that up."
    }
}

command Drop {
    aliases is ["drop"]

    pattern {
        syntax is "drop {target}"
        noMatch is "you can't drop that."
    }
}

command Give {
    aliases is ["give", "hand"]

    pattern {
        syntax is "give {instrument} to {target}"
        noMatch is "You can't give that to that."
    }

    pattern {
        syntax is "give {target} {instrument}"
        noMatch is "You can't give that to that."
    }
}