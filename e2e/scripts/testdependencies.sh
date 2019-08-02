#!/bin/bash

# Script gives a non-zero return code for any missing programs

command -v giverny > /dev/null 2>&1 ||  { echo >&2 giverny is not installed and is required. Run \"make install\" from the root of this repo. Aborting. ; exit 1; }

command -v npm > /dev/null 2>&1 ||  { echo >&2 npm is not installed and is required. Run \"make install\" from the e2e folder. Aborting. ; exit 1; }
