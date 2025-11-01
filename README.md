# Orbis Mud Engine

Orbis Mud Engine loads a collection of definition files written in the Orbis Definition Language, a domain-specific language designed for expressing entities, components, traits, and commands.
Once compiled, Orbis runs a MUD server that players can connect to via telnet or custom clients.

The Orbis Definition Language defines every part of a world — rooms, objects, creatures, traits, and their relationships — in a human-readable format. Someone wholly unfamiliar with programming can create and edit a game world.

## Installing and Running

1. Clone the repository: `git clone https://github.com/orbis-suite/orbis-mud-engine.git`

2. Run `go run .` from within the newly-cloned directory

3. All of the Orbis Definition Languagefiles are contained within the `data` folder. Edit them as you like.

4. Edit the `config.yaml` file at the root of the repository to change how the game engine handles your Orbis files.

## Orbis Definition Language
### Entities

Every noun in Orbis is an entity, whether it’s a room, an object, an NPC, or even a player. An entity only needs three pieces:
A name, something unique to the entity.
A description, what the entity looks, or sounds, or feels like.
One or more aliases -- something a user can include in a command that the engine can recognize.

Tags are optional, but can be helpful to group items.

Here’s a simple example of two entities, one is a room, and an item to go in that room.

```
entity LivingRoom {
    name is "Living Room"
    description is "A warm and orderly space, its quiet calm inviting rest."
    aliases is ["room", “area”]
    tags is ["room"]

    component Room {
        exits is {
            "north": "Kitchen",
            "west": "Hallway"
        }

        children is [
            "Couch",
            "Lamp"
        ]
    }
}

entity Couch {
    name is "Couch"
    description is "A soft, inviting couch rests here, its cushions sagged with long use."
    aliases is ["couch"]
    tags is ["furniture"]
}
```
### Components

In the example above, you’ve already seen one component, Room. Components define structure and capabilities of an entity, and an entity can have as many components as you want.

Another component is Container. Here’s an entity which contains other children items. For more details on Container fields and other components, check out the wiki, once it’s been made.

```
entity Box {
    name is “Box”
    Description is “A big old cardboard box, filled with goodies.”
    aliases is [“box”, “container”]
    tags is [“furniture”]
 
    component Container {
        prefix is "Inside the box:"
        revealed is true
        children is [
            "Book",
            "Shoe"
        ]
    }
}
```
### Reactions

Now that you have an entity, you can define how that entity reacts to different actions a player might make against it. Let’s say a player attacks the couch we defined earlier, what happens next? You can have as many reactions as you want, based on certain conditions. Then, in each reaction, you can have one or more actions to take. For a list of conditions and actions, check out the wiki, once it’s been named.

```
entity Couch {
    name is "Couch"
    description is "A soft, inviting couch rests here, its cushions sagged with long use."
    aliases is ["couch"]
    tags is ["furniture"]

    react attack {
        when {
            source has tag "player"
        } then {
            print source "You hit the couch with vicious fury."
            publish "{source} takes a mad swing at the couch."
        }
    }
}
```
### Traits

Traits allow you to define generic behavior that you can reuse across multiple components. Let’s take the example above, and have it use a trait instead. We’ll add a few more “when” blocks to the trait, so it’s more expressive.

```
entity Couch {
    name is "Couch"
    description is "A soft, inviting couch rests here, its cushions sagged with long use."
    aliases is ["couch"]
    tags is ["furniture"]

    trait Attackable
}

trait Attackable {
    react attack {
        when {
            instrument is target
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
```
### Commands

Every verb in an Orbis-defined world is a command. Just by defining a command and adding an entity or two’s reaction to that command, you can add another dimension to how a player can interact with your world. Attack comes with the standard library, but this is what it looks like. You can add multiple patterns to a command, so a user can write it out in whichever way they want.

```
command Attack {
    aliases is ["attack", "hit", "beat"]

    pattern {
        syntax is "attack {target}"
        noMatch is "You don't want to attack {target}."
    }

    pattern {
        syntax is "attack {target} with {instrument}"
        noMatch is "You don't want to attack {target} with {instrument}."
    }
}
```
