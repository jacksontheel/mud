entity A1_SunlitEdge {
  has Identity {
    name is "Sunlit Edge"
    description is "Sun-dappled grass at the forest’s edge, bright with tiny wildflowers."
    aliases is ["edge","sunlit edge","clearing"]
    tags is ["room"]
  }
  has Room {
    exits is { "east": "A2_BirchStand", "south": "B1_MossyStones" }
    children is ["A1_Wildflowers", "A1_FlintShard"]
  }
}

entity A1_Wildflowers {
  has Identity {
    name is "Wildflowers"
    description is "Delicate petals in yellow and violet sway in the breeze."
    aliases is ["flowers","wildflowers"]
    tags is ["scenery","flora"]
  }
}

entity A1_FlintShard {
  has Identity {
    name is "Flint Shard"
    description is "A sharp shard of flint, useful if struck just right."
    aliases is ["flint","shard"]
    tags is ["item","tool"]
  }
}

entity A2_BirchStand {
  has Identity {
    name is "Birch Stand"
    description is "White-barked birches cluster here; peeling curls reveal pale wood beneath."
    aliases is ["birches","stand","birch"]
    tags is ["room","forest"]
  }
  has Room {
    exits is { "west": "A1_SunlitEdge", "east": "A3_BrookBend", "south": "B2_CentralGlade" }
    children is ["A2_BirchBark", "A2_Mushrooms"]
  }
}

entity A2_BirchBark {
  has Identity {
    name is "Birch Bark"
    description is "Papery bark curls you could peel off—though the trees might not like it."
    aliases is ["bark"]
    tags is ["scenery","flora"]
  }
}

entity A2_Mushrooms {
  has Identity {
    name is "Cluster of Mushrooms"
    description is "A speckled cluster nestles among damp leaves."
    aliases is ["mushrooms","fungi"]
    tags is ["scenery","flora"]
  }
}

entity A3_BrookBend {
  has Identity {
    name is "Brook Bend"
    description is "A shallow brook bends around smooth stones, whispering over pebbles."
    aliases is ["brook","stream","bend"]
    tags is ["room","water"]
  }
  has Room {
    exits is { "west": "A2_BirchStand", "east": "A4_FernHollow", "south": "B3_HutPorch" }
    children is ["A3_SmoothStone","A3_ReedPatch"]
  }
}

entity A3_SmoothStone {
  has Identity {
    name is "Smooth Stone"
    description is "River-worn, cool to the touch."
    aliases is ["stone","pebble"]
    tags is ["item"]
  }
}

entity A3_ReedPatch {
  has Identity {
    name is "Reed Patch"
    description is "Slender reeds nod where the water slows."
    aliases is ["reeds"]
    tags is ["scenery","flora"]
  }
}

entity A4_FernHollow {
  has Identity {
    name is "Fern Hollow"
    description is "Ferns crowd under the boughs, their fronds brushing your knees."
    aliases is ["hollow","ferns"]
    tags is ["room","forest"]
  }
  has Room {
    exits is { "west": "A3_BrookBend", "south": "B4_HutInterior" }
    children is ["A4_FernFronds"]
  }
}

entity A4_FernFronds {
  has Identity {
    name is "Fern Fronds"
    description is "Coiled fiddleheads ready to unfurl."
    aliases is ["ferns","fronds"]
    tags is ["scenery","flora"]
  }
}

entity B1_MossyStones {
  has Identity {
    name is "Mossy Stones"
    description is "Rounded stones sit in plush moss; the air is cool and damp."
    aliases is ["stones","mossy stones"]
    tags is ["room"]
  }
  has Room {
    exits is { "north": "A1_SunlitEdge", "east": "B2_CentralGlade", "south": "C1_FallenLog" }
    children is ["B1_Moss","B1_Beetle"]
  }
}

entity B1_Moss {
  has Identity {
    name is "Velvet Moss"
    description is "Soft green moss springs back when pressed."
    aliases is ["moss"]
    tags is ["scenery","flora"]
  }
}

entity B1_Beetle {
  has Identity {
    name is "Iridescent Beetle"
    description is "A beetle shimmers with oil-slick colors."
    aliases is ["beetle"]
    tags is ["creature"]
  }
  when kiss {
    say "You pucker up; the beetle wisely scuttles away."
  }
}

