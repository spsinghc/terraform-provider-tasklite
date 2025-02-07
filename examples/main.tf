terraform {
  required_providers {
    tasklite = {
      source = "registry.terraform.io/providers/tasklite"
    }
  }
}

provider "tasklite" {
  host = "http://127.0.0.1:3000"
}

resource "tasklite_task" "t1" {
   title = "Task created by terraform"
   priority = 5 # default is 0
   complete = true # default is false
}
