workflow "Docker Build and Push" {
  on = "push"
  resolves = ["Docker Push"]
}

action "Docker Login" {
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_REGISTRY_URL", "DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Docker Build" {
  uses = "actions/docker/cli@master"
  args = "build -t registry.roci.nvw.io/infping ."
  needs = ["Docker Login"]
}

action "Docker Push" {
  uses = "actions/docker/cli@master"
  args = "push registry.roci.nvw.io/infping"
  needs = ["Docker Build"]
}
