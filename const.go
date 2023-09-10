package main

const DatabaseTargetPath = "/app/worker-stats.db"
const DatabaseFileName = "worker-stats.db"
const DefaultProtocol = "tcp"
const DefaultWorkerImageTag = "worker-image"
const DefaultConfigFilePath = "./config.yaml"
const WorkerSrcPath = "./worker"

var config Config
var PWD string
var ContainerIDs []string