entity B2_CentralGlade {
  has Identity {
    name is "Central Glade"
    description is "A circular opening of soft grass—a natural gathering place."
    aliases is ["glade","clearing"]
    tags is ["room"]
  }
  has Room {
    exits is { "north": "A2_BirchStand", "east": "B3_HutPorch", "south": "C2_Footpath", "west": "B1_MossyStones" }
    children is ["B2_WildHerbs","B2_WickerBasket"]
  }
}

entity B2_WildHerbs {
  has Identity {
    name is "Wild Herbs"
    description is "Fragrant sprigs of thyme and mint grow low to the ground."
    aliases is ["herbs","mint","thyme"]
    tags is ["flora","item"]
  }
}

entity B2_WickerBasket {
  has Identity {
    name is "Wicker Basket"
    description is "A woven basket left behind—empty but sturdy."
    aliases is ["basket"]
    tags is ["container","item"]
  }
  when attack {
    say "You stomp the basket flat. That felt… unnecessary."
  }
}

entity B3_HutPorch {
  has Identity {
    name is "Hut Porch"
    description is "A tiny hut with a creaking porch faces the glade."
    aliases is ["porch","hut porch"]
    tags is ["room","structure"]
  }
  has Room {
    exits is { "north": "A3_BrookBend", "east": "B4_HutInterior", "south": "C3_Spring", "west": "B2_CentralGlade" }
    children is ["B3_WoodenChair","B3_HutDoor"]
  }
}

entity B3_WoodenChair {
  has Identity {
    name is "Wooden Chair"
    description is "A simple chair with a smooth-worn seat."
    aliases is ["chair"]
    tags is ["furniture"]
  }
  when kiss by #player {
    say "You kiss the chair’s backrest. It accepts with wooden dignity."
  }
}

entity B3_HutDoor {
  has Identity {
    name is "Hut Door"
    description is "A crooked plank door with a loop of rope for a handle."
    aliases is ["door"]
    tags is ["furniture","door"]
  }
  when attack {
    say "You thump the door; dust trickles from the lintel."
  }
}

entity B4_HutInterior {
  has Identity {
    name is "Hut Interior"
    description is "Cozy and cluttered: a small table, a bedroll, and hanging bundles of herbs."
    aliases is ["hut","interior"]
    tags is ["room","structure"]
  }
  has Room {
    exits is { "north": "A4_FernHollow", "south": "C4_GardenPlot", "west": "B3_HutPorch" }
    children is ["B4_Table","B4_Bedroll","B4_HangingHerbs","B4_Lantern","B4_HenEgg"]
  }
}

entity B4_Table {
  has Identity {
    name is "Small Table"
    description is "A wobbling table carved from a single plank."
    aliases is ["table"]
    tags is ["furniture"]
  }
}

entity B4_Bedroll {
  has Identity {
    name is "Bedroll"
    description is "Blankets bundled neatly on the floor."
    aliases is ["bedroll","bed"]
    tags is ["furniture"]
  }
}

entity B4_HangingHerbs {
  has Identity {
    name is "Hanging Herbs"
    description is "Bundles of drying herbs sway from the rafters."
    aliases is ["herbs","bundles"]
    tags is ["flora","scenery"]
  }
}

entity B4_Lantern {
  has Identity {
    name is "Tin Lantern"
    description is "A pierced tin lantern; faint soot rings the top."
    aliases is ["lantern"]
    tags is ["item","light"]
  }
  when kiss {
    say "Smells vaguely of smoke and old oil."
  }
}

entity B4_HenEgg {
  has Identity {
    name is "Hen’s Egg"
    description is "A speckled egg rests in a little reed bowl."
    aliases is ["egg"]
    tags is ["item","egg"]
  }
}

entity C1_FallenLog {
  has Identity {
    name is "Fallen Log"
    description is "A broad log bridges a shallow dip in the ground."
    aliases is ["log","fallen log"]
    tags is ["room"]
  }
  has Room {
    exits is { "north": "B1_MossyStones", "east": "C2_Footpath", "south": "D1_Thicket" }
    children is ["C1_BarkBeetles"]
  }
}

entity C1_BarkBeetles {
  has Identity {
    name is "Bark Beetles"
    description is "Tiny trails lace the underside of the bark."
    aliases is ["beetles"]
    tags is ["scenery","fauna"]
  }
}

