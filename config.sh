jq --version > /dev/null 2>&1
                if [ $? -ne 0 ]; then
                        echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
                        echo
                        exit 1
                fi
  sureOrg_TOKEN=$(curl -s -X POST http://10.99.22.103:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=sureOrg&password=password&orgName=sureOrg')
 sureOrg_TOKEN=$(echo $sureOrg_TOKEN | jq ".token" | sed "s/\"//g")
 mccOrg_TOKEN=$(curl -s -X POST http://10.99.22.103:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=mccOrg&password=password&orgName=mccOrg')
 mccOrg_TOKEN=$(echo $mccOrg_TOKEN | jq ".token" | sed "s/\"//g")
 platOrg_TOKEN=$(curl -s -X POST http://10.99.22.103:4000/login -H "content-type: application/x-www-form-urlencoded" -d 'username=platOrg&password=password&orgName=platOrg')
 platOrg_TOKEN=$(echo $platOrg_TOKEN | jq ".token" | sed "s/\"//g")
  curl -s -X POST \
                  http://10.99.22.103:4000/channels \
                  -H "authorization: Bearer $sureOrg_TOKEN" \
                  -H "content-type: application/json" \
                  -d '{"channelName":"mcctvchannel"}'
 sleep 5
  curl -s -X POST \
                                http://10.99.22.103:4000/channels/peers \
                                -H "authorization: Bearer $sureOrg_TOKEN" \
                                -H "content-type: application/json" \
                                -d '{"peers": ["peer0","peer1"],"channelName":"mcctvchannel"}'
 curl -s -X POST \
                                http://10.99.22.103:4000/channels/peers \
                                -H "authorization: Bearer $mccOrg_TOKEN" \
                                -H "content-type: application/json" \
                                -d '{"peers": ["peer0","peer1"],"channelName":"mcctvchannel"}'
 curl -s -X POST \
                                http://10.99.22.103:4000/channels/peers \
                                -H "authorization: Bearer $platOrg_TOKEN" \
                                -H "content-type: application/json" \
                                -d '{"peers": ["peer0","peer1"],"channelName":"mcctvchannel"}'
  curl -s -X POST \
                        http://10.99.22.103:4000/chaincodes \
                        -H "authorization: Bearer $sureOrg_TOKEN" \
                        -H "content-type: application/json" \
                        -d '{
                          "peers": ["peer0","peer1"],
                          "channelName":"mcctvchannel",
                          "chaincodeName":"chaincode",
                          "chaincodePath":"jiakechaincode",
                          "chaincodeVersion":"v1.0"
                  }'
 curl -s -X POST \
                        http://10.99.22.103:4000/chaincodes \
                        -H "authorization: Bearer $mccOrg_TOKEN" \
                        -H "content-type: application/json" \
                        -d '{
                          "peers": ["peer0","peer1"],
                          "channelName":"mcctvchannel",
                          "chaincodeName":"chaincode",
                          "chaincodePath":"jiakechaincode",
                          "chaincodeVersion":"v1.0"
                  }'
 curl -s -X POST \
                        http://10.99.22.103:4000/chaincodes \
                        -H "authorization: Bearer $platOrg_TOKEN" \
                        -H "content-type: application/json" \
                        -d '{
                          "peers": ["peer0","peer1"],
                          "channelName":"mcctvchannel",
                          "chaincodeName":"chaincode",
                          "chaincodePath":"jiakechaincode",
                          "chaincodeVersion":"v1.0"
                  }'
  curl -s -X POST \
                http://10.99.22.103:4000/channels/chaincodes \
                -H "authorization: Bearer $sureOrg_TOKEN" \
                -H "content-type: application/json" \
                -d '{
                  "channelName":"mcctvchannel",
                  "chaincodeName":"chaincode",
                  "chaincodeVersion":"v1.0",
                  "args":[]
          }'

