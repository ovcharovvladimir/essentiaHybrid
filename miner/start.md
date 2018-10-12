

# Settings for worker
#resultQueueSize is the size of channel listening to sealing result.
export ResultQueueSize=10

# txChanSize is the size of channel listening to NewTxsEvent. // The number is referenced from the size of tx pool.
export TxChanSize=4096

# chainHeadChanSize is the size of channel listening to ChainHeadEvent.
export ChainHeadChanSize=10

# chainSideChanSize is the size of channel listening to ChainSideEvent.
export ChainSideChanSize=10

# resubmitAdjustChanSize is the size of resubmitting interval adjustment channel.
export ResubmitAdjustChanSize=10

# miningLogAtDepth is the number of confirmations before logging successful mining.
export MiningLogAtDepth=7

# minRecommitInterval is the minimal time interval to recreate the mining block with  any newly arrived transactions.
export MinRecommitInterval=1

# maxRecommitInterval is the maximum time interval to recreate the mining block with  any newly arrived transactions.
export MaxRecommitInterval=15

# intervalAdjustRatio is the impact a single interval adjustment has on sealing work resubmitting interval.
export IntervalAdjustRatio=0.1

# intervalAdjustBias is applied during the new resubmit interval calculation in favor of increasing upper limit or decreasing lower limit so that the limit can be reachable.
export IntervalAdjustBias=200000000.0

# staleThreshold is the maximum depth of the acceptable stale block.
export StaleThreshold=7

# Run gess with paramaters
./gess --testnet
