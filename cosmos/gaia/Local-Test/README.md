# Cosmos-Hub Staking Module
This project is a demonstration of the Cosmos-Hub with basic staking functionality staking module designed to get validators acquainted with staking concepts and procedures.

Currently, the validator set is updated every block. The validator set is determined as the validators with the top 100 bonded atoms. Currently, all bonding and unbonding is instantaneous (no queue). Absent features include, delegation, validator rewards, unbonding wait period.

## Installation
```
go get github.com/cosmos/gaia
cd $GOPATH/src/github.com/cosmos/gaia
make all
gaia version
```

### Local-Test Example

Here is a quick example to get you off your feet:

First, generate a new key with a name, and save the address:

```
MYNAME=<your name>
gaia client keys new $MYNAME
gaia client keys list
MYADDR=<your newly generated address>
```
Now initialize a gaia chain:

```
gaia node init $MYADDR --home=$HOME/.gaia1 --chain-id=test
```

This will create all the files necessary to run a single node chain in `$HOME/.gaia1／config`:
a `priv_validator.json` file with the validators private key, and a `genesis.json` file
with the list of validators and accounts. In this case, we have one random validator,
and `$MYADDR` is an independent account that has a bunch of coins.

We can add a second node on our local machine by initiating a node in a new directory,
and copying in the genesis:


```
gaia node init $MYADDR --home=$HOME/.gaia2 --chain-id=test
cp $HOME/.gaia1/config/genesis.json $HOME/.gaia2/config/genesis.json
```

We need to also modify `$HOME/.gaia2/config/config.toml` to set new seeds and ports. It should look like:

```
proxy_app = "tcp://127.0.0.1:46668"
moniker = "anonymous"
fast_sync = true
db_backend = "leveldb"
log_level = "main:info,state:info,*:error"

[rpc]
laddr = "tcp://0.0.0.0:46667"

[p2p]
laddr = "tcp://0.0.0.0:46666"
seeds = "ID@0.0.0.0:46656"
```

Great, now that we've initialized the chains, we can start both nodes in the background:

```
gaia node start --home=$HOME/.gaia1  &> gaia1.log &
NODE1_PID=$!
gaia node start --home=$HOME/.gaia2  &> gaia2.log &
NODE2_PID=$!
```

Note we save the `PID` so we can later kill the processes.

Of course, you can peak at your logs with `tail gaia1.log`, or follow them
for a bit with `tail -f gaia1.log`.

Now we can initialize a client for the first node, and look up our account:

```
gaia client init --chain-id=test --node=tcp://localhost:46657
gaia client query account $MYADDR
```
Nice. we can send tx to another account:
```
MYNAME1=<your name>
gaia client keys new $MYNAME1
gaia client keys list
MYADDR1=<your newly generated address>
gaia client tx send --amount=1000fermion --to=$MYADDR1 --name=$MYNAME
gaia client query $MYADDR1
gaia client query $MYADDR
```



We can also lookup the candidate/validator set:

```
gaia client query candidates
```

Notice it's empty! This is because the initial validators are special -
the app doesn't know about them, so they can't be removed. To see what
tendermint itself thinks the validator set is, use:

```
curl localhost:46657/validators
```

Ok, let's add a validator-candidate. First, we need the pubkey data:

```
cat $HOME/.gaia1/config/priv_validator.json
```

If you have a json parser like `jq`, you can get just the pubkey:

```
cat $HOME/.gaia1/config/priv_validator.json | jq .pub_key.data
```

Now we can create new validator-candidate account and delegate some coins to it
```
gaia client tx declare-candidacy --amount=1fermion --pubkey=<validator pubkey> --moniker=<moniker> --name=$MYNAME
```

we can query candidates set
```
gaia client query candidates
```
or query a validator-candidate account
```
gaia client query candidate --pubkey=<validator pubkey>
```
or edit candidacy
```
gaia client tx edit-candidacy --pubkey=<validator pubkey> --moniker=<new moniker> --name=$MYNAME
```

then we can delegate some coins to this pubkey

```
gaia client tx delegate --amount=1fermion --pubkey=<validator pubkey> --name=$MYNAME
```
or unbond it
```
gaia client tx unbond --shares=1 --pubkey=<validator pubkey> --name=$MYNAME
```

We should see our account balance decrement, and the pubkey get added to the app's list of bonds:

```
gaia client query account $MYADDR
gaia client query candidates
```

To confirm for certain the new validator is active, check tendermint:

```
curl localhost:46657/validators
```


Then let's add another validator-candidate with `$HOME/.gaia2/config/priv_validator.json`

If you now kill your second node, blocks will stop streaming in, because there aren't enough validators online.
Turn her back on and they will start streaming again.
