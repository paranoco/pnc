attr "io_mode" {
  type = "string"
  required = false
}

blockmap "tables" {
  LabelNames = ["schema", "table_name"]
  TypeName = "table"

  attr "comment" {
    required = false
    type = "string"
  }

  blockmap "columns" {
    TypeName = "column"
    LabelNames = ["name"]

    attr "comment" {
      required = false
      type = "string"
    }

    attr "type" {
      required = true
      type = "string"
    }
  }

  blockmap "rules" {
    TypeName = "rule"
    LabelNames = ["name"]
  }

  blockmap "acl" {
    TypeName = "acl"
    LabelNames = ["noop"]
  }

  blockmap "triggers" {
    TypeName = "trigger"
    LabelNames = ["name"]

    attr "fires" {
      required = false
      type = "string" // enum:BEFORE,AFTER
    }

    attr "INSERT" {
      required = false
      type = "bool"
    }

    attr "UPDATE" {
      required = false
      type = "bool"
    }

    attr "DELETE" {
      required = false
      type = "bool"
    }

    attr "TRUNCATE" {
      required = false
      type = "bool"
    }
  }

  blockmap "row_policy" {
    TypeName = "row_policy"
    LabelNames = ["name"]
  }

  blockmap "index" {
    TypeName = "index"
    LabelNames = ["name"]
  }
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