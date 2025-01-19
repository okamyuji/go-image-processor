#!/usr/bin/env bash

if [[ $# != 0 ]]; then
    echo '[Usage]
./lint.sh

[Options]
-h/--help    #Display this help message

[Exit Status]
0: Command completed successfully
1: Command failed

[Specifications]
- Errors detected by `goimports` will be automatically fixed, but `git add` must be done manually
  (For exit status purposes, successful auto-fixing is considered "success")

- Other errors are not automatically fixed and require manual correction

- If you encounter false positives in `typos`, please edit `_typos.toml`

- If you encounter false positives in `codespell`, please edit `.codespellrc`'

    exit 0
fi

function check_command_existence() {
    local commands=(
        goimports
        staticcheck
        golangci-lint
        typos
        codespell
    )
    local exit_status=0
    for command in "${commands[@]}"; do
        local does_command_exist=1
        command -v "${command}" > /dev/null || does_command_exist=0
        if [[ "${does_command_exist}" == 0 ]]; then
            echo "Command \`${command}\` is not installed"
            exit_status=1
        fi
    done
    if [[ "${exit_status}" != 0 ]]; then
        exit "${exit_status}"
    fi
}
check_command_existence

function print_header() {
    echo "========== $1 =========="
}

exit_status=0

function tear_down() {
    if [[ $? != 0 ]]; then
        exit_status=1
    fi
    echo
}

print_header '[formatter] goimports -l -w .'
goimports -l -w .
tear_down

print_header '[linter] go vet ./...'
go vet ./...
tear_down

print_header '[linter] staticcheck ./...'
staticcheck ./...
tear_down

print_header '[linter] golangci-lint ./...'
golangci-lint run
tear_down

print_header '[spell checker] typos --format brief'
typos --format brief
tear_down

print_header '[spell checker] codespell'
codespell
tear_down

exit "${exit_status}"