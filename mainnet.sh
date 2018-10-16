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

CTL="btcctl --rpcuser=a --rpcpass=a --skipverify"
CTLW="$CTL --wallet"

BTCD="btcd --singlenode --debuglevel=trace --rpcuser=a --rpcpass=a --generate"
BTCW="btcwallet --connect=localhost --username=a --password=a"

if [[ $daemon -ne 0 ]]; then
	if [[ $remove -ne 0 ]]; then
		rm -rf "$LOCALAPPDATA/btcd/data/mainnet"
	fi
	start $BTCD --miningkey=`gpg -d ~/keys/a.gpg`

	if [[ $first -ne 0 ]]; then
		sleep 10
		rm -rf "$LOCALAPPDATA/btcwallet/mainnet"
		winpty $BTCW --create
		start $BTCW
		sleep 5
		WALLET_ADDR=`$CTLW getnewaddress`
		taskkill -F -IM btcwallet.exe

		rm -rf "$LOCALAPPDATA/btcwalletTMP/mainnet"
		winpty $BTCW --appdata="$LOCALAPPDATA/btcwalletTMP" --create
		start $BTCW --appdata="$LOCALAPPDATA/btcwalletTMP"
		sleep 5
		$CTLW walletpassphrase 0 0
		$CTLW importprivkey `gpg -d ~/keys/ndr.gpg`
		$CTLW sendfrom imported $WALLET_ADDR 4 NDR
		$CTLW importprivkey `gpg -d ~/keys/stb.gpg`
		$CTLW sendfrom imported $WALLET_ADDR 6 STB
		taskkill -F -IM btcwallet.exe
	fi
fi

if [[ $wallet -ne 0 ]]; then
	sleep 2
	start $BTCW
	sleep 5
	$CTLW walletpassphrase 0 0
fi
