#
#   Read settings from test.http and execute api call. This script must be runned from tools/
# 

ADDR=$(sed -n '1p' < test.http)
AUTH=$(sed -n '2p' < test.http)

curl -v --request $ADDR --header "$AUTH"
