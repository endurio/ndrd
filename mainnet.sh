#!/bin/sh
#

# A POSIX variable
OPTIND=1         # Reset in case getopts has been used previously in the shell.

# Initialize our own variables:
remove=0
daemon=1
wallet=1
first=0
trace=0

while getopts "h?frdwt" opt; do
    case "$opt" in
    h|\?)
        echo "$(basename ""$0"") [-h] [-?] [-f|-r] [-d|-w] [t] nodes_count"
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
	t)	trace=1
		;;
    esac
done

shift $((OPTIND-1))

[ "${1:-}" = "--" ] && shift

# OS Compatibility
PTY=winpty
START=start
function killall() {
	taskkill -F -IM "$1.exe"
}

CTL="chainctl --rpcuser=a --rpcpass=a --skipverify"
CTLW="$CTL --wallet"

BTCD="btcd --singlenode --rpcuser=a --rpcpass=a --generate"
if [[ $trace -ne 0 ]]; then
	BTCD="$BTCD --debuglevel=trace"
fi
BTCW="btcwallet --connect=localhost --username=a --password=a"

if [[ $daemon -ne 0 ]]; then
	if [[ $remove -ne 0 ]]; then
		rm -rf "$LOCALAPPDATA/btcd/data/mainnet"
	fi
	$START $BTCD --miningkey=`gpg -d ~/keys/a.gpg`

	if [[ $first -ne 0 ]]; then
		sleep 10
		rm -rf "$LOCALAPPDATA/btcwallet/mainnet"
		$PTY $BTCW --create --passphrase=0
		$START $BTCW
		sleep 5
		WALLET_ADDR=`$CTLW getnewaddress`
		killall btcwallet

		rm -rf "$LOCALAPPDATA/btcwalletTMP/mainnet"
		$PTY $BTCW --appdata="$LOCALAPPDATA/btcwalletTMP" --create --passphrase=0
		$START $BTCW --appdata="$LOCALAPPDATA/btcwalletTMP"
		sleep 5
		$CTLW walletpassphrase 0 0
		$CTLW importprivkey `gpg -d ~/keys/ndr.gpg`
		$CTLW sendfrom imported $WALLET_ADDR 4 NDR
		$CTLW importprivkey `gpg -d ~/keys/stb.gpg`
		$CTLW sendfrom imported $WALLET_ADDR 6 STB
		killall btcwallet
	fi
fi

if [[ $wallet -ne 0 ]]; then
	sleep 2
	$START $BTCW
	sleep 5
	$CTLW walletpassphrase 0 0
fi
