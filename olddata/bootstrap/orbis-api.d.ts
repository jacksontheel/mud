export {}

declare global {
  type OrbisComponentMap = {
    room: {
      color: string
      icon: string
      children: EntityDef[]
      exits: Record<string, string>
    }
  }

  type OrbisComponents = Partial<{ [K in keyof OrbisComponentMap]: OrbisComponentMap[K] }>

  type Event = {
    type: string

    // functions
    print: (role: string, message: string) => void
  }

  type Entity = {
    id: string
    name: string
    description: string
    aliases: string[]
    tags?: string[]
    components?: OrbisComponents
    reactions?: Record<string, (ev: Event) => void>

  }

  namespace Orbis {
    function load(entities: Record<string, EntityDef>): void
  }

  const Orbis: typeof Orbis
}
