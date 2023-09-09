package main

const DatabaseTargetPath = "/app/worker-stats.sql"
const DatabaseFileName = "worker-stats.sql"
const DefaultProtocol = "tcp"
const DefaultWorkerImageTag = "worker-image"
const DefaultConfigFilePath = "./config.yaml"
const WorkerSrcPath = "./worker"

var config Config
var PWD string
