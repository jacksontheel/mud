command Attack {
    aliases is ["attack", "hit", "beat"]

    pattern {
        syntax is "attack {target}"
        noMatch is "blah blah blah"
    }

    pattern {
        syntax is "attack {target} with {instrument}"
        noMatch is "blah blah blah"
    }
}