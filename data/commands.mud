command Attack {
    aliases is ["attack", "hit", "beat"]

    pattern {
        syntax is "{command} {target}"
        noMatch is "blah blah blah"
    }

    pattern {
        syntax is "{command} {target} with {instrument}"
        noMatch is "blah blah blah"
    }
}