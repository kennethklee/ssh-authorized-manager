#!/bin/bash

export PORT=5000
export VERSION=`git describe --tags --always`

(cd app; go mod tidy)

case $1 in
  init)
    $0 up --build
    ;;

  build)
    shift
    docker-compose build "$@"
    ;;

  port)
    shift
    echo Server hosted on http://`$0 compose port app 8090 2> /dev/null`
    ;;

  version)
    echo $VERSION
    ;;

  up)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --remove-orphans "$@"
    $0 port
    ;;

  stop)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop "$@"
    ;;

  down)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml down --volumes "$@"
    ;;

  restart)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart ${1:-app web}
    ;;

  stage)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.test.yml up -d --build --remove-orphans "$@"
    $0 port
    ;;

  ps)
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml "$@"
    ;;

  logs)
    shift
    service=${1:-app}
    shift
    echo docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs --tail 100 -f "$@" ${service}
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs --tail 100 -f "$@" ${service}
    ;;

  test)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec -e CGO_ENABLED=0 app go test ./... "$@"
    # TODO test web container too
    # docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec web npm test
    ;;

  release)
    shift
    # TODO go releaser?
    docker build -t kennethkl/ssh-authorized-manager:${VERSION:-dev} .
    docker push kennethkl/ssh-authorized-manager:${VERSION:-dev}
    docker tag kennethkl/ssh-authorized-manager:${VERSION:-dev} kennethkl/ssh-authorized-manager:latest
    docker push kennethkl/ssh-authorized-manager:latest
    ;;

  multi-release)
    shift
    # TODO make work
    docker buildx build --platform linux/arm/v7,linux/arm64/v8,linux/amd64 -t kennethkl/ssh-authorized-manager:${VERSION:-dev} .
    docker push kennethkl/ssh-authorized-manager:${VERSION:-dev}
    docker tag kennethkl/ssh-authorized-manager:${VERSION:-dev} kennethkl/ssh-authorized-manager:latest
    docker push kennethkl/ssh-authorized-manager:latest
    ;;

  shell)
    shift
    service=${1:-app}
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec ${service} ${@:-sh}
    ;;
  sh)
    # alias
    $0 shell
    ;;

  app)
    shift
    $0 shell app ${@:-sh}
    ;;

  web)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec web ${@:-sh}

    # web sometimes generates files, change them back to ours
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec web chown -R $UID:$GROUPS /web
    ;;

  db)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec app sqlite3 -nullvalue '<null>' ${@:-/app/pb_data/data.db}
    ;;

  compose)
    shift
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml "$@"
    ;;

  npm)
    shift
    $0 web npm "$@"
    ;;
  npx)
    shift
    $0 web npx "$@"
    ;;

  *)
    echo "usage: $0 <command>"
    echo
    echo "Stack specific commands:"
    echo "init             initialize project and run stack"
    echo "port             show port for stack"
    echo "version          show version"
    echo "build            build docker image"
    echo "up [service]     run stack (try $0 up --build)"
    echo "stop [service]   stop stack"
    echo "down [service]   tear down stack"
    echo "restart [service] restart app"
    echo "staging          start staging stack"
    echo "ps               list running docker containers"
    echo "logs [service]   show app logs"
    echo "test <args>      run tests"
    echo "release          release docker image"
    echo
    echo "Helper commands:"
    echo "shell [service]  shell prompt"
    echo "sh [service]     shell prompt (alias)"
    echo "app <cmd>        run command on app container"
    echo "web <cmd>        run command on web container"
    echo "compose <args>   docker-compose commands"
    echo "npm <args>       npm commands"
    echo "npx <args>       npx commands"
    ;;

esac
