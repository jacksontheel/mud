const _entities = new Map<string, EntityDef>()

// Ambient function authors will call in their entity files:
;(globalThis as any).registerEntity = function(def: EntityDef): EntityDef {
  if (!def?.id) throw new Error("registerEntity: 'id' is required")
  if (_entities.has(def.id)) {
    throw new Error(`Entity '${def.id}' already registered`)
  }
  // store a clone to avoid accidental later mutation
  _entities.set(def.id, { ...def })
  return def
}

// Utility for entry.ts to produce a plain object
export function collect(): Record<string, EntityDef> {
  const out: Record<string, EntityDef> = {}
  for (const [id, ent] of _entities) out[id] = ent
  return out
}
