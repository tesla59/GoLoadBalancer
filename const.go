package main

const DatabaseTargetPath = "/app/worker-stats.db"
const DefaultProtocol = "tcp"
const DefaultWorkerImageTag = "worker-image"
const DefaultConfigFileName = "config.yaml"
const ConfigTargetPath = "/app/config.yaml"
const WorkerSrcPath = "./worker"

var DatabaseFileName = "worker-stats.db"
var config Config
var PWD string
var ContainerIDs []string
