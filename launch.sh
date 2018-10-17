#!/bin/sh
#

# A POSIX variable
OPTIND=1         # Reset in case getopts has been used previously in the shell.

# Initialize our own variables:
remove=0
first=0
debug=0

while getopts "h?frdw" opt; do
    case "$opt" in
    h|\?)
        echo "$(basename ""$0"") [-h] [-?] [-f|-r] -d"
        exit 0
        ;;
	r)	remove=1
		;;
	f)	first=1
		remove=1
		;;
	d)	debug=1
    esac
done

shift $((OPTIND-1))

[ "${1:-}" = "--" ] && shift

BTCD="btcd --singlenode --generate"
if [[ $debug -ne 0 ]]; then
	BTCD=$BTCD --debuglevel=trace
fi

if [[ $remove -ne 0 ]]; then
	rm -rf "~/.btcd/data/mainnet"
fi
$BTCD --miningkey=`gpg -d ~/keys/a.gpg`
