package internal

import "os"

var JUSTDEPLOY_FOLDER = os.Getenv("HOME") + "/xs.config/" + "justdeploy"

var CERT_DOCKER_FOLDER = JUSTDEPLOY_FOLDER + "/cert-docker"

var DATABASE_FILE_PATH = JUSTDEPLOY_FOLDER + "/database.json"
