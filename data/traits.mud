trait Standard {
    trait Kissable
    trait Hittable
}

trait Kissable {
    react kiss {
        then {
            print source "You kiss the {target}"
            publish "{source} kisses the {target}."
        }
    }
}

trait Hittable {
    react attack {
        when {
            not not instrument is target
        } then {
            print source "You can't hit something with itself."
        }

        when {
            instrument exists
        } then {
            print source "You hit the {target} with {instrument}"
            publish "{source} hits the {target} with {instrument}."
        }

        then {
            print source "You hit the {target}"
            publish "{source} hits the {target}."
        }
    }
}