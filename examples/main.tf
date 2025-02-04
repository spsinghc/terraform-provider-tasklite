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
}
