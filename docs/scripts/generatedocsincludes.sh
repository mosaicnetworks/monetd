#!/bin/bash


set -eu

mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

outdir="$mydir/../_static/includes/"


# $1 is output file name without a path
# $2 is the command line
function generatedocinsert() {
    outfile=$outdir$1
    cmd="$2"

 #   echo ".. code:: bash"> $outfile
 #   echo ""  >> $outfile
    echo "[..monetd] \$ $cmd"  > $outfile
    $cmd  >> $outfile
}


generatedocinsert monetd_version.txt "monetd version"
generatedocinsert monetd_help.txt "monetd help"
generatedocinsert monetd_keys_help.txt "monetd keys help"
generatedocinsert monetd_help_keys_new.txt "monetd help keys new"
generatedocinsert monetd_help_keys_display.txt "monetd help keys display"
generatedocinsert monetd_help_keys_update.txt "monetd help keys update"
generatedocinsert monetd_help_keys_list.txt "monetd help keys list"
generatedocinsert monetd_help_config_contract.txt "monetd help config contract"
generatedocinsert monetd_help_config_location.txt "monetd help config location"
generatedocinsert monetd_help_config_build.txt "monetd help config build"
generatedocinsert monetd_help_config_pull.txt "monetd help config pull"
generatedocinsert monetd_help_run.txt "monetd help run"


generatedocinsert giverny_help_network_new.txt "giverny help network new"
generatedocinsert giverny_version.txt "giverny version"
generatedocinsert giverny_help_keys_import.txt "giverny help keys import"



# Special cases

giverny help keys | grep -A10 "Global Flags:" > $outdir"giverney_keys_flags.txt"


