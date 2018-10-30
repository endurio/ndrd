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

# addresses and keys
MINING_ADDR=ame7CaXCbBV4YvpLUX3fNsGSd7y3ryBfKf
MINING_SKEY=Fw28Hpjon65S4XT8uyfh7w7UFxWVExTs8oDyQZXwB1fTgwwzxnVY

# process OPTs
if [[ $remove -ne 0 ]]; then
	rm -rf "$LOCALAPPDATA/btcd/data/simnet"
	rm -rf "$LOCALAPPDATA/btcwallet/simnet"
fi
rm -rf "$LOCALAPPDATA/btcwalletTMP/simnet"

CTL="ndrctl --simnet --rpcuser=a --rpcpass=a --skipverify"
CTLW="$CTL --wallet"

BTCW="btcwallet --simnet --connect=localhost --username=a --password=a --createtemp"

if [[ $daemon -ne 0 ]]; then
	start btcd --simnet --rpcuser=a --rpcpass=a --miningkey=$MINING_SKEY

	if [[ $first -ne 0 ]]; then
		start $BTCW --appdata="$LOCALAPPDATA/btcwallet"
		sleep 5
		WALLET_ADDR=`$CTLW getnewaddress`
		taskkill -IM btcwallet.exe

		start $BTCW --appdata="$LOCALAPPDATA/btcwalletTMP"
		sleep 5
		$CTLW walletpassphrase "password" 0
		$CTLW importprivkey $MINING_SKEY
		$CTLW sendfrom imported $WALLET_ADDR 4 NDR
		$CTLW sendfrom imported $WALLET_ADDR 6 STB
		$CTL generate 1
		taskkill -IM btcwallet.exe
	fi
fi

if [[ $wallet -ne 0 ]]; then
	sleep 2
	start $BTCW --appdata="$LOCALAPPDATA/btcwallet"
	sleep 5
	$CTLW walletpassphrase "password" 0
fi
