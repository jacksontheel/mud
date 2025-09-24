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