#!/bin/sh

function cdx() {
    command="$(go-cdx $@)"
    eval "${command}"
}
