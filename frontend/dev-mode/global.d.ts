declare global {
  // custom class props for vue components
  type ClassProp = string | Record<string, boolean>
  type ClassProps = ClassProp | ClassProp[]
}

export {}