entity C2_Footpath {
  has Identity {
    name is "Footpath"
    description is "A narrow path trodden into the grass, pointing hut-ward."
    aliases is ["path","trail","footpath"]
    tags is ["room","path"]
  }
  has Room {
    exits is { "north": "B2_CentralGlade", "east": "C3_Spring", "south": "D2_WillowShade", "west": "C1_FallenLog" }
    children is ["C2_Feather"]
  }
}

entity C2_Feather {
  has Identity {
    name is "Jay Feather"
    description is "A bright blue feather—a lucky find."
    aliases is ["feather"]
    tags is ["item"]
  }
}

entity C3_Spring {
  has Identity {
    name is "Clear Spring"
    description is "Cold water wells up from the earth, glassy and pure."
    aliases is ["spring","water"]
    tags is ["room","water"]
  }
  has Room {
    exits is { "north": "B3_HutPorch", "east": "C4_GardenPlot", "south": "D3_BrookFord", "west": "C2_Footpath" }
    children is ["C3_WaterBucket"]
  }
}

entity C3_WaterBucket {
  has Identity {
    name is "Wooden Bucket"
    description is "A small bucket with iron hoops; damp inside."
    aliases is ["bucket"]
    tags is ["item","container"]
  }
}

entity C4_GardenPlot {
  has Identity {
    name is "Garden Plot"
    description is "Rows of neat soil, with sprouts just peeking through."
    aliases is ["garden","plot"]
    tags is ["room","structure"]
  }
  has Room {
    exits is { "north": "B4_HutInterior", "south": "D4_BackPath", "west": "C3_Spring" }
    children is ["C4_Sapling","C4_Trowel"]
  }
}

entity C4_Sapling {
  has Identity {
    name is "Sapling"
    description is "A young fruit tree is staked against the wind."
    aliases is ["tree","sapling"]
    tags is ["flora","scenery"]
  }
}

entity C4_Trowel {
  has Identity {
    name is "Garden Trowel"
    description is "A little iron trowel with a smooth wooden handle."
    aliases is ["trowel","spade"]
    tags is ["item","tool"]
  }
}

entity D1_Thicket {
  has Identity {
    name is "Thicket"
    description is "Tangled undergrowth forms a scraggly wall of stems and thorns."
    aliases is ["thicket","brush"]
    tags is ["room","forest"]
  }
  has Room {
    exits is { "north": "C1_FallenLog", "east": "D2_WillowShade" }
    children is ["D1_BerryBush"]
  }
}

entity D1_BerryBush {
  has Identity {
    name is "Berry Bush"
    description is "Dark berries hide beneath serrated leaves."
    aliases is ["berries","bush"]
    tags is ["flora","item"]
  }
  when attack {
    say "A few berries squish in your hand. Sticky."
  }
}

entity D2_WillowShade {
  has Identity {
    name is "Willow Shade"
    description is "Long willow fronds trail like curtains."
    aliases is ["willow","shade"]
    tags is ["room","forest"]
  }
  has Room {
    exits is { "north": "C2_Footpath", "east": "D3_BrookFord", "west": "D1_Thicket" }
    children is ["D2_WillowSwitch"]
  }
}

entity D2_WillowSwitch {
  has Identity {
    name is "Willow Switch"
    description is "A springy, supple branch—thin but strong."
    aliases is ["switch","branch"]
    tags is ["item","tool"]
  }
}

entity D3_BrookFord {
  has Identity {
    name is "Brook Ford"
    description is "Flat stones break the stream for an easy crossing."
    aliases is ["ford","brook"]
    tags is ["room","water"]
  }
  has Room {
    exits is { "north": "C3_Spring", "east": "D4_BackPath", "west": "D2_WillowShade" }
    children is ["D3_WetFootprints"]
  }
}

entity D3_WetFootprints {
  has Identity {
    name is "Wet Footprints"
    description is "Fresh prints glisten on the stones, heading toward the hut."
    aliases is ["footprints","prints"]
    tags is ["scenery","clue"]
  }
}

entity D4_BackPath {
  has Identity {
    name is "Back Path"
    description is "A narrow track behind the hut, half-overgrown."
    aliases is ["path","back path"]
    tags is ["room","path"]
  }
  has Room {
    exits is { "north": "C4_GardenPlot", "west": "D3_BrookFord" }
    children is ["D4_Pinecones"]
  }
}

entity D4_Pinecones {
  has Identity {
    name is "Pinecones"
    description is "A scatter of prickly cones; resin clings to the scales."
    aliases is ["cones","pinecones"]
    tags is ["scenery","flora"]
  }
}
