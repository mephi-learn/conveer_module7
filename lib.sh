#!/bin/bash

delay=1
jpg_dir=$( cd -- "$(dirname -- "${BASH_SOURCE[0]}" )" &>/dev/null && pwd )

prompt_to_top() {
  tput cup 0 0
  for ((i=0;i<${1};i++));do echo -e "\033[2K";done
  tput cup 0 0
}

printexec() {
    name="${1}"
    timestamp=$(date --iso-8601=second | awk -F+ '{gsub(/:/, "-");sub(/T/, "_"); print $1}')
    shift
    clear
    lines=$(echo -e "${hint}\n$@\n"|wc -l)
    echo -e "${hint}\n$@\n"
    if [[ "${name}" != "no-exec" ]]; then
        eval "${@}"
    fi;
    prompt_to_top $lines
    echo -e "${hint}\n$@\n"
    sleep "${delay}"
    tput civis
    xfce4-screenshooter -w -d 0 -s "${jpg_dir}/${timestamp}-${name}.jpg"
    tput cnorm
}
