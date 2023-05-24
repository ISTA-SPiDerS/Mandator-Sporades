This repo implements an asynchronous consensus protocol and a mem-pool.

This branch implements vanilla Mandator and Sporades (without workers)

Our protocol uses the [Mage](https://magefile.org/) build tool. Therefore it needs the ```mage``` command to be installed.


Our protocol uses [Protocol Buffers](https://developers.google.com/protocol-buffers/).
It requires the ```protoc``` compiler with the ```go``` output plugin installed.


Our protocol uses [Redis](https://redis.io/topics/quickstart) and it should be installed with default options.

All implementations are tested in ```Ubuntu 20.04.3 LTS```

Some build dependencies can be installed by running ```mage builddeps```

Download code dependencies by running ```mage deps``` or ```go mod vendor```.

Build the code using ```mage generate && mage build``` in the root directory

All the commands to run replicas and the clients are available in the respective directories


Remote repositories

Asynchronous consensus repo    : https://github.com/PasinduTennage/async-consensus

Modified Rabia repo		: https://github.com/PasinduTennage/rabia

Modified Epaxos repo		: https://github.com/PasinduTennage/epaxos

Experiments repo		: https://github.com/PasinduTennage/async_consensus_experiments


