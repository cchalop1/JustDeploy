package internal

import "os"

var JUSTDEPLOY_FOLDER = os.Getenv("HOME") + "/.config/" + "justdeploy"

var CERT_DOCKER_FOLDER = JUSTDEPLOY_FOLDER + "/cert-docker"

var DATABASE_SQLITE_PATH = JUSTDEPLOY_FOLDER + "/database.db"
