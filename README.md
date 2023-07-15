This repo implements Sporades consensus protocol and a Mandator memory pool.


Our protocol uses the [Mage](https://magefile.org/) build tool. Therefore it needs the ```mage``` command to be installed.


Our protocol uses [Protocol Buffers](https://developers.google.com/protocol-buffers/).
It requires the ```protoc``` compiler with the ```go``` output plugin installed.


Our protocol uses [Redis](https://redis.io/topics/quickstart) and it should be installed with default options.

All implementations are tested in ```Ubuntu 20.04.3 LTS```

Some build dependencies can be installed by running ```mage builddeps```

Download code dependencies by running ```mage deps``` or ```go mod vendor```.

Build the code using ```mage generate && mage build``` in the 2 directories, separately