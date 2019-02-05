# Clone repositories

The provided script allows user to clone multiple repository from the output of the dataset filters.

## Prerequisites

In order to run this script, you need to install the following tools:

```
Go (1.11.4)
Git
```

## Executing

In order to run this script, you have to execute the following commands:

```
cd src
go build -o clone_repositories
./clone_repositories <output_dataset_filter>
```

<output_dataset_filter> is the path of the output file of the script which process and filter the dataset.