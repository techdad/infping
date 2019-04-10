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
  runs = "build"
  args = "-t infping ."
  needs = ["Docker Login"]
}

action "Docker Push" {
  uses = "actions/docker/cli@master"
  runs = "push"
  args = "infping"
  needs = ["Docker Build"]
}
