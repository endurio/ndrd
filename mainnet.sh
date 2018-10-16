#!/bin/sh
#

# A POSIX variable
OPTIND=1         # Reset in case getopts has been used previously in the shell.

# Initialize our own variables:
remove=0
daemon=1
wallet=1
first=0

while getopts "h?frdw" opt; do
    case "$opt" in
    h|\?)
        echo "$(basename ""$0"") [-h] [-?] [-f|-r] [-d|-w] nodes_count"
        exit 0
        ;;
	r)	remove=1
		;;
	f)	first=1
		remove=1
		;;
	d)	wallet=0
		;;
	w)	daemon=0
		;;
    esac
done

shift $((OPTIND-1))

[ "${1:-}" = "--" ] && shift

# process OPTs
if [[ $remove -ne 0 ]]; then
	rm -rf "$LOCALAPPDATA/btcd/data/mainnet"
	rm -rf "$LOCALAPPDATA/btcwallet/mainnet"
fi
rm -rf "$LOCALAPPDATA/btcwalletTMP/mainnet"

CTL="btcctl --rpcuser=a --rpcpass=a --skipverify"
CTLW="$CTL --wallet"

BTCW="btcwallet --connect=localhost --username=a --password=a --create --walletpass=password"

if [[ $daemon -ne 0 ]]; then
	# addresses and keys
	MINING_SKEY=`gpg -d ~/keys/a.gpg`
	start btcd --debuglevel=trace --rpcuser=a --rpcpass=a --generate --miningkey=$MINING_SKEY
fi

if [[ $wallet -ne 0 ]]; then
	sleep 2
	start $BTCW --appdata="$LOCALAPPDATA/btcwallet"
	sleep 5
	$CTLW walletpassphrase "password" 0
fi
