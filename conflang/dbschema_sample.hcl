table "public" "user" {
  column "email" {
    type = "varchar(${abs(128*2-1)})"
  }

  rule "r" {

  }

  trigger "n" {
    fires = "BEFORE"
    INSERT = true
  }

  index "n" {

  }

  row_policy "n" {

  }

  acl "noop" {

  }
}