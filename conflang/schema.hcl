attr "io_mode" {
  type = "string"
  required = false
}

blockmap "services" {
  LabelNames = ["type", "label"]
  TypeName = "service"

  attr "listen_addr" {
    type = "string"
    required = true
  }

  blockmap "processes" {
    TypeName = "process"
    LabelNames = ["name"]

    attr "command" {
      required = true
      type = "list[string]"
    }
  }
}